//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureOrganization(t *testing.T) {
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
				ResourceName:      "sysdig_secure_organization.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureOrgWithAccountID(accountID string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "sample" {
  provider_id   = "sample-%s"
  provider_type = "PROVIDER_GCP"
  enabled       = "true"
}
resource "sysdig_secure_organization" "sample" {
  management_account_id		= sysdig_secure_cloud_auth_account.sample.id
}
`, accountID)
}
