//go:build tf_acc_sysdig_secure

package sysdig_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccTrustedCloudIdentityDataSource(t *testing.T) {
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
				Config: `data "sysdig_secure_trusted_cloud_identity" "trusted_identity" {	cloud_provider = "aws" }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_trusted_cloud_identity.trusted_identity", "cloud_provider", "aws"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_cloud_identity.trusted_identity", "aws_account_id"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_cloud_identity.trusted_identity", "aws_role_name"),
				),
			},
			{
				Config: `data "sysdig_secure_trusted_cloud_identity" "trusted_identity" {	cloud_provider = "gcp" }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_trusted_cloud_identity.trusted_identity", "cloud_provider", "gcp"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_cloud_identity.trusted_identity", "aws_account_id"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_cloud_identity.trusted_identity", "aws_role_name"),
				),
			},
			{
				Config: `data "sysdig_secure_trusted_cloud_identity" "trusted_identity" { cloud_provider = "azure" }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_trusted_cloud_identity.trusted_identity", "cloud_provider", "azure"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_cloud_identity.trusted_identity", "azure_tenant_id"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_cloud_identity.trusted_identity", "azure_service_principal_id"),
				),
			},
		},
	})
}

func TestAccTrustedAzureAppDataSource(t *testing.T) {
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
				Config:      `data "sysdig_secure_trusted_azure_app" "config_posture" {	name = "invalid" }`,
				ExpectError: regexp.MustCompile(`.*expected name to be one of.*`),
			},
			{
				Config: `data "sysdig_secure_trusted_azure_app" "config_posture" {	name = "config_posture" }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_trusted_azure_app.config_posture", "name", "config_posture"),
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.config_posture", "application_id"), // uncomment to assert a non empty value
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.config_posture", "tenant_id"), // uncomment to assert a non empty value
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.config_posture", "service_principal_id"), // uncomment to assert a non empty value
				),
			},
			{
				Config: `data "sysdig_secure_trusted_azure_app" "onboarding" {	name = "onboarding" }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_trusted_azure_app.onboarding", "name", "onboarding"),
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.onboarding", "application_id"), // uncomment to assert a non empty value
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.onboarding", "tenant_id"), // uncomment to assert a non empty value
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.onboarding", "service_principal_id"), // uncomment to assert a non empty value
				),
			},
			{
				Config: `data "sysdig_secure_trusted_azure_app" "threat_detection" { name = "threat_detection" }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_trusted_azure_app.threat_detection", "name", "threat_detection"),
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.threat_detection", "application_id"), // uncomment to assert a non empty value
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.threat_detection", "tenant_id"), // uncomment to assert a non empty value
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.threat_detection", "service_principal_id"), // uncomment to assert a non empty value
				),
			},
		},
	})
}

func TestAccTenantExternalIDDataSource(t *testing.T) {
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
				Config: `data "sysdig_secure_tenant_external_id" "external_id" {}`,
			},
		},
	})
}

func TestAccAgentlessScanningAssetsDataSource(t *testing.T) {
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
				Config: `data "sysdig_secure_agentless_scanning_assets" "assets" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_agentless_scanning_assets.assets", "aws.%", "1"),
					resource.TestCheckResourceAttr("data.sysdig_secure_agentless_scanning_assets.assets", "azure.%", "2"),
					resource.TestCheckResourceAttr("data.sysdig_secure_agentless_scanning_assets.assets", "backend.%", "2"),
					resource.TestCheckResourceAttr("data.sysdig_secure_agentless_scanning_assets.assets", "gcp.%", "1"),
				),
			},
		},
	})
}

func TestAccCloudIngestionAssetsDataSource(t *testing.T) {
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
				Config: `data "sysdig_secure_cloud_ingestion_assets" "assets" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_cloud_ingestion_assets.assets", "aws.%", "1"),
					resource.TestCheckResourceAttr("data.sysdig_secure_cloud_ingestion_assets.assets", "gcp.%", "2"),
				),
			},
		},
	})
}
