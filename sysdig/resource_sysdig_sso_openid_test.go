//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_onprem_monitor || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSSOOpenID_Basic(t *testing.T) {
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
				Config: ssoOpenIDBasicConfig(integrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test",
						"issuer_url",
						"https://accounts.google.com",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test",
						"client_id",
						"test-client-id",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test",
						"integration_name",
						integrationName,
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test",
						"is_active",
						"true",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test",
						"is_metadata_discovery_enabled",
						"true",
					),
					resource.TestCheckResourceAttrSet(
						"sysdig_sso_openid.test",
						"version",
					),
				),
			},
			{
				ResourceName:            "sysdig_sso_openid.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_secret"},
			},
		},
	})
}

func TestAccSSOOpenID_WithMetadata(t *testing.T) {
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
				Config: ssoOpenIDWithMetadataConfig(integrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test_metadata",
						"is_metadata_discovery_enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test_metadata",
						"metadata.0.issuer",
						"https://idp.example.com",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test_metadata",
						"metadata.0.authorization_endpoint",
						"https://idp.example.com/oauth2/authorize",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test_metadata",
						"metadata.0.token_endpoint",
						"https://idp.example.com/oauth2/token",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test_metadata",
						"metadata.0.jwks_uri",
						"https://idp.example.com/.well-known/jwks.json",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test_metadata",
						"metadata.0.token_auth_method",
						"CLIENT_SECRET_BASIC",
					),
				),
			},
			{
				ResourceName:            "sysdig_sso_openid.test_metadata",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_secret"},
			},
		},
	})
}

func TestAccSSOOpenID_Update(t *testing.T) {
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
				Config: ssoOpenIDBasicConfig(integrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test",
						"integration_name",
						integrationName,
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test",
						"is_group_mapping_enabled",
						"false",
					),
				),
			},
			{
				Config: ssoOpenIDUpdatedConfig(integrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test",
						"integration_name",
						integrationName, // integration_name cannot be updated (ForceNew)
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test",
						"is_group_mapping_enabled",
						"true",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test",
						"group_mapping_attribute_name",
						"custom_groups",
					),
				),
			},
		},
	})
}

func TestAccSSOOpenID_WithAdditionalScopes(t *testing.T) {
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
				Config: ssoOpenIDWithAdditionalScopesConfig(integrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test_scopes",
						"is_additional_scopes_check_enabled",
						"true",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test_scopes",
						"additional_scopes.#",
						"2",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test_scopes",
						"additional_scopes.0",
						"groups",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_openid.test_scopes",
						"additional_scopes.1",
						"roles",
					),
				),
			},
		},
	})
}

func ssoOpenIDBasicConfig(integrationName string) string {
	return fmt.Sprintf(`
resource "sysdig_sso_openid" "test" {
  issuer_url       = "https://accounts.google.com"
  client_id        = "test-client-id"
  client_secret    = "test-client-secret"
  integration_name = "%s"
  is_active        = true
}
`, integrationName)
}

func ssoOpenIDUpdatedConfig(integrationName string) string {
	return fmt.Sprintf(`
resource "sysdig_sso_openid" "test" {
  issuer_url                   = "https://accounts.google.com"
  client_id                    = "test-client-id"
  client_secret                = "test-client-secret"
  integration_name             = "%s"
  is_active                    = true
  is_group_mapping_enabled     = true
  group_mapping_attribute_name = "custom_groups"
  group_attribute_name         = "custom_groups"
}
`, integrationName)
}

func ssoOpenIDWithMetadataConfig(integrationName string) string {
	return fmt.Sprintf(`
resource "sysdig_sso_openid" "test_metadata" {
  issuer_url                     = "https://idp.example.com"
  client_id                      = "test-client-id"
  client_secret                  = "test-client-secret"
  integration_name               = "%s"
  is_metadata_discovery_enabled  = false

  metadata {
    issuer                 = "https://idp.example.com"
    authorization_endpoint = "https://idp.example.com/oauth2/authorize"
    token_endpoint         = "https://idp.example.com/oauth2/token"
    jwks_uri               = "https://idp.example.com/.well-known/jwks.json"
    token_auth_method      = "CLIENT_SECRET_BASIC"
    end_session_endpoint   = "https://idp.example.com/oauth2/logout"
    user_info_endpoint     = "https://idp.example.com/userinfo"
  }
}
`, integrationName)
}

func ssoOpenIDWithAdditionalScopesConfig(integrationName string) string {
	return fmt.Sprintf(`
resource "sysdig_sso_openid" "test_scopes" {
  issuer_url                          = "https://accounts.google.com"
  client_id                           = "test-client-id"
  client_secret                       = "test-client-secret"
  integration_name                    = "%s"
  is_additional_scopes_check_enabled  = true
  additional_scopes                   = ["groups", "roles"]
}
`, integrationName)
}
