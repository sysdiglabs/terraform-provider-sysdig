//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
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
	// generated once so every step renders the same tenant id, otherwise the reorder step
	// would also diff on provider_tenant_id
	tenantID := acctest.RandStringFromCharSet(36, acctest.CharSetAlphaNum)
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
				Config: secureAzureWithServicePrincipalFeature(accID, tenantID, false),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.azure_sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// reordering the components set elements must not produce a diff
				Config:   secureAzureWithServicePrincipalFeature(accID, tenantID, true),
				PlanOnly: true,
			},
		},
	})
}

func secureAzureWithServicePrincipalFeature(accountID string, randomTenantId string, reorderComponents bool) string {
	// to replicate user behavior, snippet creates an actual azure account and
	// two actual Service Principal components. It then passes cloudauth returned account_id
	// as input to the account feature calls.
	components := `["COMPONENT_SERVICE_PRINCIPAL/secure-posture", "COMPONENT_SERVICE_PRINCIPAL/secure-posture-2"]`
	if reorderComponents {
		components = `["COMPONENT_SERVICE_PRINCIPAL/secure-posture-2", "COMPONENT_SERVICE_PRINCIPAL/secure-posture"]`
	}

	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "azure_sample" {
  provider_id        = "azure-cspm-test-%s"
  provider_type      = "PROVIDER_AZURE"
  enabled            = true
  provider_tenant_id = "%s"
  provider_alias     = "some-alias"
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

resource "sysdig_secure_cloud_auth_account_component" "azure_service_principal_2" {
  account_id		         = sysdig_secure_cloud_auth_account.azure_sample.id
  type                       = "COMPONENT_SERVICE_PRINCIPAL"
  instance                   = "secure-posture-2"
  service_principal_metadata = jsonencode({
	  azure = {
		  active_directory_service_principal = {
				id                        = "some-id-2"
				account_enabled           = true
				display_name              = "some-display-name-2"
				app_display_name          = "some-app-display-name-2"
				app_id                    = "some-app-id-2"
				app_owner_organization_id = "some-app-owner-organization-id-2"
		  }
	  }
  })
}

resource "sysdig_secure_cloud_auth_account_feature" "azure_config_posture" {
  account_id		         = sysdig_secure_cloud_auth_account.azure_sample.id
  type                       = "FEATURE_SECURE_CONFIG_POSTURE"
  enabled                    = true
  components                 = %s

  depends_on = [
	sysdig_secure_cloud_auth_account_component.azure_service_principal,
	sysdig_secure_cloud_auth_account_component.azure_service_principal_2,
  ]
}
`, accountID, randomTenantId, components)
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

func TestAccSecureCloudAuthAccountUpdateKeepsFeatures(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accID := rText()
	tenantID := acctest.RandStringFromCharSet(36, acctest.CharSetAlphaNum)
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
				Config: secureAzureAccountWithFeatureResources(accID, tenantID, "some-alias"),
			},
			{
				// the account PUT must keep the flags and the feature types the account resource does
				// not map, otherwise the feature resources plan to restore them after this apply
				Config: secureAzureAccountWithFeatureResources(accID, tenantID, "renamed-alias"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.azure_sample", "provider_alias", "renamed-alias"),
					// the feature resource ignores a 404 on read, so only the API shows a deleted feature
					testAccCheckCloudAuthAccountFeatureExists("sysdig_secure_cloud_auth_account_feature.azure_workload_scanning_containers"),
				),
			},
		},
	})
}

func secureAzureAccountWithFeatureResources(accountID, tenantID, alias string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "azure_sample" {
  provider_id        = "azure-update-test-%s"
  provider_type      = "PROVIDER_AZURE"
  enabled            = true
  provider_tenant_id = "%s"
  provider_alias     = "%s"
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
  account_id = sysdig_secure_cloud_auth_account.azure_sample.id
  type       = "FEATURE_SECURE_CONFIG_POSTURE"
  enabled    = true
  components = ["COMPONENT_SERVICE_PRINCIPAL/secure-posture"]
  flags      = {
      "POSTURE_VALIDATED_EXPOSURE_ENABLED": "true"
  }

  depends_on = [ sysdig_secure_cloud_auth_account_component.azure_service_principal ]
}

resource "sysdig_secure_cloud_auth_account_feature" "azure_workload_scanning_containers" {
  account_id = sysdig_secure_cloud_auth_account.azure_sample.id
  type       = "FEATURE_SECURE_WORKLOAD_SCANNING_CONTAINERS"
  enabled    = true
  components = ["COMPONENT_SERVICE_PRINCIPAL/secure-posture"]

  depends_on = [ sysdig_secure_cloud_auth_account_feature.azure_config_posture ]
}
`, accountID, tenantID, alias)
}

func testAccCheckCloudAuthAccountFeatureExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s not found in state", resourceName)
		}
		url := os.Getenv("SYSDIG_SECURE_URL")
		if url == "" {
			url = "https://secure.sysdig.com"
		}
		client := v2.NewSysdigSecure(v2.WithURL(url), v2.WithToken(os.Getenv("SYSDIG_SECURE_API_TOKEN")))
		_, errStatus, err := client.GetCloudauthAccountFeatureSecure(context.Background(), rs.Primary.Attributes["account_id"], rs.Primary.Attributes["type"])
		if err != nil {
			return fmt.Errorf("%s: %s %w", resourceName, errStatus, err)
		}
		return nil
	}
}
