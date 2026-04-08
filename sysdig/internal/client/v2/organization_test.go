//go:build unit

package v2

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
)

func TestMarshalOrg(t *testing.T) {
	t.Parallel()
	c := Client{}

	given := &OrganizationSecure{
		cloudauth.CloudOrganization{
			ManagementAccountId: "58ca66a5-ac87-497b-a501-7a4c934b3017",
			Provider:            cloudauth.Provider_PROVIDER_GCP,
		},
	}
	expected := `{"managementAccountId":"58ca66a5-ac87-497b-a501-7a4c934b3017", "provider":"PROVIDER_GCP"}`

	payload, err := c.marshalCloudauthProto(given)
	if err != nil {
		t.Errorf("failed to marshal payload, err: %v", err)
	}

	buf := &strings.Builder{}
	_, err = io.Copy(buf, payload)
	if err != nil {
		t.Errorf("failed to populate buffer, err: %v", err)
	}
	marshaled := buf.String()

	if strings.ReplaceAll(marshaled, " ", "") != strings.ReplaceAll(expected, " ", "") {
		t.Errorf("expected %v, got %v", expected, marshaled)
	}
}

func TestUpdateOrganizationSecureUsesCamelCase(t *testing.T) {
	t.Parallel()

	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		// Return a valid protobuf JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"managementAccountId":"31ebd166-82ef-4ca4-baf9-ce760afa46fb","organizationRootId":"r-8llw"}`))
	}))
	defer server.Close()

	c := newSysdigClient(
		WithURL(server.URL),
		WithToken("test-token"),
	)

	org := &OrganizationSecure{
		cloudauth.CloudOrganization{
			ManagementAccountId: "31ebd166-82ef-4ca4-baf9-ce760afa46fb",
			OrganizationRootId:  "r-8llw",
		},
	}

	_, _, err := c.UpdateOrganizationSecure(context.Background(), "e69f12fd-934d-43cc-8b8c-2964aba20003", org)
	if err != nil {
		t.Fatalf("UpdateOrganizationSecure failed: %v", err)
	}

	if strings.Contains(receivedBody, "management_account_id") {
		t.Errorf("request body uses snake_case (json.Marshal), expected camelCase (protojson.Marshal): %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "managementAccountId") {
		t.Errorf("request body missing camelCase field 'managementAccountId': %s", receivedBody)
	}
	if strings.Contains(receivedBody, "organization_root_id") {
		t.Errorf("request body uses snake_case 'organization_root_id', expected camelCase 'organizationRootId': %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "organizationRootId") {
		t.Errorf("request body missing camelCase field 'organizationRootId': %s", receivedBody)
	}
}

func TestUnmarshalOrg(t *testing.T) {
	t.Parallel()
	c := Client{}
	given := `{"managementAccountId":"58ca66a5-ac87-497b-a501-7a4c934b3017","provider":"PROVIDER_GCP"}`
	expected := &OrganizationSecure{
		cloudauth.CloudOrganization{
			ManagementAccountId: "58ca66a5-ac87-497b-a501-7a4c934b3017",
			Provider:            cloudauth.Provider_PROVIDER_GCP,
		},
	}

	unmarshalled := &OrganizationSecure{}
	err := c.unmarshalCloudauthProto(io.NopCloser(strings.NewReader(given)), unmarshalled)
	if err != nil {
		t.Errorf("got error while unmarshaling, err: %v", err)
	}

	if expected.String() != unmarshalled.String() {
		t.Errorf("expected %v, got %v", expected, unmarshalled)
	}
}
