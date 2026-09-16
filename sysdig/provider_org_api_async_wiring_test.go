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
	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
)

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

// Covers the whole seam the attribute has to cross: provider schema -> getSysdigSecureVariables ->
// client options -> the query string actually put on the wire.
func TestOrgAPIAsyncReachesTheWire(t *testing.T) {
	tests := []struct {
		name      string
		async     any
		wantQuery map[string]string
	}{
		{
			name:  "default leaves every request untouched",
			async: nil,
			wantQuery: map[string]string{
				http.MethodPost: "", http.MethodPut: "", http.MethodDelete: "", http.MethodGet: "",
			},
		},
		{
			name:  "explicit false leaves every request untouched",
			async: false,
			wantQuery: map[string]string{
				http.MethodPost: "", http.MethodPut: "", http.MethodDelete: "", http.MethodGet: "",
			},
		},
		{
			name:  "enabled marks the mutating requests only",
			async: true,
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
			seen := map[string]string{}
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				seen[r.Method] = r.URL.RawQuery
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"managementAccountId":"58ca66a5-ac87-497b-a501-7a4c934b3017"}`))
			}))
			defer srv.Close()

			config := map[string]any{
				"sysdig_secure_url":       srv.URL,
				"sysdig_secure_api_token": "fake-token",
			}
			if tt.async != nil {
				config["sysdig_secure_org_api_async"] = tt.async
			}

			clients := &sysdigClients{ctx: context.Background(), d: providerData(t, config)}
			client, err := clients.sysdigSecureClientV2()
			if err != nil {
				t.Fatalf("sysdigSecureClientV2: %v", err)
			}

			const orgID = "4c53102d-6846-447b-bfd1-4c0d5002cddf"
			org := &v2.OrganizationSecure{CloudOrganization: cloudauth.CloudOrganization{ManagementAccountId: "58ca66a5-ac87-497b-a501-7a4c934b3017"}}

			if _, _, err := client.CreateOrganizationSecure(context.Background(), org); err != nil {
				t.Fatalf("CreateOrganizationSecure: %v", err)
			}
			if _, _, err := client.UpdateOrganizationSecure(context.Background(), orgID, org); err != nil {
				t.Fatalf("UpdateOrganizationSecure: %v", err)
			}
			if _, err := client.DeleteOrganizationSecure(context.Background(), orgID); err != nil {
				t.Fatalf("DeleteOrganizationSecure: %v", err)
			}
			if _, _, err := client.GetOrganizationSecure(context.Background(), orgID); err != nil {
				t.Fatalf("GetOrganizationSecure: %v", err)
			}

			for method, want := range tt.wantQuery {
				if got := seen[method]; got != want {
					t.Errorf("%s query = %q, want %q", method, got, want)
				}
			}
		})
	}
}

// The delete path can now provoke a 202, so it has to accept one.
func TestDeleteOrganizationAcceptsAccepted(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	client := v2.NewSysdigSecure(v2.WithURL(srv.URL), v2.WithToken("fake-token"), v2.WithOrgAPIAsync(true))
	if errStatus, err := client.DeleteOrganizationSecure(context.Background(), "4c53102d-6846-447b-bfd1-4c0d5002cddf"); err != nil {
		t.Fatalf("DeleteOrganizationSecure on 202: %v (%s)", err, errStatus)
	}
}
