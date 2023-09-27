//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
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
	type sample_service_account_key struct {
		ProjectId    string `json:"project_id"`
		PrivateKeyId string `json:"private_key_id"`
		PrivateKey   string `json:"private_key"`
	}
	test_service_account_key := &sample_service_account_key{
		ProjectId:    fmt.Sprintf("sample-1-%s", accountID),
		PrivateKeyId: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		PrivateKey:   "-----BEGIN PRIVATE KEY-----\nxxxxxxxxxxxxxxxxxxxxxxxxxxx\n-----END PRIVATE KEY-----\n",
	}
	test_service_account_keyJSON, _ := json.Marshal(test_service_account_key)
	test_service_account_key_encoded := b64.StdEncoding.EncodeToString([]byte(string(test_service_account_keyJSON)))

	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "sample-1" {
  provider_id   = "sample-1-%s"
  provider_type = "PROVIDER_GCP"
  enabled       = "true"
  feature {
	secure_config_posture {
	  enabled    = "true"
	  components = ["COMPONENT_SERVICE_PRINCIPAL/secure-service-principal"]
	}
  }
  component {
	type                       = "COMPONENT_SERVICE_PRINCIPAL"
	instance                   = "secure-service-principal"
	service_principal_metadata = jsonencode({
      gcp = {
        key = "%s"
      }
    })
  }
}
`, accountID, test_service_account_key_encoded)
}
