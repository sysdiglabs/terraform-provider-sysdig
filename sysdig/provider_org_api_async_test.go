//go:build unit

package sysdig

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
)

const orgAPIAsyncKey = "sysdig_secure_org_api_async"

// clearOrgAPIAsyncEnv keeps the ambient environment out of the tests: the variables below are the
// ones an engineer working on large-organization onboarding has exported.
func clearOrgAPIAsyncEnv(t *testing.T) {
	t.Helper()
	for _, name := range orgAPIAsyncEnvVars {
		t.Setenv(name, "")
	}
}

func providerData(t *testing.T, config map[string]any) *schema.ResourceData {
	t.Helper()

	sm := schema.InternalMap(Provider().Schema)
	diff, err := sm.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil, nil, true)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	d, err := sm.Data(nil, diff)
	if err != nil {
		t.Fatalf("Data: %v", err)
	}
	return d
}

// Config has to beat the environment, and the diff layer only coerces to TypeBool after the
// default func runs, so an explicit false could otherwise lose to the variable.
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

// An unparseable value must stay disabled rather than fail provider configuration, which is what
// schema.EnvDefaultFunc would do with it.
func TestProviderOrgAPIAsyncEnvValues(t *testing.T) {
	cases := map[string]bool{
		"true": true, "TRUE": true, "True": true, "1": true, "t": true,
		"": false, "false": false, "0": false, "f": false,
		"yes": false, "on": false, "enabled": false, "true ": false,
	}

	for value, want := range cases {
		t.Run(value, func(t *testing.T) {
			clearOrgAPIAsyncEnv(t)
			t.Setenv("SYSDIG_ORG_API_ASYNC", value)

			if got := providerData(t, map[string]any{}).Get(orgAPIAsyncKey); got != want {
				t.Errorf("env %q resolved to %#v, want %v", value, got, want)
			}
		})
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
	data := schema.TestResourceDataRaw(t, resourceSysdigSecureOrganization().Schema, map[string]any{})
	data.SetId(orgID)

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
	data := schema.TestResourceDataRaw(t, resourceSysdigSecureOrganization().Schema, map[string]any{})

	diags := resourceSysdigSecureOrganizationCreate(context.Background(), data, clients)
	if !diags.HasError() {
		t.Fatal("expected an error when the acknowledgement carries no id")
	}
	if data.Id() != "" {
		t.Errorf("id = %q, want empty so the resource is not tracked", data.Id())
	}
}
