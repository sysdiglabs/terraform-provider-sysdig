//go:build unit

package sysdig

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
)

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
