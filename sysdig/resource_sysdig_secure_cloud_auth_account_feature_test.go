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

// TF acceptance tests for secure account feature need an actual azure account
// onboarded to use its account_id (uuid based) as input to the account feature CRUD calls.
// They also need related valid component(s) to be onboarded for account feature to work.

/************
* Azure tests
************/
func TestAccSecureCloudAuthAccountFeature(t *testing.T) {
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
				Config: secureAzureWithServicePrincipalFeature(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.azure_sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureAzureWithServicePrincipalFeature(accountID string) string {
	// to replicate user behavior, snippet creates an actual azure account and
	// an actual Service Principal component. It then passes cloudauth returned account_id
	// as input to the account feature calls.
	rID := func() string { return acctest.RandStringFromCharSet(36, acctest.CharSetAlphaNum) }
	randomTenantId := rID()

	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "azure_sample" {
  provider_id        = "azure-cspm-test-%s"
  provider_type      = "PROVIDER_AZURE"
  enabled            = true
  provider_tenant_id = "%s"
  provider_alias     = "some-alias"
  regulatory_framework = "REGULATORY_FRAMEWORK_UNSPECIFIED"
  lifecycle {
	ignore_changes = [
	  component,
	  feature
	]
  }
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

resource "sysdig_secure_cloud_auth_account_feature" "azure_config_posture" {
  account_id		         = sysdig_secure_cloud_auth_account.azure_sample.id
  type                       = "FEATURE_SECURE_CONFIG_POSTURE"
  enabled                    = true
  components                 = ["COMPONENT_SERVICE_PRINCIPAL/secure-posture"]

  depends_on = [ sysdig_secure_cloud_auth_account_component.azure_service_principal ]
}
`, accountID, randomTenantId)
}

func TestAccSecureCloudAuthAccountFeatureWithFlags(t *testing.T) {
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
				Config: secureAzureWithScanningFeatureWithFlags(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.azure_sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureAzureWithScanningFeatureWithFlags(accountID string) string {
	// to replicate user behavior, snippet creates an actual azure account and
	// an actual Service Principal component. It then passes cloudauth returned account_id
	// as input to the account feature calls.
	rID := func() string { return acctest.RandStringFromCharSet(36, acctest.CharSetAlphaNum) }
	randomTenantId := rID()

	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "azure_sample" {
  provider_id        = "azure-vmscan-test-%s"
  provider_type      = "PROVIDER_AZURE"
  enabled            = true
  provider_tenant_id = "%s"
  provider_alias     = "some-alias"
  regulatory_framework = "REGULATORY_FRAMEWORK_UNSPECIFIED"
  lifecycle {
	ignore_changes = [
	  component,
	  feature
	]
  }
}

resource "sysdig_secure_cloud_auth_account_component" "azure_service_principal" {
  account_id		         = sysdig_secure_cloud_auth_account.azure_sample.id
  type                       = "COMPONENT_SERVICE_PRINCIPAL"
  instance                   = "secure-scanning"
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

resource "sysdig_secure_cloud_auth_account_feature" "azure_agentless_scanning" {
  account_id		         = sysdig_secure_cloud_auth_account.azure_sample.id
  type                       = "FEATURE_SECURE_AGENTLESS_SCANNING"
  enabled                    = true
  components                 = ["COMPONENT_SERVICE_PRINCIPAL/secure-scanning"]
  flags                      = {
      "SCANNING_HOST_CONTAINER_ENABLED": "true"
  }

  depends_on = [ sysdig_secure_cloud_auth_account_component.azure_service_principal ]
}
`, accountID, randomTenantId)
}
