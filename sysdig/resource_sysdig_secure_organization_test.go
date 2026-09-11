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

func TestAccSecureOrganization(t *testing.T) {
	// The organization API scrapes every folder and project under the organization with the service
	// principal from the config, so it only succeeds against a real GCP organization. Gate the test
	// on that fixture: recognising the resulting failure by its message cannot be told apart from a
	// genuine outage against the same endpoint.
	if os.Getenv("SYSDIG_SECURE_GCP_ORG_TESTS") == "" {
		t.Skip("Skipping tests on sysdig_secure_organization resource because SYSDIG_SECURE_GCP_ORG_TESTS is not set")
	}

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
				Config: secureOrgWithAccountID(accID),
			},
			{
				ResourceName:            "sysdig_secure_cloud_auth_account.sample",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"component"},
			},
			{
				// reordering the include/exclude set elements must not produce a diff
				Config:   secureOrgWithAccountIDReordered(accID),
				PlanOnly: true,
			},
		},
	})
}

// the values have to satisfy cloudauth's per-provider validation: organizational groups are
// `folders/<id>` and cloud accounts are GCP project ids. organizational_unit_ids is left unset on
// purpose, the API rejects it when the include/exclude fields are used.
var (
	secureOrgIncludedGroups   = []string{"folders/111111111111", "folders/222222222222"}
	secureOrgExcludedGroups   = []string{"folders/333333333333", "folders/444444444444"}
	secureOrgIncludedAccounts = []string{"sample-project-one", "sample-project-two"}
	secureOrgExcludedAccounts = []string{"sample-project-three", "sample-project-four"}
)

func secureOrgWithAccountID(accountID string) string {
	return secureOrgConfig(accountID,
		secureOrgIncludedGroups, secureOrgExcludedGroups,
		secureOrgIncludedAccounts, secureOrgExcludedAccounts,
	)
}

// derived from the same slices so the two configs cannot drift apart and stop exercising a reorder
func secureOrgWithAccountIDReordered(accountID string) string {
	return secureOrgConfig(accountID,
		reversed(secureOrgIncludedGroups), reversed(secureOrgExcludedGroups),
		reversed(secureOrgIncludedAccounts), reversed(secureOrgExcludedAccounts),
	)
}

func reversed(values []string) []string {
	out := make([]string, 0, len(values))
	for i := len(values) - 1; i >= 0; i-- {
		out = append(out, values[i])
	}
	return out
}

func secureOrgConfig(accountID string, includedGroups, excludedGroups, includedAccounts, excludedAccounts []string) string {
	// this is a base64 encoded service account key
	test_service_account_key_encoded := getEncodedGCPServiceAccountKeyForOrg("sample", accountID)

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
  organization_root_id 		     = "organizations/123456789012"
  automatic_onboarding           = false
  included_organizational_groups = [%s]
  excluded_organizational_groups = [%s]
  included_cloud_accounts        = [%s]
  excluded_cloud_accounts        = [%s]
}
`, accountID, test_service_account_key_encoded, test_service_account_key_encoded,
		quoteJoin(includedGroups), quoteJoin(excludedGroups), quoteJoin(includedAccounts), quoteJoin(excludedAccounts))
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
