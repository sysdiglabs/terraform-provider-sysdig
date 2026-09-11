//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

const (
	gcpOrgServiceAccountKeyEnv = "SYSDIG_SECURE_GCP_ORG_SERVICE_ACCOUNT_KEY"
	gcpOrgRootIDEnv            = "SYSDIG_SECURE_GCP_ORG_ROOT_ID"
	gcpOrgProjectIDEnv         = "SYSDIG_SECURE_GCP_ORG_PROJECT_ID"
)

type gcpOrgFixture struct {
	projectID          string
	organizationRootID string
	serviceAccountKey  string
}

func TestAccSecureOrganization(t *testing.T) {
	// Creating the organization scrapes every folder and project under it with this service
	// principal, so the key, the organization root and the management project all have to be real.
	fixture := gcpOrgFixture{
		projectID:          os.Getenv(gcpOrgProjectIDEnv),
		organizationRootID: os.Getenv(gcpOrgRootIDEnv),
		serviceAccountKey:  os.Getenv(gcpOrgServiceAccountKeyEnv),
	}
	if fixture.projectID == "" || fixture.organizationRootID == "" || fixture.serviceAccountKey == "" {
		t.Skipf("Skipping tests on sysdig_secure_organization resource because %s, %s and %s are not all set",
			gcpOrgProjectIDEnv, gcpOrgRootIDEnv, gcpOrgServiceAccountKeyEnv)
	}

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
				Config: secureOrgConfigInOrder(fixture),
			},
			{
				ResourceName:            "sysdig_secure_cloud_auth_account.sample",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"component"},
			},
			{
				// reordering the include/exclude set elements must not produce a diff
				Config:   secureOrgConfigReordered(fixture),
				PlanOnly: true,
			},
		},
	})
}

// The include/exclude values stay synthetic: cloudauth validates their shape, not their
// existence, so well-formed folder ids and project ids are enough and keep the reorder under the
// test's control. organizational_unit_ids is left unset, the API rejects it alongside these.
var (
	secureOrgIncludedGroups   = []string{"folders/111111111111", "folders/222222222222"}
	secureOrgExcludedGroups   = []string{"folders/333333333333", "folders/444444444444"}
	secureOrgIncludedAccounts = []string{"sample-project-one", "sample-project-two"}
	secureOrgExcludedAccounts = []string{"sample-project-three", "sample-project-four"}
)

func secureOrgConfigInOrder(fixture gcpOrgFixture) string {
	return secureOrgConfig(fixture,
		secureOrgIncludedGroups, secureOrgExcludedGroups,
		secureOrgIncludedAccounts, secureOrgExcludedAccounts,
	)
}

// derived from the same slices so the two configs cannot drift apart and stop exercising a reorder
func secureOrgConfigReordered(fixture gcpOrgFixture) string {
	return secureOrgConfig(fixture,
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

func secureOrgConfig(fixture gcpOrgFixture, includedGroups, excludedGroups, includedAccounts, excludedAccounts []string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "sample" {
  provider_id   = %q
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
        key = %q
      }
    })
  }
	component {
    type                       = "COMPONENT_SERVICE_PRINCIPAL"
    instance                   = "secure-onboarding"
    service_principal_metadata = jsonencode({
      gcp = {
        key = %q
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
  organization_root_id 		     = %q
  automatic_onboarding           = false
  included_organizational_groups = [%s]
  excluded_organizational_groups = [%s]
  included_cloud_accounts        = [%s]
  excluded_cloud_accounts        = [%s]
}
`, fixture.projectID, fixture.serviceAccountKey, fixture.serviceAccountKey, fixture.organizationRootID,
		quoteJoin(includedGroups), quoteJoin(excludedGroups), quoteJoin(includedAccounts), quoteJoin(excludedAccounts))
}
