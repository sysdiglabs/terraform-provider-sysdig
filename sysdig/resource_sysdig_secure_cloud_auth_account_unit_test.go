//go:build unit

package sysdig

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
)

const (
	testCloudauthAccountID    = "a1b2c3d4-0000-0000-0000-000000000001"
	testPostureComponentID    = "COMPONENT_TRUSTED_ROLE/secure-posture"
	testOnboardingComponentID = "COMPONENT_TRUSTED_ROLE/secure-onboarding"
	testRuntimeComponentID    = "COMPONENT_EVENT_BRIDGE/secure-runtime"
)

func trustedRoleComponent(instance string) *cloudauth.AccountComponent {
	return &cloudauth.AccountComponent{
		Type:     cloudauth.Component_COMPONENT_TRUSTED_ROLE,
		Instance: instance,
		Metadata: &cloudauth.AccountComponent_TrustedRoleMetadata{TrustedRoleMetadata: &cloudauth.TrustedRoleMetadata{
			Provider: &cloudauth.TrustedRoleMetadata_Aws{Aws: &cloudauth.TrustedRoleMetadata_AWS{RoleName: "sysdig-secure"}},
		}},
	}
}

// storedCloudauthAccount carries what the resource state cannot: feature flags, a feature type the
// resource does not map, and a component no feature references.
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
				Components: []string{testPostureComponentID},
				Flags:      map[string]string{"POSTURE_VALIDATED_EXPOSURE_ENABLED": "true"},
			},
			SecureIdentityEntitlement: &cloudauth.AccountFeature{
				Type:       cloudauth.Feature_FEATURE_SECURE_IDENTITY_ENTITLEMENT,
				Enabled:    true,
				Components: []string{testPostureComponentID},
				Flags:      map[string]string{"CIEM_FEATURE_MODE": "basic"},
			},
			SecureWorkloadScanningContainers: &cloudauth.AccountFeature{
				Type:       cloudauth.Feature_FEATURE_SECURE_WORKLOAD_SCANNING_CONTAINERS,
				Enabled:    true,
				Components: []string{testPostureComponentID},
			},
		},
		Components: []*cloudauth.AccountComponent{
			trustedRoleComponent("secure-posture"),
			{Type: cloudauth.Component_COMPONENT_EVENT_BRIDGE, Instance: "secure-runtime"},
		},
	}
}

func testSysdigClients(t *testing.T, url string) *sysdigClients {
	t.Helper()
	return &sysdigClients{ctx: context.Background(), d: providerData(t, map[string]any{
		"sysdig_secure_url":       url,
		"sysdig_secure_api_token": "fake-token",
	})}
}

// serveStoredCloudauthAccount serves stored like cloudauth does and returns the body of the account PUT.
func serveStoredCloudauthAccount(t *testing.T, stored *cloudauth.CloudAccount) (*sysdigClients, func() *cloudauth.CloudAccount) {
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
			// without decryption the metadata carries placeholders, which must never be written back
			if r.URL.Query().Get("decrypt") != "true" {
				t.Errorf("GET %s, want decrypt=true", r.URL)
			}
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
	t.Cleanup(srv.Close)

	return testSysdigClients(t, srv.URL), func() *cloudauth.CloudAccount {
		t.Helper()
		if put == nil {
			t.Fatal("update sent no PUT")
		}
		return put
	}
}

// plannedAccountUpdate builds the ResourceData Terraform hands to Update: the prior state applied from
// prior, diffed against the planned configuration.
func plannedAccountUpdate(t *testing.T, prior, planned map[string]any) *schema.ResourceData {
	t.Helper()
	r := resourceSysdigSecureCloudauthAccount()
	priorData := schema.TestResourceDataRaw(t, r.Schema, prior)
	priorData.SetId(testCloudauthAccountID)
	state := priorData.State()

	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(planned), nil)
	if err != nil {
		t.Fatal(err)
	}
	data, err := schema.InternalMap(r.Schema).Data(state, diff)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func accountConfig(alias string, features map[string]any, components ...map[string]any) map[string]any {
	config := map[string]any{
		SchemaCloudProviderID:    "123456789012",
		SchemaCloudProviderType:  cloudauth.Provider_PROVIDER_AWS.String(),
		SchemaEnabled:            true,
		SchemaCloudProviderAlias: alias,
	}
	if features != nil {
		config[SchemaFeature] = []any{features}
	}
	if len(components) > 0 {
		list := []any{}
		for _, component := range components {
			list = append(list, component)
		}
		config[SchemaComponent] = list
	}
	return config
}

func featureConfig(components ...any) []any {
	return []any{map[string]any{SchemaEnabled: true, SchemaComponents: components}}
}

func trustedRoleComponentConfig(instance string) map[string]any {
	return map[string]any{
		SchemaType:                cloudauth.Component_COMPONENT_TRUSTED_ROLE.String(),
		SchemaInstance:            instance,
		SchemaTrustedRoleMetadata: `{"aws":{"role_name":"sysdig-secure"}}`,
	}
}

func runAccountUpdate(t *testing.T, clients *sysdigClients, data *schema.ResourceData) {
	t.Helper()
	if diags := resourceSysdigSecureCloudauthAccountUpdate(context.Background(), data, clients); diags.HasError() {
		t.Fatalf("update returned an error: %v", diags)
	}
}

func TestCloudauthAccountUpdateSendsStoredBlocksWithoutPlannedChange(t *testing.T) {
	stored := storedCloudauthAccount()
	clients, put := serveStoredCloudauthAccount(t, stored)
	features := map[string]any{SchemaSecureConfigPosture: featureConfig(testPostureComponentID)}
	component := trustedRoleComponentConfig("secure-posture")
	data := plannedAccountUpdate(t, accountConfig("old", features, component), accountConfig("renamed", features, component))
	if data.HasChange(SchemaFeature) || data.HasChange(SchemaComponent) {
		t.Fatal("the plan must change only the alias")
	}

	runAccountUpdate(t, clients, data)

	if !proto.Equal(put().GetFeature(), stored.GetFeature()) {
		t.Errorf("PUT features = %v, want the stored ones with their flags", put().GetFeature())
	}
	if len(put().GetComponents()) != len(stored.GetComponents()) {
		t.Errorf("PUT components = %v, want the stored ones", put().GetComponents())
	}
	for i, component := range put().GetComponents() {
		if !proto.Equal(component, stored.GetComponents()[i]) {
			t.Errorf("PUT component %d = %v, want %v", i, component, stored.GetComponents()[i])
		}
	}
	if put().GetProviderAlias() != "renamed" {
		t.Errorf("PUT provider alias = %q, want the planned one", put().GetProviderAlias())
	}
}

func TestCloudauthAccountUpdateMergesStoredIntoPlannedFeatures(t *testing.T) {
	clients, put := serveStoredCloudauthAccount(t, storedCloudauthAccount())
	component := trustedRoleComponentConfig("secure-posture")
	prior := map[string]any{
		SchemaSecureConfigPosture:       featureConfig(testPostureComponentID),
		SchemaSecureIdentityEntitlement: featureConfig(testPostureComponentID),
	}
	planned := map[string]any{
		SchemaSecureConfigPosture:   featureConfig(testPostureComponentID),
		SchemaSecureThreatDetection: featureConfig(testRuntimeComponentID),
	}
	data := plannedAccountUpdate(t, accountConfig("", prior, component), accountConfig("", planned, component))

	runAccountUpdate(t, clients, data)

	want := &cloudauth.AccountFeatures{
		SecureConfigPosture: &cloudauth.AccountFeature{
			Type:       cloudauth.Feature_FEATURE_SECURE_CONFIG_POSTURE,
			Enabled:    true,
			Components: []string{testPostureComponentID},
			Flags:      map[string]string{"POSTURE_VALIDATED_EXPOSURE_ENABLED": "true"},
		},
		SecureThreatDetection: &cloudauth.AccountFeature{
			Type:       cloudauth.Feature_FEATURE_SECURE_THREAT_DETECTION,
			Enabled:    true,
			Components: []string{testRuntimeComponentID},
		},
		SecureWorkloadScanningContainers: storedCloudauthAccount().GetFeature().GetSecureWorkloadScanningContainers(),
	}
	if !proto.Equal(put().GetFeature(), want) {
		t.Errorf("PUT features = %v, want %v: planned features with stored flags, the removed one dropped, the unmapped one kept", put().GetFeature(), want)
	}
}

func TestCloudauthAccountUpdateKeepsReferencedComponents(t *testing.T) {
	clients, put := serveStoredCloudauthAccount(t, storedCloudauthAccount())
	data := plannedAccountUpdate(t,
		accountConfig("", nil, trustedRoleComponentConfig("secure-posture")),
		accountConfig("", nil, trustedRoleComponentConfig("secure-onboarding")))

	runAccountUpdate(t, clients, data)

	var got []string
	for _, component := range put().GetComponents() {
		got = append(got, component.GetType().String()+"/"+component.GetInstance())
	}
	// the stored features still point at secure-posture, while nothing references the event bridge
	want := []string{testOnboardingComponentID, testPostureComponentID}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("PUT components = %v, want %v", got, want)
	}
}

func TestMergeStoredAccountFeaturesWithoutStoredFeatures(t *testing.T) {
	planned := &cloudauth.AccountFeatures{SecureConfigPosture: &cloudauth.AccountFeature{Enabled: true}}

	if got := mergeStoredAccountFeatures(planned, nil); !proto.Equal(got, planned) {
		t.Errorf("merged = %v, want the planned features", got)
	}
}

// serveCloudauthAccountFeature answers the feature endpoint with the given statuses and counts PUTs.
func serveCloudauthAccountFeature(t *testing.T, getStatus, putStatus int) (*sysdigClients, *int) {
	t.Helper()
	puts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := getStatus
		if r.Method == http.MethodPut {
			puts++
			status = putStatus
		}
		if status != http.StatusOK {
			http.Error(w, `{"error":"the requested entity was not found"}`, status)
			return
		}
		_, _ = w.Write([]byte(`{"type":"FEATURE_SECURE_WORKLOAD_SCANNING_CONTAINERS","enabled":true}`))
	}))
	t.Cleanup(srv.Close)
	return testSysdigClients(t, srv.URL), &puts
}

func featureData(t *testing.T) *schema.ResourceData {
	t.Helper()
	data := schema.TestResourceDataRaw(t, resourceSysdigSecureCloudauthAccountFeature().Schema, map[string]any{
		SchemaAccountID:  testCloudauthAccountID,
		SchemaType:       cloudauth.Feature_FEATURE_SECURE_WORKLOAD_SCANNING_CONTAINERS.String(),
		SchemaEnabled:    true,
		SchemaComponents: []any{testPostureComponentID},
	})
	data.SetId(testCloudauthAccountID + "/" + cloudauth.Feature_FEATURE_SECURE_WORKLOAD_SCANNING_CONTAINERS.String())
	return data
}

// A feature deleted outside Terraform has to leave the state, so the next plan recreates it.
func TestCloudauthAccountFeatureReadDropsDeletedFeature(t *testing.T) {
	clients, _ := serveCloudauthAccountFeature(t, http.StatusNotFound, http.StatusOK)
	data := featureData(t)

	if diags := resourceSysdigSecureCloudauthAccountFeatureRead(context.Background(), data, clients); diags.HasError() {
		t.Fatalf("read returned an error: %v", diags)
	}
	if data.Id() != "" {
		t.Errorf("id = %q, want it cleared", data.Id())
	}
}

func TestCloudauthAccountFeatureUpdateWritesDeletedFeature(t *testing.T) {
	clients, puts := serveCloudauthAccountFeature(t, http.StatusNotFound, http.StatusOK)

	if diags := resourceSysdigSecureCloudauthAccountFeatureUpdate(context.Background(), featureData(t), clients); diags.HasError() {
		t.Fatalf("update returned an error: %v", diags)
	}
	if *puts != 1 {
		t.Errorf("PUTs = %d, want the feature written again", *puts)
	}
}

func TestCloudauthAccountFeatureUpdateReportsMissingAccount(t *testing.T) {
	clients, _ := serveCloudauthAccountFeature(t, http.StatusOK, http.StatusNotFound)

	if diags := resourceSysdigSecureCloudauthAccountFeatureUpdate(context.Background(), featureData(t), clients); !diags.HasError() {
		t.Error("update succeeded, want the failed PUT reported")
	}
}
