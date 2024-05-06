//go:build unit

package v2

import (
	"io"
	"strings"
	"testing"

	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
)

func TestMarshalProto(t *testing.T) {
	t.Parallel()
	c := Client{}
	given := &CloudauthAccountSecure{
		CloudAccount: cloudauth.CloudAccount{
			Enabled:    true,
			ProviderId: "test-project",
			Provider:   cloudauth.Provider_PROVIDER_GCP,
		},
	}
	expected := `{"enabled":true, "providerId":"test-project", "provider":"PROVIDER_GCP"}`

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

func TestUnmarshalProto(t *testing.T) {
	t.Parallel()
	c := Client{}
	given := `{"enabled":true, "providerId":"test-project", "provider":"PROVIDER_GCP"}`
	expected := &CloudauthAccountSecure{
		CloudAccount: cloudauth.CloudAccount{
			Enabled:    true,
			ProviderId: "test-project",
			Provider:   cloudauth.Provider_PROVIDER_GCP,
		},
	}

	unmarshalled := &CloudauthAccountSecure{}
	err := c.unmarshalCloudauthProto(io.NopCloser(strings.NewReader(given)), unmarshalled)
	if err != nil {
		t.Errorf("got error while unmarshaling, err: %v", err)
	}

	if expected.String() != unmarshalled.String() {
		t.Errorf("expected %v, got %v", expected, unmarshalled)
	}
}
