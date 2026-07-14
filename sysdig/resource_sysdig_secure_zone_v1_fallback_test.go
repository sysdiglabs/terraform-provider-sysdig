package sysdig

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

// fakeZoneBackend simulates the zones API of a backend. When exposeV2 is
// false it behaves like a deployment where only /platform/v1/zones is
// exposed and any /platform/v2/* request returns a plain 404.
type fakeZoneBackend struct {
	mu       sync.Mutex
	exposeV2 bool
	zones    map[int]*v2.Zone
	nextID   int

	v1Hits []string // "METHOD path" of every /platform/v1/zones request
	v2Hits []string // "METHOD path" of every /platform/v2/zones request
}

func newFakeZoneBackend(exposeV2 bool) *fakeZoneBackend {
	return &fakeZoneBackend{
		exposeV2: exposeV2,
		zones:    map[int]*v2.Zone{},
		nextID:   1,
	}
}

func (f *fakeZoneBackend) seed(zone *v2.Zone) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.zones[zone.ID] = zone
	if zone.ID >= f.nextID {
		f.nextID = zone.ID + 1
	}
}

func (f *fakeZoneBackend) has(id int) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.zones[id]
	return ok
}

func (f *fakeZoneBackend) server(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/platform/v1/zones", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.v1Hits = append(f.v1Hits, r.Method+" "+r.URL.Path)

		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		var req v2.ZoneRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode zone request: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		zone := &v2.Zone{
			ID:          f.nextID,
			Name:        req.Name,
			Description: req.Description,
			Scopes:      req.Scopes,
		}
		f.nextID++
		f.zones[zone.ID] = zone
		writeJSON(t, w, zone)
	})

	mux.HandleFunc("/platform/v1/zones/", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.v1Hits = append(f.v1Hits, r.Method+" "+r.URL.Path)

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/platform/v1/zones/"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		zone, ok := f.zones[id]

		switch r.Method {
		case http.MethodGet:
			if !ok {
				http.NotFound(w, r)
				return
			}
			writeJSON(t, w, zone)
		case http.MethodPut:
			if !ok {
				http.NotFound(w, r)
				return
			}
			var req v2.ZoneRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode zone request: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			zone.Name = req.Name
			zone.Description = req.Description
			zone.Scopes = req.Scopes
			writeJSON(t, w, zone)
		case http.MethodDelete:
			if !ok {
				http.NotFound(w, r)
				return
			}
			delete(f.zones, id)
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	})

	if f.exposeV2 {
		toV2 := func(zone *v2.Zone) *v2.ZoneV2 {
			out := &v2.ZoneV2{
				ID:          zone.ID,
				Name:        zone.Name,
				Description: zone.Description,
			}
			scope := v2.ScopeV2{}
			for _, s := range zone.Scopes {
				scope.Filters = append(scope.Filters, v2.FilterV2{
					ID:           s.ID,
					ResourceType: s.TargetType,
					Rules:        s.Rules,
				})
			}
			if len(scope.Filters) > 0 {
				out.Scopes = []v2.ScopeV2{scope}
			}
			return out
		}

		fromV2 := func(zone *v2.ZoneV2) *v2.Zone {
			out := &v2.Zone{
				ID:          zone.ID,
				Name:        zone.Name,
				Description: zone.Description,
			}
			for _, s := range zone.Scopes {
				for _, filter := range s.Filters {
					out.Scopes = append(out.Scopes, v2.ZoneScope{
						ID:         filter.ID,
						TargetType: filter.ResourceType,
						Rules:      filter.Rules,
					})
				}
			}
			return out
		}

		mux.HandleFunc("/platform/v2/zones", func(w http.ResponseWriter, r *http.Request) {
			f.mu.Lock()
			defer f.mu.Unlock()
			f.v2Hits = append(f.v2Hits, r.Method+" "+r.URL.Path)

			if r.Method != http.MethodPost {
				http.NotFound(w, r)
				return
			}
			var req v2.ZoneV2
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode zone v2 request: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			zone := fromV2(&req)
			zone.ID = f.nextID
			f.nextID++
			f.zones[zone.ID] = zone
			writeJSON(t, w, toV2(zone))
		})

		mux.HandleFunc("/platform/v2/zones/", func(w http.ResponseWriter, r *http.Request) {
			f.mu.Lock()
			defer f.mu.Unlock()
			f.v2Hits = append(f.v2Hits, r.Method+" "+r.URL.Path)

			id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/platform/v2/zones/"))
			if err != nil {
				http.NotFound(w, r)
				return
			}
			zone, ok := f.zones[id]

			switch r.Method {
			case http.MethodGet:
				if !ok {
					http.NotFound(w, r)
					return
				}
				writeJSON(t, w, toV2(zone))
			case http.MethodPut:
				if !ok {
					http.NotFound(w, r)
					return
				}
				var req v2.ZoneV2
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("failed to decode zone v2 request: %v", err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				upd := fromV2(&req)
				zone.Name = upd.Name
				zone.Description = upd.Description
				zone.Scopes = upd.Scopes
				writeJSON(t, w, toV2(zone))
			case http.MethodDelete:
				if !ok {
					http.NotFound(w, r)
					return
				}
				delete(f.zones, id)
				w.WriteHeader(http.StatusNoContent)
			default:
				http.NotFound(w, r)
			}
		})
	}

	// Anything else (including /platform/v2/* when exposeV2 is false)
	// falls through to the mux default 404.
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

// writeJSON encodes v as the JSON response body. It reports failures with
// t.Errorf, which is safe to call from the server goroutine (unlike
// require.NoError, whose FailNow must run on the test goroutine).
func writeJSON(t *testing.T, w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Errorf("failed to encode response: %v", err)
	}
}

func zoneTestClients(t *testing.T, url string) (v2.ZoneInterface, v2.ZoneV2Interface) {
	t.Helper()
	client := v2.NewSysdigSecure(v2.WithURL(url), v2.WithToken("test-token"))
	return client, client
}

// v2-syntax rules (no labels/labelValues/agentTags): categorizeZone treats
// this as v2 even though the backend may only expose the v1 API.
func zoneRulesResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceSysdigSecureZone().Schema, map[string]interface{}{
		"name":        "Empty zone",
		"description": "Exclude all cluster and host runtime scan results",
		"scope": []interface{}{
			map[string]interface{}{
				"target_type": "kubernetes",
				"rules":       `clusterId in ("non-existent")`,
			},
		},
	})
}

func zoneExpressionResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceSysdigSecureZone().Schema, map[string]interface{}{
		"name": "Expression zone",
		"scope": []interface{}{
			map[string]interface{}{
				"target_type": "kubernetes",
				"expression": []interface{}{
					map[string]interface{}{
						"field":    "clusterId",
						"operator": "in",
						"values":   []interface{}{"prod"},
					},
				},
			},
		},
	})
}

func TestZoneV1Fallback_CreateUsesV1WhenV2NotExposed(t *testing.T) {
	backend := newFakeZoneBackend(false)
	srv := backend.server(t)
	clientV1, clientV2 := zoneTestClients(t, srv.URL)

	d := zoneRulesResourceData(t)

	id, diags := createZone(context.Background(), d, clientV1, clientV2)
	require.False(t, diags.HasError(), "diags: %v", diags)
	require.True(t, backend.has(id), "zone should have been created through the v1 endpoint")
}

func TestZoneV1Fallback_CreateWithExpressionsFailsClearly(t *testing.T) {
	backend := newFakeZoneBackend(false)
	srv := backend.server(t)
	clientV1, clientV2 := zoneTestClients(t, srv.URL)

	d := zoneExpressionResourceData(t)

	_, diags := createZone(context.Background(), d, clientV1, clientV2)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "/platform/v2/zones")
}

func TestZoneV1Fallback_ReadFallsBackToV1(t *testing.T) {
	backend := newFakeZoneBackend(false)
	backend.seed(&v2.Zone{
		ID:   22,
		Name: "Empty zone",
		Scopes: []v2.ZoneScope{
			{ID: 1, TargetType: "kubernetes", Rules: `clusterId in ("non-existent")`},
		},
	})
	srv := backend.server(t)
	clientV1, clientV2 := zoneTestClients(t, srv.URL)

	d := zoneRulesResourceData(t)
	d.SetId("22")

	diags := readZone(context.Background(), d, clientV1, clientV2)
	require.False(t, diags.HasError(), "diags: %v", diags)
	require.Equal(t, "22", d.Id(), "zone must not be removed from state when only the v2 endpoint is missing")
	require.Equal(t, "Empty zone", d.Get("name"))
}

func TestZoneV1Fallback_ReadRemovesFromStateWhenGoneEverywhere(t *testing.T) {
	backend := newFakeZoneBackend(false)
	srv := backend.server(t)
	clientV1, clientV2 := zoneTestClients(t, srv.URL)

	d := zoneRulesResourceData(t)
	d.SetId("22")

	diags := readZone(context.Background(), d, clientV1, clientV2)
	require.False(t, diags.HasError(), "diags: %v", diags)
	require.Empty(t, d.Id(), "zone missing on both endpoints must be removed from state")
}

func TestZoneV1Fallback_UpdateFallsBackToV1(t *testing.T) {
	backend := newFakeZoneBackend(false)
	backend.seed(&v2.Zone{
		ID:   22,
		Name: "Old name",
		Scopes: []v2.ZoneScope{
			{ID: 1, TargetType: "kubernetes", Rules: `clusterId in ("non-existent")`},
		},
	})
	srv := backend.server(t)
	clientV1, clientV2 := zoneTestClients(t, srv.URL)

	d := zoneRulesResourceData(t)
	d.SetId("22")

	diags := updateZone(context.Background(), d, clientV1, clientV2)
	require.False(t, diags.HasError(), "diags: %v", diags)

	backend.mu.Lock()
	name := backend.zones[22].Name
	backend.mu.Unlock()
	require.Equal(t, "Empty zone", name, "update should have gone through the v1 endpoint")
}

func TestZoneV1Fallback_DeleteFallsBackToV1(t *testing.T) {
	backend := newFakeZoneBackend(false)
	backend.seed(&v2.Zone{
		ID:   22,
		Name: "Empty zone",
		Scopes: []v2.ZoneScope{
			{ID: 1, TargetType: "kubernetes", Rules: `clusterId in ("non-existent")`},
		},
	})
	srv := backend.server(t)
	clientV1, clientV2 := zoneTestClients(t, srv.URL)

	d := zoneRulesResourceData(t)
	d.SetId("22")

	diags := deleteZone(context.Background(), d, clientV1, clientV2)
	require.False(t, diags.HasError(), "diags: %v", diags)
	require.False(t, backend.has(22), "zone should have been deleted through the v1 endpoint")
}

// Regression guard for SaaS: when the v2 endpoint is available, the provider
// must keep using it and never touch v1 for non-legacy zones.
func TestZoneV2Available_ReadDoesNotTouchV1(t *testing.T) {
	backend := newFakeZoneBackend(true)
	backend.seed(&v2.Zone{
		ID:   22,
		Name: "Empty zone",
		Scopes: []v2.ZoneScope{
			{ID: 1, TargetType: "kubernetes", Rules: `clusterId in ("non-existent")`},
		},
	})
	srv := backend.server(t)
	clientV1, clientV2 := zoneTestClients(t, srv.URL)

	d := zoneRulesResourceData(t)
	d.SetId("22")

	diags := readZone(context.Background(), d, clientV1, clientV2)
	require.False(t, diags.HasError(), "diags: %v", diags)
	require.Equal(t, "22", d.Id())

	backend.mu.Lock()
	defer backend.mu.Unlock()
	require.Empty(t, backend.v1Hits, "v1 endpoint must not be called when v2 is available")
	require.NotEmpty(t, backend.v2Hits)
}

// Regression guard for SaaS: deleting a zone that is already gone must still
// succeed (DeleteZoneV2 returns 404, the v1 fallback tolerates it).
func TestZoneV2Available_DeleteAlreadyGoneSucceeds(t *testing.T) {
	backend := newFakeZoneBackend(true)
	srv := backend.server(t)
	clientV1, clientV2 := zoneTestClients(t, srv.URL)

	d := zoneRulesResourceData(t)
	d.SetId("22")

	diags := deleteZone(context.Background(), d, clientV1, clientV2)
	require.False(t, diags.HasError(), "diags: %v", diags)
	require.Empty(t, d.Id())
}

// A zone whose state uses expression blocks cannot be represented by the v1
// API: the read fallback must error explicitly instead of silently rewriting
// the scopes as rules.
func TestZoneV1Fallback_ReadWithExpressionsFailsClearly(t *testing.T) {
	backend := newFakeZoneBackend(false)
	backend.seed(&v2.Zone{
		ID:   22,
		Name: "Expression zone",
		Scopes: []v2.ZoneScope{
			{ID: 1, TargetType: "kubernetes", Rules: `clusterId in ("prod")`},
		},
	})
	srv := backend.server(t)
	clientV1, clientV2 := zoneTestClients(t, srv.URL)

	d := zoneExpressionResourceData(t)
	d.SetId("22")

	diags := readZone(context.Background(), d, clientV1, clientV2)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "/platform/v2/zones")
	require.Equal(t, "22", d.Id(), "zone must not be removed from state")
}

// Regression guard for backends with v2 available: create must go through the
// v2 endpoint and never touch v1.
func TestZoneV2Available_CreateDoesNotTouchV1(t *testing.T) {
	backend := newFakeZoneBackend(true)
	srv := backend.server(t)
	clientV1, clientV2 := zoneTestClients(t, srv.URL)

	d := zoneRulesResourceData(t)

	id, diags := createZone(context.Background(), d, clientV1, clientV2)
	require.False(t, diags.HasError(), "diags: %v", diags)
	require.True(t, backend.has(id))

	backend.mu.Lock()
	defer backend.mu.Unlock()
	require.Empty(t, backend.v1Hits, "v1 endpoint must not be called when v2 is available")
	require.NotEmpty(t, backend.v2Hits)
}

// Regression guard for backends with v2 available: update must go through the
// v2 endpoint and never touch v1.
func TestZoneV2Available_UpdateDoesNotTouchV1(t *testing.T) {
	backend := newFakeZoneBackend(true)
	backend.seed(&v2.Zone{
		ID:   22,
		Name: "Old name",
		Scopes: []v2.ZoneScope{
			{ID: 1, TargetType: "kubernetes", Rules: `clusterId in ("non-existent")`},
		},
	})
	srv := backend.server(t)
	clientV1, clientV2 := zoneTestClients(t, srv.URL)

	d := zoneRulesResourceData(t)
	d.SetId("22")

	diags := updateZone(context.Background(), d, clientV1, clientV2)
	require.False(t, diags.HasError(), "diags: %v", diags)

	backend.mu.Lock()
	defer backend.mu.Unlock()
	require.Equal(t, "Empty zone", backend.zones[22].Name, "update should have gone through the v2 endpoint")
	require.Empty(t, backend.v1Hits, "v1 endpoint must not be called when v2 is available")
	require.NotEmpty(t, backend.v2Hits)
}

// A zone configured with expression blocks cannot be represented by the v1
// API: the update fallback must error explicitly instead of PUTting scopes
// with empty rules to v1, which would wipe the zone's scopes server-side.
func TestZoneV1Fallback_UpdateWithExpressionsFailsClearly(t *testing.T) {
	backend := newFakeZoneBackend(false)
	backend.seed(&v2.Zone{
		ID:   22,
		Name: "Expression zone",
		Scopes: []v2.ZoneScope{
			{ID: 1, TargetType: "kubernetes", Rules: `clusterId in ("prod")`},
		},
	})
	srv := backend.server(t)
	clientV1, clientV2 := zoneTestClients(t, srv.URL)

	d := zoneExpressionResourceData(t)
	d.SetId("22")

	diags := updateZone(context.Background(), d, clientV1, clientV2)
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "/platform/v2/zones")

	backend.mu.Lock()
	defer backend.mu.Unlock()
	require.Equal(t, `clusterId in ("prod")`, backend.zones[22].Scopes[0].Rules,
		"zone scopes must not be modified through the v1 endpoint")
}
