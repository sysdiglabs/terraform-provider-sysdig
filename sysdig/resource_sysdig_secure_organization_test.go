//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"regexp"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureOrganization(t *testing.T) {
	// XXX: TF acceptance tests for secure org onboarding need an actual existing gcp project
	// along with an actual service_principal_key to scrape all folders and projects under the org.
	// Without it POST /organizations call will fail with 500 error.
	// Skipping the test based on this error when it occurs.
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accID := rText()
	organizationApiUrl := fmt.Sprintf(`%s/api/cloudauth/v1/organizations`, os.Getenv("SYSDIG_SECURE_URL"))
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
		ErrorCheck: func(err error) error {
			// if regex matches with the expected error, do t.Skip
			re := regexp.MustCompile(fmt.Sprintf(`POST %s giving up after 5 attempt(s)`, organizationApiUrl))
			if re.MatchString(err.Error()) {
				t.Skipf("skipping test; this POST call is not supported without actual existing GCP projects and service principal.")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: secureOrgWithAccountID(accID),
			},
			{
				ResourceName:            "sysdig_secure_cloud_auth_account.sample",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"component"},
			},
		},
	})
}

func secureOrgWithAccountID(accountID string) string {
	// this is a base64 encoded service account key, that should exist apriori in the project (mycapitalprojet)
	test_service_account_key_encoded := getEncodedServiceAccountKey("sample", accountID)

	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "sample" {
  provider_id   = "%s"
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
  component {
    type     				   = "COMPONENT_SERVICE_PRINCIPAL"
    instance 				   = "secure-onboarding"
    service_principal_metadata = jsonencode({
      gcp = {
        key = "%s"
      }
    })
  }
}
resource "sysdig_secure_organization" "sample-org" {
  management_account_id		= sysdig_secure_cloud_auth_account.sample.id
}
`, accountID, test_service_account_key_encoded, test_service_account_key_encoded)
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

	test_service_account_key_bytes, err := json.Marshal(test_service_account_key)
	if err != nil {
		fmt.Printf("Failed to marshal test_service_account_key: %v", err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, test_service_account_key_bytes, "", "  ")
	if err != nil {
		fmt.Printf("Failed to indent test_service_account_key: %v", err)
	}
	out.WriteByte('\n')

	test_service_account_key_encoded := b64.StdEncoding.EncodeToString(out.Bytes())
	return test_service_account_key_encoded
}
