package sysdig

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
)

const orgAPIAsyncKey = "sysdig_secure_org_api_async"

// clearOrgAPIAsyncEnv keeps the ambient environment out of the tests: the variables below are the
// ones an engineer working on large-organization onboarding has exported.
func clearOrgAPIAsyncEnv(t *testing.T) {
	t.Helper()
	for _, env := range orgAPIAsyncEnvVars {
		t.Setenv(env.name, "")
	}
}

func providerData(t *testing.T, config map[string]any) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, Provider().Schema, config)
}

// Config has to beat the environment, and the diff layer only coerces to TypeBool after the
// default func runs, so an explicit false could otherwise lose to the variable.
// TestResourceDataRaw does not attach the resource's Timeouts, so data.Timeout() would fall back
// to the SDK's generic default and a regression in the delete budget would go unnoticed. Going
// through the resource itself pins the real value into whatever the code under test reads.
func organizationData(t *testing.T, id string, attrs map[string]string) *schema.ResourceData {
	t.Helper()
	res := resourceSysdigSecureOrganization()
	data := res.Data(&terraform.InstanceState{ID: id, Attributes: attrs})
	if got, want := data.Timeout(schema.TimeoutDelete), 30*time.Minute; got != want {
		t.Fatalf("delete timeout the wait will use = %v, want %v", got, want)
	}
	return data
}

func TestProviderOrgAPIAsync(t *testing.T) {
	tests := []struct {
		name   string
		env    map[string]string
		config map[string]any
		want   bool
	}{
		{name: "unset everywhere", config: map[string]any{}, want: false},
		{name: "legacy env var enables it", env: map[string]string{"SYSDIG_ORG_API_ASYNC": "true"}, config: map[string]any{}, want: true},
		{name: "prefixed env var enables it", env: map[string]string{"SYSDIG_SECURE_ORG_API_ASYNC": "true"}, config: map[string]any{}, want: true},
		{
			name:   "prefixed env var wins over the legacy one",
			env:    map[string]string{"SYSDIG_SECURE_ORG_API_ASYNC": "false", "SYSDIG_ORG_API_ASYNC": "true"},
			config: map[string]any{},
			want:   false,
		},
		{
			name:   "explicit false beats env var",
			env:    map[string]string{"SYSDIG_ORG_API_ASYNC": "true"},
			config: map[string]any{orgAPIAsyncKey: false},
			want:   false,
		},
		{name: "explicit true with no env var", config: map[string]any{orgAPIAsyncKey: true}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearOrgAPIAsyncEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got, ok := providerData(t, tt.config).Get(orgAPIAsyncKey).(bool)
			if !ok {
				t.Fatalf("%s did not resolve to a bool", orgAPIAsyncKey)
			}
			if got != tt.want {
				t.Errorf("%s = %v, want %v", orgAPIAsyncKey, got, tt.want)
			}
		})
	}
}

// The legacy name keeps its exact-string rule so an upgrade cannot flip anyone from sync to async,
// while the new name accepts the conventional forms. Neither may fail provider configuration.
func TestProviderOrgAPIAsyncEnvValues(t *testing.T) {
	legacy := map[string]bool{
		"true": true,
		"":     false,
		"1":    false, // exact match only: this was a no-op before the attribute existed
		"TRUE": false,
		"yes":  false,
	}
	prefixed := map[string]bool{
		"true":  true,
		"1":     true,
		"TRUE":  true,
		"":      false,
		"yes":   false, // unparseable stays disabled instead of failing provider configuration
		"true ": false, // not trimmed
	}

	for value, want := range legacy {
		t.Run("legacy/"+value, func(t *testing.T) {
			clearOrgAPIAsyncEnv(t)
			t.Setenv("SYSDIG_ORG_API_ASYNC", value)
			if got := providerData(t, map[string]any{}).Get(orgAPIAsyncKey); got != want {
				t.Errorf("legacy env %q resolved to %#v, want %v", value, got, want)
			}
		})
	}
	for value, want := range prefixed {
		t.Run("prefixed/"+value, func(t *testing.T) {
			clearOrgAPIAsyncEnv(t)
			t.Setenv("SYSDIG_SECURE_ORG_API_ASYNC", value)
			if got := providerData(t, map[string]any{}).Get(orgAPIAsyncKey); got != want {
				t.Errorf("prefixed env %q resolved to %#v, want %v", value, got, want)
			}
		})
	}
}

// An unparseable value in the new name must not shadow a working legacy one.
func TestProviderOrgAPIAsyncUnparseableFallsThrough(t *testing.T) {
	clearOrgAPIAsyncEnv(t)
	t.Setenv("SYSDIG_SECURE_ORG_API_ASYNC", "tru")
	t.Setenv("SYSDIG_ORG_API_ASYNC", "true")

	if got := providerData(t, map[string]any{}).Get(orgAPIAsyncKey); got != true {
		t.Errorf("resolved to %#v, want the legacy value to win over an unparseable one", got)
	}
}

// Covers the whole seam: provider schema -> getSysdigSecureVariables -> client options -> the
// query string actually put on the wire.
func TestOrgAPIAsyncReachesTheWire(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodGet}

	tests := []struct {
		name      string
		config    map[string]any
		wantQuery map[string]string
	}{
		{
			name:   "disabled leaves every request untouched",
			config: map[string]any{orgAPIAsyncKey: false},
			wantQuery: map[string]string{
				http.MethodPost: "", http.MethodPut: "", http.MethodDelete: "", http.MethodGet: "",
			},
		},
		{
			name:   "enabled marks the mutating requests only",
			config: map[string]any{orgAPIAsyncKey: true},
			wantQuery: map[string]string{
				http.MethodPost:   "async=true",
				http.MethodPut:    "async=true",
				http.MethodDelete: "async=true",
				http.MethodGet:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearOrgAPIAsyncEnv(t)

			query := map[string]string{}
			sent := map[string]bool{}
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				query[r.Method] = r.URL.RawQuery
				sent[r.Method] = true
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"id":"4c53102d","managementAccountId":"m"}`))
			}))
			defer srv.Close()

			config := map[string]any{"sysdig_secure_url": srv.URL, "sysdig_secure_api_token": "fake-token"}
			for k, v := range tt.config {
				config[k] = v
			}

			clients := &sysdigClients{ctx: context.Background(), d: providerData(t, config)}
			client, err := clients.sysdigSecureClientV2()
			if err != nil {
				t.Fatalf("sysdigSecureClientV2: %v", err)
			}

			const orgID = "4c53102d"
			org := &v2.OrganizationSecure{}

			if _, _, err := client.CreateOrganizationSecure(context.Background(), org); err != nil {
				t.Fatalf("CreateOrganizationSecure: %v", err)
			}
			if _, _, err := client.UpdateOrganizationSecure(context.Background(), orgID, org); err != nil {
				t.Fatalf("UpdateOrganizationSecure: %v", err)
			}
			// a 404 delete is the short-circuit, so it never reaches the deletion wait
			if _, err := client.DeleteOrganizationSecure(context.Background(), orgID); err == nil {
				t.Fatal("DeleteOrganizationSecure on 404: expected an error")
			}
			if _, _, err := client.GetOrganizationSecure(context.Background(), orgID); err != nil {
				t.Fatalf("GetOrganizationSecure: %v", err)
			}

			for _, method := range methods {
				if !sent[method] {
					t.Errorf("%s was never sent", method)
					continue
				}
				if got, want := query[method], tt.wantQuery[method]; got != want {
					t.Errorf("%s query = %q, want %q", method, got, want)
				}
			}
		})
	}
}

// An async delete is only acknowledged, so destroy has to wait for the organization to go away.
// Driven through the resource so that dropping the wait from Delete fails here.
func TestOrganizationDeleteWaitsForRemoval(t *testing.T) {
	const orgID = "4c53102d"

	var gets int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			gets++
			if gets < 3 {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"id":"` + orgID + `"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	clearOrgAPIAsyncEnv(t)
	config := map[string]any{
		"sysdig_secure_url":       srv.URL,
		"sysdig_secure_api_token": "fake-token",
		orgAPIAsyncKey:            true,
	}
	clients := &sysdigClients{ctx: context.Background(), d: providerData(t, config)}
	data := organizationData(t, orgID, nil)

	if diags := resourceSysdigSecureOrganizationDelete(context.Background(), data, clients); diags.HasError() {
		t.Fatalf("delete returned an error: %v", diags)
	}
	if gets < 3 {
		t.Errorf("polled %d times, expected the delete to keep polling until the organization was gone", gets)
	}
}

// A 202 without a usable body must not surface as a proto parse error, but create still needs an id.
func TestOrganizationCreateRejectsAckWithoutID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// a body that parses cleanly but carries no id, so the guard is what rejects it
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	clearOrgAPIAsyncEnv(t)
	config := map[string]any{
		"sysdig_secure_url":       srv.URL,
		"sysdig_secure_api_token": "fake-token",
		orgAPIAsyncKey:            true,
	}
	clients := &sysdigClients{ctx: context.Background(), d: providerData(t, config)}
	data := schema.TestResourceDataRaw(t, resourceSysdigSecureOrganization().Schema, map[string]any{})

	diags := resourceSysdigSecureOrganizationCreate(context.Background(), data, clients)
	if !diags.HasError() {
		t.Fatal("expected an error when the acknowledgement carries no id")
	}
	if !strings.Contains(diags[0].Summary, "no id was returned") {
		t.Errorf("diagnostic = %q, want the missing-id guard rather than a decode failure", diags[0].Summary)
	}
	if data.Id() != "" {
		t.Errorf("id = %q, want empty so the resource is not tracked", data.Id())
	}
}

// A permanent failure on the confirmation read must not be retried until the delete timeout.
func TestOrganizationDeleteWaitClassifiesErrors(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		wantGets    int
		wantSucceed bool
	}{
		{name: "gone", status: http.StatusNotFound, wantGets: 1, wantSucceed: true},
		{name: "forbidden fails fast", status: http.StatusForbidden, wantGets: 1},
		{name: "unauthorized fails fast", status: http.StatusUnauthorized, wantGets: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gets int
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				gets++
				w.WriteHeader(tt.status)
			}))
			defer srv.Close()

			client := v2.NewSysdigSecure(v2.WithURL(srv.URL), v2.WithToken("fake-token"))
			err := waitForOrganizationDeletion(context.Background(), client, "4c53102d", 2*time.Second)

			if tt.wantSucceed && err != nil {
				t.Fatalf("expected success, got %v", err)
			}
			if !tt.wantSucceed && err == nil {
				t.Fatal("expected an error")
			}
			if gets < tt.wantGets {
				t.Errorf("issued %d reads, expected at least %d", gets, tt.wantGets)
			}
			if tt.wantGets == 1 && gets > 1 {
				t.Errorf("issued %d reads, expected the failure to be reported without retrying", gets)
			}
		})
	}
}

// The wait has to give up once the timeout is exhausted, reporting what it was waiting for.
func TestOrganizationDeleteWaitTimesOut(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"4c53102d"}`))
	}))
	defer srv.Close()

	client := v2.NewSysdigSecure(v2.WithURL(srv.URL), v2.WithToken("fake-token"))
	err := waitForOrganizationDeletion(context.Background(), client, "4c53102d", time.Second)
	if err == nil {
		t.Fatal("expected the wait to give up")
	}
	if !strings.Contains(err.Error(), "has not been deleted yet") {
		t.Errorf("error = %q, want it to say the organization was still there", err)
	}
}

// The async cascade outlasts the other operations, so only delete gets the longer budget.
func TestOrganizationDeleteTimeoutBudget(t *testing.T) {
	timeouts := resourceSysdigSecureOrganization().Timeouts
	if got, want := *timeouts.Delete, 30*time.Minute; got != want {
		t.Errorf("delete timeout = %v, want %v", got, want)
	}
	if got, want := *timeouts.Create, 5*time.Minute; got != want {
		t.Errorf("create timeout = %v, want %v", got, want)
	}
}

// An organization deleted outside Terraform has to leave state, so a later plan recreates it.
func TestOrganizationReadDropsDeletedOrganization(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	clients := &sysdigClients{ctx: context.Background(), d: providerData(t, map[string]any{
		"sysdig_secure_url":       srv.URL,
		"sysdig_secure_api_token": "fake-token",
	})}
	data := schema.TestResourceDataRaw(t, resourceSysdigSecureOrganization().Schema, map[string]any{})
	data.SetId("4c53102d")

	if diags := resourceSysdigSecureOrganizationRead(context.Background(), data, clients); diags.HasError() {
		t.Fatalf("read returned an error: %v", diags)
	}
	if data.Id() != "" {
		t.Errorf("id = %q, want empty so the resource leaves state", data.Id())
	}
}

// The PUT already returns the persisted organization, so Update must take state from that response
// instead of paying for a second read that could also outlive the update timeout.
func TestOrganizationUpdateUsesResponseBody(t *testing.T) {
	const managementAccountID = "58ca66a5-ac87-497b-a501-7a4c934b3017"

	var methods []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"4c53102d","managementAccountId":"` + managementAccountID + `"}`))
	}))
	defer srv.Close()

	clients := &sysdigClients{ctx: context.Background(), d: providerData(t, map[string]any{
		"sysdig_secure_url":       srv.URL,
		"sysdig_secure_api_token": "fake-token",
	})}
	data := schema.TestResourceDataRaw(t, resourceSysdigSecureOrganization().Schema, map[string]any{})
	data.SetId("4c53102d")

	if diags := resourceSysdigSecureOrganizationUpdate(context.Background(), data, clients); diags.HasError() {
		t.Fatalf("update returned an error: %v", diags)
	}
	if len(methods) != 1 || methods[0] != http.MethodPut {
		t.Errorf("requests = %v, want just the PUT", methods)
	}
	if got := data.Get(SchemaManagementAccountID); got != managementAccountID {
		t.Errorf("management account id = %q, want the value the update returned", got)
	}
}

// A 404 on the update itself is the same situation as one on the read-back, and has to be
// reported the same way rather than recorded as a successful update.
func TestOrganizationUpdateReportsMissingOrganization(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	clients := &sysdigClients{ctx: context.Background(), d: providerData(t, map[string]any{
		"sysdig_secure_url":       srv.URL,
		"sysdig_secure_api_token": "fake-token",
	})}
	data := schema.TestResourceDataRaw(t, resourceSysdigSecureOrganization().Schema, map[string]any{})
	data.SetId("4c53102d")

	diags := resourceSysdigSecureOrganizationUpdate(context.Background(), data, clients)
	if !diags.HasError() {
		t.Fatal("expected an error when the organization no longer exists")
	}
	if data.Id() == "" {
		t.Error("id was cleared, which Terraform rejects as an inconsistent result for an update")
	}
}

// A synchronous delete is complete when it returns, so confirming it would only cost a call
// and would fail the destroy on a 403 from a token that may no longer need read access.
func TestOrganizationDeleteSkipsWaitWhenSync(t *testing.T) {
	var methods []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	clearOrgAPIAsyncEnv(t)
	clients := &sysdigClients{ctx: context.Background(), d: providerData(t, map[string]any{
		"sysdig_secure_url":       srv.URL,
		"sysdig_secure_api_token": "fake-token",
		orgAPIAsyncKey:            false,
	})}
	data := organizationData(t, "4c53102d", nil)

	if diags := resourceSysdigSecureOrganizationDelete(context.Background(), data, clients); diags.HasError() {
		t.Fatalf("delete returned an error: %v", diags)
	}
	if len(methods) != 1 || methods[0] != http.MethodDelete {
		t.Errorf("requests = %v, want just the DELETE with no confirmation read", methods)
	}
}

// The transport already retries 409 and 5xx on its own, so those statuses cannot be driven through
// a test server in reasonable time; the classification is asserted directly instead.
func TestOrganizationDeleteProbeClassification(t *testing.T) {
	tests := []struct {
		name          string
		status        string
		err           error
		wantDone      bool
		wantRetryable bool
	}{
		{name: "not found is done", status: "404 Not Found", err: errors.New("nope"), wantDone: true},
		{name: "gone is done", status: "410 Gone", err: errors.New("nope"), wantDone: true},
		{name: "conflict retries", status: "409 Conflict", err: errors.New("nope"), wantRetryable: true},
		{name: "too many requests retries", status: "429 Too Many Requests", err: errors.New("nope"), wantRetryable: true},
		{name: "server error retries", status: "503 Service Unavailable", err: errors.New("nope"), wantRetryable: true},
		{name: "not implemented is fatal", status: "501 Not Implemented", err: errors.New("nope")},
		{name: "unexpected success keeps waiting", status: "202 Accepted", err: errors.New("nope"), wantRetryable: true},
		{name: "forbidden is fatal", status: "403 Forbidden", err: errors.New("nope")},
		{name: "bad request is fatal", status: "400 Bad Request", err: errors.New("nope")},
		{name: "transport failure retries", err: &url.Error{Op: "Get", Err: errors.New("connection reset")}, wantRetryable: true},
		{name: "bad scheme is fatal", err: &url.Error{Op: "Get", Err: errors.New(`unsupported protocol scheme "foo"`)}},
		{name: "untrusted certificate is fatal", err: &url.Error{Op: "Get", Err: &tls.CertificateVerificationError{Err: x509.UnknownAuthorityError{}}}},
		// a body that could not be read says nothing about the organization, and often clears
		{name: "unreadable answer keeps waiting", err: errors.New("unable to read response body"), wantRetryable: true},
		// one that arrived whole and does not parse will not parse on the next attempt either
		{name: "malformed answer is fatal", err: &v2.MalformedBodyError{Err: errors.New("unexpected EOF")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyDeletionProbe("4c53102d", tt.status, tt.err)
			switch {
			case tt.wantDone:
				if got != nil {
					t.Fatalf("classified as %+v, want the organization treated as gone", got)
				}
			case got == nil:
				t.Fatal("classified as gone, want the poll to keep going or fail")
			case got.Retryable != tt.wantRetryable:
				t.Errorf("retryable = %v, want %v (%v)", got.Retryable, tt.wantRetryable, got.Err)
			}
		})
	}
}

// The client tolerates an acknowledgement without a body, so Update must not map that zero value
// into state: it would clear the management account and the collections the plan just set, which
// Terraform reports as an inconsistent result rather than as the server's answer.
func TestOrganizationUpdateKeepsStateOnBodylessAck(t *testing.T) {
	const managementAccountID = "58ca66a5-ac87-497b-a501-7a4c934b3017"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	clearOrgAPIAsyncEnv(t)
	clients := &sysdigClients{ctx: context.Background(), d: providerData(t, map[string]any{
		"sysdig_secure_url":       srv.URL,
		"sysdig_secure_api_token": "fake-token",
		orgAPIAsyncKey:            true,
	})}
	data := schema.TestResourceDataRaw(t, resourceSysdigSecureOrganization().Schema, map[string]any{
		SchemaManagementAccountID: managementAccountID,
		SchemaAutomaticOnboarding: true,
	})
	data.SetId("4c53102d")

	if diags := resourceSysdigSecureOrganizationUpdate(context.Background(), data, clients); diags.HasError() {
		t.Fatalf("update returned an error: %v", diags)
	}
	if got := data.Get(SchemaManagementAccountID); got != managementAccountID {
		t.Errorf("management account id = %q, want the planned value kept", got)
	}
	if got := data.Get(SchemaAutomaticOnboarding); got != true {
		t.Errorf("automatic onboarding = %v, want the planned value kept", got)
	}
}

// An acknowledgement is not a statement of final state, whatever it happens to carry: a body with
// an id and a management account but nothing else must not wipe the collections from state.
func TestOrganizationUpdateKeepsStateOnPartialAck(t *testing.T) {
	const managementAccountID = "58ca66a5-ac87-497b-a501-7a4c934b3017"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"id":"4c53102d","managementAccountId":"` + managementAccountID + `"}`))
	}))
	defer srv.Close()

	clearOrgAPIAsyncEnv(t)
	clients := &sysdigClients{ctx: context.Background(), d: providerData(t, map[string]any{
		"sysdig_secure_url":       srv.URL,
		"sysdig_secure_api_token": "fake-token",
		orgAPIAsyncKey:            true,
	})}
	data := organizationData(t, "4c53102d", map[string]string{
		SchemaManagementAccountID:   managementAccountID,
		SchemaAutomaticOnboarding:   "true",
		"organizational_unit_ids.#": "1",
		"organizational_unit_ids.0": "ou-1234",
	})

	if diags := resourceSysdigSecureOrganizationUpdate(context.Background(), data, clients); diags.HasError() {
		t.Fatalf("update returned an error: %v", diags)
	}
	if got := data.Get(SchemaManagementAccountID); got != managementAccountID {
		t.Errorf("management account id = %q, want the planned value kept", got)
	}
	if got := data.Get(SchemaAutomaticOnboarding); got != true {
		t.Errorf("automatic onboarding = %v, want the planned value kept", got)
	}
	if got := data.Get(SchemaOrganizationalUnitIds).(*schema.Set).List(); len(got) != 1 {
		t.Errorf("organizational unit ids = %v, want the planned value kept", got)
	}
}

// A 202 on delete is only meaningful when async was requested. Accepting it with the flag off
// would end the destroy on an acknowledgement that nothing is waiting for.
func TestOrganizationDeleteAcceptsAckOnlyWhenAsync(t *testing.T) {
	for _, tt := range []struct {
		name      string
		async     bool
		wantError bool
	}{
		{name: "async accepts the acknowledgement", async: true},
		{name: "sync rejects it", wantError: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					w.WriteHeader(http.StatusNotFound) // already gone, so the wait ends at once
					return
				}
				w.WriteHeader(http.StatusAccepted)
			}))
			defer srv.Close()

			clearOrgAPIAsyncEnv(t)
			clients := &sysdigClients{ctx: context.Background(), d: providerData(t, map[string]any{
				"sysdig_secure_url":       srv.URL,
				"sysdig_secure_api_token": "fake-token",
				orgAPIAsyncKey:            tt.async,
			})}
			data := organizationData(t, "4c53102d", nil)

			diags := resourceSysdigSecureOrganizationDelete(context.Background(), data, clients)
			if got := diags.HasError(); got != tt.wantError {
				t.Fatalf("delete error = %v, want %v (%v)", got, tt.wantError, diags)
			}
		})
	}
}

// 410 means the organization is gone just as 404 does, and the deletion poll already knows that;
// Read and Delete have to agree, or a gone organization stays in state and a destroy errors out.
func TestOrganizationGoneIsHandledForBothStatuses(t *testing.T) {
	for _, status := range []int{http.StatusNotFound, http.StatusGone} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(status)
			}))
			defer srv.Close()

			clearOrgAPIAsyncEnv(t)
			clients := &sysdigClients{ctx: context.Background(), d: providerData(t, map[string]any{
				"sysdig_secure_url":       srv.URL,
				"sysdig_secure_api_token": "fake-token",
			})}

			read := organizationData(t, "4c53102d", nil)
			if diags := resourceSysdigSecureOrganizationRead(context.Background(), read, clients); diags.HasError() {
				t.Fatalf("read returned an error: %v", diags)
			}
			if read.Id() != "" {
				t.Errorf("id = %q, want empty so the resource leaves state", read.Id())
			}

			del := organizationData(t, "4c53102d", nil)
			if diags := resourceSysdigSecureOrganizationDelete(context.Background(), del, clients); diags.HasError() {
				t.Errorf("delete returned an error for an organization that is already gone: %v", diags)
			}
		})
	}
}

// A malformed confirmation read has to be reported at once rather than polled for the whole
// delete timeout, since asking again cannot make the same body parse.
func TestOrganizationDeleteWaitStopsOnMalformedAnswer(t *testing.T) {
	var gets int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		gets++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":`))
	}))
	defer srv.Close()

	client := v2.NewSysdigSecure(v2.WithURL(srv.URL), v2.WithToken("fake-token"))
	err := waitForOrganizationDeletion(context.Background(), client, "4c53102d", 10*time.Second)
	if err == nil {
		t.Fatal("expected the wait to report the malformed answer")
	}
	if gets != 1 {
		t.Errorf("issued %d reads, want the failure reported without retrying", gets)
	}
}

// G4: one confirmation read must not be able to spend the transport's whole backoff budget inside
// a single tick, or a cascade that is progressing fine runs the wait out of time.
func TestOrganizationDeleteProbeIsBounded(t *testing.T) {
	probe := &deadlineRecordingClient{errStatus: "404 Not Found"}

	if err := waitForOrganizationDeletion(context.Background(), probe, "4c53102d", 30*time.Minute); err != nil {
		t.Fatalf("wait returned an error: %v", err)
	}
	if len(probe.budgets) == 0 {
		t.Fatal("the wait never issued a confirmation read")
	}
	switch got := probe.budgets[0]; {
	case got < 0:
		t.Error("the confirmation read was given no deadline of its own, so one probe can spend the whole delete budget backing off")
	case got > organizationDeletionProbeTimeout:
		t.Errorf("first read had %v to answer in, want no more than %v out of the whole delete budget", got, organizationDeletionProbeTimeout)
	}
}

// A transport failure has to keep the wait going until the budget is out, rather than ending it
// as a permanent error and blaming the network for what is really a timeout.
func TestOrganizationDeleteWaitKeepsProbingOnTransportErrors(t *testing.T) {
	probe := &deadlineRecordingClient{err: &url.Error{Op: "Get", Err: errors.New("connection reset by peer")}}

	err := waitForOrganizationDeletion(context.Background(), probe, "4c53102d", 2*time.Second)
	if err == nil {
		t.Fatal("expected the wait to give up once the budget was out")
	}
	if len(probe.budgets) < 2 {
		t.Errorf("issued %d reads, want the wait to keep probing instead of failing on the first transport error", len(probe.budgets))
	}
}

// records the budget each confirmation read is given
type deadlineRecordingClient struct {
	v2.OrganizationSecureInterface
	errStatus string
	err       error
	budgets   []time.Duration
}

func (c *deadlineRecordingClient) GetOrganizationSecure(ctx context.Context, _ string) (*v2.OrganizationSecure, string, error) {
	if deadline, ok := ctx.Deadline(); ok {
		c.budgets = append(c.budgets, time.Until(deadline))
	} else {
		c.budgets = append(c.budgets, -1)
	}
	if c.err != nil {
		return nil, c.errStatus, c.err
	}
	return nil, c.errStatus, errors.New("organization not found")
}
