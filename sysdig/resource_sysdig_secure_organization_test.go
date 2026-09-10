//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
			re := regexp.MustCompile(fmt.Sprintf(`POST %s giving up after \d+ attempt\(s\)`, regexp.QuoteMeta(organizationApiUrl)))
			if re.MatchString(err.Error()) {
				t.Skipf("skipping test; this POST call is not supported without actual existing GCP projects and service principal.")
			}
			// anything else is a real failure; returning nil here would make the test pass vacuously
			return err
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
			{
				// reordering the include/exclude set elements must not produce a diff (SSPROD-71536)
				Config:             secureOrgWithAccountIDReordered(accID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func secureOrgWithAccountID(accountID string) string {
	return secureOrgConfig(accountID,
		[]string{"ou-1111", "ou-2222"},
		[]string{"group-a", "group-b"},
		[]string{"group-c", "group-d"},
		[]string{"111111111111", "222222222222"},
		[]string{"333333333333", "444444444444"},
	)
}

func secureOrgWithAccountIDReordered(accountID string) string {
	return secureOrgConfig(accountID,
		[]string{"ou-2222", "ou-1111"},
		[]string{"group-b", "group-a"},
		[]string{"group-d", "group-c"},
		[]string{"222222222222", "111111111111"},
		[]string{"444444444444", "333333333333"},
	)
}

func secureOrgConfig(accountID string, organizationalUnitIDs, includedGroups, excludedGroups, includedAccounts, excludedAccounts []string) string {
	// this is a base64 encoded service account key
	test_service_account_key_encoded := getEncodedGCPServiceAccountKeyForOrg("sample", accountID)
	quote := func(values []string) string { return `"` + strings.Join(values, `", "`) + `"` }

	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "sample" {
  provider_id   = "%s"
  provider_type = "PROVIDER_GCP"
  enabled       = "true"
  feature {
    secure_config_posture {
      enabled    = true
      components = ["COMPONENT_SERVICE_PRINCIPAL/secure-posture"]
    }
    secure_identity_entitlement {
      enabled    = true
      components = ["COMPONENT_WEBHOOK_DATASOURCE/secure-runtime"]
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
    type                       = "COMPONENT_SERVICE_PRINCIPAL"
    instance                   = "secure-onboarding"
    service_principal_metadata = jsonencode({
      gcp = {
        key = "%s"
      }
    })
  }
	component {
		type                        = "COMPONENT_WEBHOOK_DATASOURCE"
		instance                    = "secure-runtime"
		webhook_datasource_metadata = jsonencode({
			gcp = {
				webhook_datasource = {
					pubsub_topic_name      = "pubsub_topic_name_value"
					sink_name              = "sink_name_value"
					push_subscription_name = "push_subscription_name_value"
					push_endpoint          = "push_endpoint_value"
				}
			  service_principal = {
					workload_identity_federation = {
						pool_id          = "pool_id_value"
						pool_provider_id = "pool_provider_id_value"
						project_number   = "123456789011"
					}
					email = "email_value"
				}
			}
		})
	}
}
resource "sysdig_secure_organization" "sample-org" {
  management_account_id		     = sysdig_secure_cloud_auth_account.sample.id
  organization_root_id 		     = "test-id"
  automatic_onboarding           = false
  organizational_unit_ids        = [%s]
  included_organizational_groups = [%s]
  excluded_organizational_groups = [%s]
  included_cloud_accounts        = [%s]
  excluded_cloud_accounts        = [%s]
}
`, accountID, test_service_account_key_encoded, test_service_account_key_encoded,
		quote(organizationalUnitIDs), quote(includedGroups), quote(excludedGroups), quote(includedAccounts), quote(excludedAccounts))
}

func getEncodedGCPServiceAccountKeyForOrg(resourceName string, accountID string) string {
	test_service_account_key_bytes, err := json.Marshal(map[string]any{
		"type":                        "service_account",
		"project_id":                  fmt.Sprintf("%s-%s", resourceName, accountID),
		"private_key_id":              "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		"private_key":                 "-----BEGIN PRIVATE KEY-----\nxxxxxxxxxxxxxxxxxxxxxxxxxxx\n-----END PRIVATE KEY-----\n",
		"client_email":                fmt.Sprintf("some-sa-name@%s-%s.iam.gserviceaccount.com", resourceName, accountID),
		"client_id":                   "some-client-id",
		"auth_uri":                    "https://some-auth-uri",
		"token_uri":                   "https://some-token-uri",
		"auth_provider_x509_cert_url": "https://some-authprovider-cert-url",
		"client_x509_cert_url":        "https://some-client-cert-url",
		"universe_domain":             "googleapis.com",
	})
	if err != nil {
		fmt.Printf("Failed to marshal test_service_account_key: %v", err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, test_service_account_key_bytes, "", "  ")
	if err != nil {
		fmt.Printf("Failed to indent test_service_account_key: %v", err)
	}
	out.WriteByte('\n')

	return b64.StdEncoding.EncodeToString(out.Bytes())
}
