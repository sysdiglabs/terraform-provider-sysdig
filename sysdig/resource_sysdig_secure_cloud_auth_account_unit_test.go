package sysdig

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
)

const testCloudauthAccountID = "a1b2c3d4-0000-0000-0000-000000000001"

// storedCloudauthAccount carries what the resource state cannot: feature flags and feature types the
// resource does not map.
func storedCloudauthAccount() *cloudauth.CloudAccount {
	return &cloudauth.CloudAccount{
		Id:         testCloudauthAccountID,
		Enabled:    true,
		Provider:   cloudauth.Provider_PROVIDER_AWS,
		ProviderId: "123456789012",
		Feature: &cloudauth.AccountFeatures{
			SecureConfigPosture: &cloudauth.AccountFeature{
				Type:       cloudauth.Feature_FEATURE_SECURE_CONFIG_POSTURE,
				Enabled:    true,
				Components: []string{"COMPONENT_TRUSTED_ROLE/secure-posture"},
				Flags:      map[string]string{"POSTURE_VALIDATED_EXPOSURE_ENABLED": "true"},
			},
			SecureIdentityEntitlement: &cloudauth.AccountFeature{
				Type:    cloudauth.Feature_FEATURE_SECURE_IDENTITY_ENTITLEMENT,
				Enabled: true,
				Flags:   map[string]string{"CIEM_FEATURE_MODE": "advanced"},
			},
			SecureWorkloadScanningContainers: &cloudauth.AccountFeature{
				Type:    cloudauth.Feature_FEATURE_SECURE_WORKLOAD_SCANNING_CONTAINERS,
				Enabled: true,
			},
		},
		Components: []*cloudauth.AccountComponent{{
			Type:     cloudauth.Component_COMPONENT_TRUSTED_ROLE,
			Instance: "secure-posture",
			Metadata: &cloudauth.AccountComponent_TrustedRoleMetadata{TrustedRoleMetadata: &cloudauth.TrustedRoleMetadata{
				Provider: &cloudauth.TrustedRoleMetadata_Aws{Aws: &cloudauth.TrustedRoleMetadata_AWS{RoleName: "sysdig-secure"}},
			}},
		}},
	}
}

// updateCloudauthAccount runs the resource Update against a backend that serves stored and returns
// the account body of the PUT.
func updateCloudauthAccount(t *testing.T, stored *cloudauth.CloudAccount, config map[string]any) *cloudauth.CloudAccount {
	t.Helper()
	storedBody, err := protojson.Marshal(stored)
	if err != nil {
		t.Fatal(err)
	}

	var put *cloudauth.CloudAccount
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cloudauth/v1/accounts/"+testCloudauthAccountID {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write(storedBody)
		case http.MethodPut:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("reading PUT body: %v", err)
			}
			put = &cloudauth.CloudAccount{}
			if err := protojson.Unmarshal(body, put); err != nil {
				t.Errorf("decoding PUT body: %v", err)
			}
			_, _ = w.Write(body)
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
		}
	}))
	defer srv.Close()

	clients := &sysdigClients{ctx: context.Background(), d: providerData(t, map[string]any{
		"sysdig_secure_url":       srv.URL,
		"sysdig_secure_api_token": "fake-token",
	})}
	data := schema.TestResourceDataRaw(t, resourceSysdigSecureCloudauthAccount().Schema, config)
	data.SetId(testCloudauthAccountID)

	if diags := resourceSysdigSecureCloudauthAccountUpdate(context.Background(), data, clients); diags.HasError() {
		t.Fatalf("update returned an error: %v", diags)
	}
	if put == nil {
		t.Fatal("update sent no PUT")
	}
	return put
}

func baseCloudauthAccountConfig() map[string]any {
	return map[string]any{
		SchemaCloudProviderID:    "123456789012",
		SchemaCloudProviderType:  cloudauth.Provider_PROVIDER_AWS.String(),
		SchemaEnabled:            true,
		SchemaCloudProviderAlias: "renamed",
	}
}

func TestCloudauthAccountUpdateKeepsStoredFeaturesAndComponents(t *testing.T) {
	stored := storedCloudauthAccount()

	put := updateCloudauthAccount(t, stored, baseCloudauthAccountConfig())

	if !proto.Equal(put.GetFeature(), stored.GetFeature()) {
		t.Errorf("PUT features = %v, want the stored ones with their flags", put.GetFeature())
	}
	if len(put.GetComponents()) != 1 || !proto.Equal(put.GetComponents()[0], stored.GetComponents()[0]) {
		t.Errorf("PUT components = %v, want the stored ones", put.GetComponents())
	}
	if put.GetProviderAlias() != "renamed" {
		t.Errorf("PUT provider alias = %q, want the planned one", put.GetProviderAlias())
	}
}

func TestCloudauthAccountUpdateSendsPlannedFeaturesAndComponents(t *testing.T) {
	config := baseCloudauthAccountConfig()
	config[SchemaFeature] = []any{map[string]any{
		SchemaSecureThreatDetection: []any{map[string]any{
			SchemaEnabled:    true,
			SchemaComponents: []any{"COMPONENT_EVENT_BRIDGE/secure-runtime"},
		}},
	}}
	config[SchemaComponent] = []any{map[string]any{
		SchemaType:     cloudauth.Component_COMPONENT_EVENT_BRIDGE.String(),
		SchemaInstance: "secure-runtime",
	}}

	put := updateCloudauthAccount(t, storedCloudauthAccount(), config)

	want := &cloudauth.AccountFeatures{SecureThreatDetection: &cloudauth.AccountFeature{
		Type:       cloudauth.Feature_FEATURE_SECURE_THREAT_DETECTION,
		Enabled:    true,
		Components: []string{"COMPONENT_EVENT_BRIDGE/secure-runtime"},
	}}
	if !proto.Equal(put.GetFeature(), want) {
		t.Errorf("PUT features = %v, want only the planned ones", put.GetFeature())
	}
	if len(put.GetComponents()) != 1 || put.GetComponents()[0].GetType() != cloudauth.Component_COMPONENT_EVENT_BRIDGE {
		t.Errorf("PUT components = %v, want only the planned one", put.GetComponents())
	}
}
