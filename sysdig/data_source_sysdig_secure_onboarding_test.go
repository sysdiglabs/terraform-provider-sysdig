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
				Config:      `data "sysdig_secure_trusted_cloud_identity" "trusted_identity" {	cloud_provider = "invalid" }`,
				ExpectError: regexp.MustCompile(`.*expected cloud_provider to be one of.*`),
			},
			{
				Config: `data "sysdig_secure_trusted_cloud_identity" "trusted_identity" {	cloud_provider = "aws" }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_trusted_cloud_identity.trusted_identity", "cloud_provider", "aws"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_cloud_identity.trusted_identity", "aws_account_id"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_cloud_identity.trusted_identity", "aws_role_name"),
					// not asserting the gov exported fields because not every backend environment is gov supported and thus will have empty values
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
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.config_posture", "application_id"),       // uncomment to assert a non empty value
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.config_posture", "tenant_id"),            // uncomment to assert a non empty value
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.config_posture", "service_principal_id"), // uncomment to assert a non empty value
				),
			},
			{
				Config: `data "sysdig_secure_trusted_azure_app" "onboarding" {	name = "onboarding" }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_trusted_azure_app.onboarding", "name", "onboarding"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.onboarding", "application_id"),       // uncomment to assert a non empty value
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.onboarding", "tenant_id"),            // uncomment to assert a non empty value
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.onboarding", "service_principal_id"), // uncomment to assert a non empty value
				),
			},
			{
				Config: `data "sysdig_secure_trusted_azure_app" "threat_detection" { name = "threat_detection" }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_trusted_azure_app.threat_detection", "name", "threat_detection"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.threat_detection", "application_id"),       // uncomment to assert a non empty value
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.threat_detection", "tenant_id"),            // uncomment to assert a non empty value
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.threat_detection", "service_principal_id"), // uncomment to assert a non empty value
				),
			},
			{
				Config: `data "sysdig_secure_trusted_azure_app" "vm_workload_scanning" { name = "vm_workload_scanning" }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_trusted_azure_app.vm_workload_scanning", "name", "vm_workload_scanning"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.vm_workload_scanning", "application_id"),       // uncomment to assert a non empty value
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.vm_workload_scanning", "tenant_id"),            // uncomment to assert a non empty value
					resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_azure_app.vm_workload_scanning", "service_principal_id"), // uncomment to assert a non empty value
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
				Config: `data "sysdig_secure_cloud_ingestion_assets" "assets" {
						cloud_provider = "invalid"
						cloud_provider_id = "123"
						}`,
				ExpectError: regexp.MustCompile(`.*expected cloud_provider to be one of.*`),
			},
			{
				Config: `data "sysdig_secure_cloud_ingestion_assets" "assets" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_cloud_ingestion_assets.assets", "aws.%", "4"),
					// not asserting the gov exported fields because not every backend environme nt is gov supported and thus will have empty values

					resource.TestCheckResourceAttrSet("data.sysdig_secure_cloud_ingestion_assets.assets", "gcp_routing_key"),
					// metadata fields are opaque to api backend; cloudingestion controls what f ields are passed
					// asserts ingestionType and ingestionURL in metadata since it is required
					resource.TestCheckResourceAttr("data.sysdig_secure_cloud_ingestion_assets.assets", "gcp_metadata.ingestionType", "gcp"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_cloud_ingestion_assets.assets", "gcp_metadata.ingestionURL"),
				),
			},
			{
				Config: `data "sysdig_secure_cloud_ingestion_assets" "assets" {
						cloud_provider = "aws" 
						cloud_provider_id = "012345678901"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sysdig_secure_cloud_ingestion_assets.assets", "aws.sns_routing_key"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_cloud_ingestion_assets.assets", "aws.sns_routing_url"),
				),
			},
		},
	})
}

func TestAccTrustedOracleAppDataSource(t *testing.T) {
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
				Config:      `data "sysdig_secure_trusted_oracle_app" "invalid" {	name = "invalid" }`,
				ExpectError: regexp.MustCompile(`.*expected name to be one of.*`),
			},
			{
				Config: `data "sysdig_secure_trusted_oracle_app" "config_posture" {	name = "config_posture" }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_trusted_oracle_app.config_posture", "name", "config_posture"),
					// not asserting the oci exported fields because not every backend environment is oci supported yet and thus will have empty values
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_oracle_app.config_posture", "tenancy_ocid"), // uncomment to assert a non empty value
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_oracle_app.config_posture", "group_ocid"),   // uncomment to assert a non empty value
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_oracle_app.config_posture", "user_ocid"),    // uncomment to assert a non empty value
				),
			},
			{
				Config: `data "sysdig_secure_trusted_oracle_app" "onboarding" {	name = "onboarding" }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_trusted_oracle_app.onboarding", "name", "onboarding"),
					// not asserting the oci exported fields because not every backend environment is oci supported yet and thus will have empty values
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_oracle_app.onboarding", "tenancy_ocid"), // uncomment to assert a non empty value
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_oracle_app.onboarding", "group_ocid"),   // uncomment to assert a non empty value
					// resource.TestCheckResourceAttrSet("data.sysdig_secure_trusted_oracle_app.onboarding", "user_ocid"),    // uncomment to assert a non empty value
				),
			},
		},
	})
}
