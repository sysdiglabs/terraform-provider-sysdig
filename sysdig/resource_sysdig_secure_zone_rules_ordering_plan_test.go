//go:build tf_acc_sysdig_secure

package sysdig_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

// These tests drive real plan/apply cycles against an in-process zones backend,
// so they assert on what a user actually sees rather than on the shape of a diff.
// They need no credentials: the provider is pointed at the fake server.

// orderingBackend is a zones backend for the ordering tests.
//
//   - exposeV2 == false reproduces a deployment where only /platform/v1/zones is
//     routed, so the provider goes through its v1 fallback.
//   - sortValuesOnRead makes every read return the values of an `in (...)` list in
//     a different order than they were written, reproducing a backend that does
//     not preserve the order.
type orderingBackend struct {
	mu               sync.Mutex
	exposeV2         bool
	sortValuesOnRead bool
	zones            map[int]*v2.ZoneV2
	nextZoneID       int
	nextScopeID      int
}

func newOrderingBackend(t *testing.T, exposeV2, sortValuesOnRead bool) *httptest.Server {
	t.Helper()

	b := &orderingBackend{
		exposeV2:         exposeV2,
		sortValuesOnRead: sortValuesOnRead,
		zones:            map[int]*v2.ZoneV2{},
		nextZoneID:       1,
		nextScopeID:      1000,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/platform/v1/zones", b.handleV1Collection(t))
	mux.HandleFunc("/platform/v1/zones/", b.handleV1Item(t))
	if exposeV2 {
		mux.HandleFunc("/platform/v2/zones", b.handleV2Collection(t))
		mux.HandleFunc("/platform/v2/zones/", b.handleV2Item(t))
	}

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

var orderingValueListPattern = regexp.MustCompile(`\(([^()]*)\)`)

// reorderValues returns the same rules expression with the values of every list
// in a different order, so that a test can tell a semantic change apart from a
// cosmetic one.
func reorderValues(rules string) string {
	return orderingValueListPattern.ReplaceAllStringFunc(rules, func(match string) string {
		body := strings.TrimSuffix(strings.TrimPrefix(match, "("), ")")
		values := strings.Split(body, ",")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}
		sort.Sort(sort.Reverse(sort.StringSlice(values)))
		return "(" + strings.Join(values, ", ") + ")"
	})
}

func (b *orderingBackend) store(zone *v2.ZoneV2) {
	for i := range zone.Scopes {
		for j := range zone.Scopes[i].Filters {
			zone.Scopes[i].Filters[j].ID = b.nextScopeID
			b.nextScopeID++
		}
	}
	b.zones[zone.ID] = zone
}

// read returns the stored zone as the API would serialize it.
func (b *orderingBackend) read(id int) (*v2.ZoneV2, bool) {
	stored, ok := b.zones[id]
	if !ok {
		return nil, false
	}

	out := *stored
	out.Scopes = make([]v2.ScopeV2, len(stored.Scopes))
	for i, scope := range stored.Scopes {
		filters := make([]v2.FilterV2, len(scope.Filters))
		for j, filter := range scope.Filters {
			if b.sortValuesOnRead && filter.Rules != "" {
				filter.Rules = reorderValues(filter.Rules)
			}
			filters[j] = filter
		}
		out.Scopes[i] = v2.ScopeV2{Filters: filters}
	}
	return &out, true
}

func (b *orderingBackend) handleV2Collection(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b.mu.Lock()
		defer b.mu.Unlock()

		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		var zone v2.ZoneV2
		require.NoError(t, json.NewDecoder(r.Body).Decode(&zone))
		zone.ID = b.nextZoneID
		b.nextZoneID++
		b.store(&zone)

		stored, _ := b.read(zone.ID)
		writeOrderingJSON(t, w, stored)
	}
}

func (b *orderingBackend) handleV2Item(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b.mu.Lock()
		defer b.mu.Unlock()

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/platform/v2/zones/"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			stored, ok := b.read(id)
			if !ok {
				http.NotFound(w, r)
				return
			}
			writeOrderingJSON(t, w, stored)
		case http.MethodPut:
			var zone v2.ZoneV2
			require.NoError(t, json.NewDecoder(r.Body).Decode(&zone))
			zone.ID = id
			b.store(&zone)
			stored, _ := b.read(id)
			writeOrderingJSON(t, w, stored)
		case http.MethodDelete:
			delete(b.zones, id)
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}
}

func (b *orderingBackend) handleV1Collection(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b.mu.Lock()
		defer b.mu.Unlock()

		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		var req v2.ZoneRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))

		zone := &v2.ZoneV2{ID: b.nextZoneID, Name: req.Name, Description: req.Description}
		b.nextZoneID++
		zone.Scopes = []v2.ScopeV2{{Filters: filtersFromV1(req.Scopes)}}
		b.store(zone)

		stored, _ := b.read(zone.ID)
		writeOrderingJSON(t, w, zoneV1From(stored))
	}
}

func (b *orderingBackend) handleV1Item(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b.mu.Lock()
		defer b.mu.Unlock()

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/platform/v1/zones/"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			stored, ok := b.read(id)
			if !ok {
				http.NotFound(w, r)
				return
			}
			writeOrderingJSON(t, w, zoneV1From(stored))
		case http.MethodPut:
			if _, ok := b.zones[id]; !ok {
				http.NotFound(w, r)
				return
			}
			var req v2.ZoneRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
			zone := &v2.ZoneV2{ID: id, Name: req.Name, Description: req.Description}
			zone.Scopes = []v2.ScopeV2{{Filters: filtersFromV1(req.Scopes)}}
			b.store(zone)
			stored, _ := b.read(id)
			writeOrderingJSON(t, w, zoneV1From(stored))
		case http.MethodDelete:
			if _, ok := b.zones[id]; !ok {
				http.NotFound(w, r)
				return
			}
			delete(b.zones, id)
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}
}

func filtersFromV1(scopes []v2.ZoneScope) []v2.FilterV2 {
	filters := make([]v2.FilterV2, 0, len(scopes))
	for _, scope := range scopes {
		filters = append(filters, v2.FilterV2{ResourceType: scope.TargetType, Rules: scope.Rules})
	}
	return filters
}

func zoneV1From(zone *v2.ZoneV2) *v2.Zone {
	out := &v2.Zone{ID: zone.ID, Name: zone.Name, Description: zone.Description}
	for _, scope := range zone.Scopes {
		for _, filter := range scope.Filters {
			out.Scopes = append(out.Scopes, v2.ZoneScope{
				ID:         filter.ID,
				TargetType: filter.ResourceType,
				Rules:      filter.Rules,
			})
		}
	}
	return out
}

func writeOrderingJSON(t *testing.T, w http.ResponseWriter, body any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	require.NoError(t, json.NewEncoder(w).Encode(body))
}

func orderingProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"sysdig": func() (*schema.Provider, error) {
			return sysdig.Provider(), nil
		},
	}
}

func zoneRulesConfig(url, name, rules string) string {
	return fmt.Sprintf(`
provider "sysdig" {
  sysdig_secure_url       = %q
  sysdig_secure_api_token = "fake-token"
}

resource "sysdig_secure_zone" "test" {
  name        = %q
  description = "value ordering"

  scope {
    target_type = "kubernetes"
    rules       = %q
  }
}
`, url, name, rules)
}

const (
	zoneRulesOriginalOrder = `clusterId in ("c2", "c1", "c4", "c3")`
	zoneRulesSameValues    = `clusterId in ("c1", "c3", "c2", "c4")`
	zoneRulesRealChange    = `clusterId in ("c1", "c3", "c2", "brand-new")`
)

// Reordering the values in the configuration must not plan any change. A
// configuration that builds the list from a variable whose order is not stable
// across runs would otherwise be dirty on every plan.
func TestAccZoneRulesOrderingConfigReorderIsANoOp(t *testing.T) {
	srv := newOrderingBackend(t, true, false)
	name := randomText(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: orderingProviderFactories(),
		Steps: []resource.TestStep{
			{Config: zoneRulesConfig(srv.URL, name, zoneRulesOriginalOrder)},
			{
				Config:             zoneRulesConfig(srv.URL, name, zoneRulesSameValues),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// The same must hold when it is the backend that returns the values in another
// order: a refresh alone used to be enough to make the plan permanently dirty.
func TestAccZoneRulesOrderingBackendReorderIsANoOp(t *testing.T) {
	srv := newOrderingBackend(t, true, true)
	name := randomText(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: orderingProviderFactories(),
		Steps: []resource.TestStep{
			{Config: zoneRulesConfig(srv.URL, name, zoneRulesOriginalOrder)},
			{
				Config:             zoneRulesConfig(srv.URL, name, zoneRulesOriginalOrder),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// Same as above on a deployment that only routes /platform/v1/zones, where the
// provider falls back to the v1 API.
func TestAccZoneRulesOrderingV1OnlyBackend(t *testing.T) {
	srv := newOrderingBackend(t, false, false)
	name := randomText(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: orderingProviderFactories(),
		Steps: []resource.TestStep{
			{Config: zoneRulesConfig(srv.URL, name, zoneRulesOriginalOrder)},
			{
				Config:             zoneRulesConfig(srv.URL, name, zoneRulesSameValues),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// Ignoring the order must not turn into ignoring the values.
func TestAccZoneRulesOrderingRealChangeIsApplied(t *testing.T) {
	srv := newOrderingBackend(t, true, false)
	name := randomText(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: orderingProviderFactories(),
		Steps: []resource.TestStep{
			{Config: zoneRulesConfig(srv.URL, name, zoneRulesOriginalOrder)},
			{
				Config: zoneRulesConfig(srv.URL, name, zoneRulesRealChange),
				Check: resource.TestCheckTypeSetElemNestedAttrs(
					"sysdig_secure_zone.test", "scope.*",
					map[string]string{"rules": zoneRulesRealChange},
				),
			},
			{
				Config:             zoneRulesConfig(srv.URL, name, zoneRulesRealChange),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func zoneExpressionConfig(url, name string, values ...string) string {
	quoted := make([]string, len(values))
	for i, v := range values {
		quoted[i] = fmt.Sprintf("%q", v)
	}

	return fmt.Sprintf(`
provider "sysdig" {
  sysdig_secure_url       = %q
  sysdig_secure_api_token = "fake-token"
}

resource "sysdig_secure_zone" "test" {
  name        = %q
  description = "value ordering"

  scope {
    target_type = "kubernetes"
    expression {
      field    = "kubernetes.cluster.name"
      operator = "in"
      values   = [%s]
    }
  }

  scope {
    target_type = "host"
    expression {
      field    = "host.hostName"
      operator = "in"
      values   = ["h1", "h2"]
    }
  }
}
`, url, name, strings.Join(quoted, ", "))
}

// Expression-based scopes keep the stock set hash. Two scope blocks must stay
// distinct, and updating the values of one must still work.
func TestAccZoneRulesOrderingExpressionScopesAreUnaffected(t *testing.T) {
	srv := newOrderingBackend(t, true, false)
	name := randomText(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: orderingProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: zoneExpressionConfig(srv.URL, name, "c1", "c2"),
				Check:  resource.TestCheckResourceAttr("sysdig_secure_zone.test", "scope.#", "2"),
			},
			{
				Config: zoneExpressionConfig(srv.URL, name, "c1", "brand-new"),
				Check:  resource.TestCheckResourceAttr("sysdig_secure_zone.test", "scope.#", "2"),
			},
			{
				Config:             zoneExpressionConfig(srv.URL, name, "c1", "brand-new"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
