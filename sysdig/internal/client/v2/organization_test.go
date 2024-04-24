//go:build unit

package v2

import (
	"io"
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

	if strings.Replace(marshaled, " ", "", -1) != strings.Replace(expected, " ", "", -1) {
		t.Errorf("expected %v, got %v", expected, marshaled)
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

	unmarshalled := &OrganizationSecure
	err := c.unmarshalCloudauthProto(io.NopCloser(strings.NewReader(given)), unmarshalled)
	if err != nil {
		t.Errorf("got error while unmarshaling, err: %v", err)
	}

	if expected.String() != unmarshalled.String() {
		t.Errorf("expected %v, got %v", expected, unmarshalled)
	}
}
