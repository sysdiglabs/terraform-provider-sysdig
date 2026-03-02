//go:build tf_acc_onprem_monitor || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSSOSamlOnprem_WithMetadataURL(t *testing.T) {
	integrationName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			monitor := os.Getenv("SYSDIG_MONITOR_API_TOKEN")
			secure := os.Getenv("SYSDIG_SECURE_API_TOKEN")
			if monitor == "" && secure == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN or SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: ssoSamlOnpremWithMetadataURLConfig(integrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test",
						"is_system",
						"true",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test",
						"metadata_url",
						"https://idp.example.com/metadata",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test",
						"email_parameter",
						"email",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test",
						"integration_name",
						integrationName,
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test",
						"is_active",
						"true",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test",
						"is_signature_validation_enabled",
						"true",
					),
					resource.TestCheckResourceAttrSet(
						"sysdig_sso_saml.test",
						"version",
					),
				),
			},
			{
				ResourceName:      "sysdig_sso_saml.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return "system/" + s.RootModule().Resources["sysdig_sso_saml.test"].Primary.ID, nil
				},
			},
		},
	})
}

func TestAccSSOSamlOnprem_WithMetadataXML(t *testing.T) {
	integrationName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			monitor := os.Getenv("SYSDIG_MONITOR_API_TOKEN")
			secure := os.Getenv("SYSDIG_SECURE_API_TOKEN")
			if monitor == "" && secure == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN or SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: ssoSamlOnpremWithMetadataXMLConfig(integrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test_xml",
						"is_system",
						"true",
					),
					resource.TestCheckResourceAttrSet(
						"sysdig_sso_saml.test_xml",
						"metadata_xml",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test_xml",
						"email_parameter",
						"email",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test_xml",
						"integration_name",
						integrationName,
					),
				),
			},
			{
				ResourceName:      "sysdig_sso_saml.test_xml",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return "system/" + s.RootModule().Resources["sysdig_sso_saml.test_xml"].Primary.ID, nil
				},
			},
		},
	})
}

func TestAccSSOSamlOnprem_Update(t *testing.T) {
	integrationName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			monitor := os.Getenv("SYSDIG_MONITOR_API_TOKEN")
			secure := os.Getenv("SYSDIG_SECURE_API_TOKEN")
			if monitor == "" && secure == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN or SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: ssoSamlOnpremWithMetadataURLConfig(integrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test",
						"integration_name",
						integrationName,
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test",
						"is_group_mapping_enabled",
						"false",
					),
				),
			},
			{
				Config: ssoSamlOnpremUpdatedConfig(integrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test",
						"integration_name",
						fmt.Sprintf("%s-updated", integrationName),
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test",
						"is_group_mapping_enabled",
						"true",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test",
						"group_mapping_attribute_name",
						"custom_groups",
					),
				),
			},
		},
	})
}

func TestAccSSOSamlOnprem_SecuritySettings(t *testing.T) {
	integrationName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			monitor := os.Getenv("SYSDIG_MONITOR_API_TOKEN")
			secure := os.Getenv("SYSDIG_SECURE_API_TOKEN")
			if monitor == "" && secure == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN or SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: ssoSamlOnpremWithSecuritySettingsConfig(integrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test_security",
						"is_system",
						"true",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test_security",
						"is_signature_validation_enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test_security",
						"is_signed_assertion_enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test_security",
						"is_destination_verification_enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_saml.test_security",
						"is_encryption_support_enabled",
						"true",
					),
				),
			},
		},
	})
}

func ssoSamlOnpremWithMetadataURLConfig(integrationName string) string {
	return fmt.Sprintf(`
resource "sysdig_sso_saml" "test" {
  is_system        = true
  metadata_url     = "https://idp.example.com/metadata"
  email_parameter  = "email"
  integration_name = "%s"
  is_active        = true
}
`, integrationName)
}

func ssoSamlOnpremWithMetadataXMLConfig(integrationName string) string {
	return fmt.Sprintf(`
resource "sysdig_sso_saml" "test_xml" {
  is_system        = true
  metadata_xml     = <<-EOF
<?xml version="1.0" encoding="UTF-8"?>
<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://idp.example.com">
  <IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://idp.example.com/sso"/>
  </IDPSSODescriptor>
</EntityDescriptor>
EOF
  email_parameter  = "email"
  integration_name = "%s"
  is_active        = true
}
`, integrationName)
}

func ssoSamlOnpremUpdatedConfig(integrationName string) string {
	return fmt.Sprintf(`
resource "sysdig_sso_saml" "test" {
  is_system                     = true
  metadata_url                  = "https://idp.example.com/metadata"
  email_parameter               = "email"
  integration_name              = "%s-updated"
  is_active                     = true
  is_group_mapping_enabled      = true
  group_mapping_attribute_name  = "custom_groups"
}
`, integrationName)
}

func ssoSamlOnpremWithSecuritySettingsConfig(integrationName string) string {
	return fmt.Sprintf(`
resource "sysdig_sso_saml" "test_security" {
  is_system                           = true
  metadata_url                        = "https://idp.example.com/metadata"
  email_parameter                     = "email"
  integration_name                    = "%s"
  is_active                           = true
  is_signature_validation_enabled     = false
  is_signed_assertion_enabled         = false
  is_destination_verification_enabled = false
  is_encryption_support_enabled       = true
}
`, integrationName)
}
