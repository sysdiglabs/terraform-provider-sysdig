//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureCloudAuthAccount(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accID := rText()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureCloudAuthAccountMinimumConfiguration(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureCloudAuthAccountMinimumConfiguration(accountID string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "sample" {
  provider_id   = "sample-%s"
  provider_type = "PROVIDER_GCP"
  enabled       = "true"
}`, accountID)
}

func TestAccSecureCloudAuthAccountFC(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accID := rText()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureCloudAuthAccountWithFC(accID),
			},
			{
				ResourceName:            "sysdig_secure_cloud_auth_account.sample-1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"component"},
			},
		},
	})
}

func secureCloudAuthAccountWithFC(accountID string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "sample-1" {
  provider_id   = "sample-1-%s"
  provider_type = "PROVIDER_GCP"
  enabled       = "true"
  feature {
	secure_config_posture {
	  enabled    = "true"
	  components = ["COMPONENT_SERVICE_PRINCIPAL/secure-posture"]
	}
	secure_identity_entitlement {
	  enabled    = true
	  components = ["COMPONENT_SERVICE_PRINCIPAL/secure-posture"]
	}
  }
  component {
	type                       = "COMPONENT_SERVICE_PRINCIPAL"
	instance                   = "secure-posture"
	service_principal_metadata = jsonencode({
      gcp = {
        key = "%s"
      }
    })
  }
}
`, accountID, getEncodedServiceAccountKey("sample-1", accountID))
}

func getEncodedServiceAccountKey(resourceName string, accountID string) string {
	type sample_service_account_key struct {
		Type                    string `json:"type"`
		ProjectId               string `json:"project_id"`
		PrivateKeyId            string `json:"private_key_id"`
		PrivateKey              string `json:"private_key"`
		ClientEmail             string `json:"client_email"`
		ClientId                string `json:"client_id"`
		AuthUri                 string `json:"auth_uri"`
		TokenUri                string `json:"token_uri"`
		AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
		ClientX509CertUrl       string `json:"client_x509_cert_url"`
		UniverseDomain          string `json:"universe_domain"`
	}
	test_service_account_key := &sample_service_account_key{
		Type:                    "service_account",
		ProjectId:               fmt.Sprintf("%s-%s", resourceName, accountID),
		PrivateKeyId:            "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		PrivateKey:              "-----BEGIN PRIVATE KEY-----\nxxxxxxxxxxxxxxxxxxxxxxxxxxx\n-----END PRIVATE KEY-----\n",
		ClientEmail:             fmt.Sprintf("some-sa-name@%s-%s.iam.gserviceaccount.com", resourceName, accountID),
		ClientId:                "some-client-id",
		AuthUri:                 "https://some-auth-uri",
		TokenUri:                "https://some-token-uri",
		AuthProviderX509CertUrl: "https://some-authprovider-cert-url",
		ClientX509CertUrl:       "https://some-client-cert-url",
		UniverseDomain:          "googleapis.com",
	}

	test_service_account_key_bytes, _ := json.Marshal(test_service_account_key)
	var out bytes.Buffer
	json.Indent(&out, test_service_account_key_bytes, "", "  ")
	out.WriteByte('\n')

	test_service_account_key_encoded := b64.StdEncoding.EncodeToString(out.Bytes())
	return test_service_account_key_encoded
}
