//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common

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

/************
* Azure tests
************/
func TestAccSecureCloudAuthAccountComponent(t *testing.T) {
	// TF acceptance tests for secure account component need an actual azure account
	// onboarded to use its account_id (uuid based) as input to the account component CRUD calls.
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
				Config: secureAzureWithServicePrincipalComponent(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account_component.azure_service_principal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureAzureWithServicePrincipalComponent(accountID string) string {
	// to replicate user behavior, snippet creates an actual azure account and passes
	// cloudauth returned account_id as input to the account component calls.
	rID := func() string { return acctest.RandStringFromCharSet(36, acctest.CharSetAlphaNum) }
	randomTenantId := rID()

	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "azure_sample" {
  provider_id        = "azure-sp-test-%s"
  provider_type      = "PROVIDER_AZURE"
  enabled            = true
  provider_tenant_id = "%s"
  provider_alias     = "some-alias"
}

resource "sysdig_secure_cloud_auth_account_component" "azure_service_principal" {
  account_id		         = sysdig_secure_cloud_auth_account.azure_sample.id
  type                       = "COMPONENT_SERVICE_PRINCIPAL"
  instance                   = "secure-posture"
  service_principal_metadata = jsonencode({
	  azure = {
		  active_directory_service_principal = {
				id                        = "some-id"
				account_enabled           = true
				display_name              = "some-display-name"
				app_display_name          = "some-app-display-name"
				app_id                    = "some-app-id"
				app_owner_organization_id = "some-app-owner-organization-id"
		  }
	  }
  })
}
`, accountID, randomTenantId)
}
