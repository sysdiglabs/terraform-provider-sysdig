package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureCloudAccount(t *testing.T) {
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
				Config: secureCloudAccountWithID(accID),
			},
			{
				Config: secureCloudAccountMinimumConfiguration(accID),
			},
			{
				Config: secureCloudAccountWithWID(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_account.sample-1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "sysdig_secure_cloud_account.sample-2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "sysdig_secure_cloud_account.sample-3",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureCloudAccountWithID(accountID string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_account" "sample-1" {
  account_id          = "sample1-%s"
  cloud_provider      = "aws"
  alias               = "%s"
  role_enabled        = "false"
  role_name            = "CustomRoleName"
}
`, accountID, accountID)
}

func secureCloudAccountMinimumConfiguration(accountID string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_account" "sample-2" {
  account_id      = "sample2-%s"
  cloud_provider  = "aws"
}`, accountID)
}

func secureCloudAccountWithWID(accountID string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_account" "sample-3" {
  account_id          = "sample3-%s"
  cloud_provider      = "gcp"
  alias               = "%s"
  role_enabled        = "false"
  role_name            = "CustomRoleName"
  workload_identity_account_id = "sample3-%s"
}
`, accountID, accountID, accountID)
}
