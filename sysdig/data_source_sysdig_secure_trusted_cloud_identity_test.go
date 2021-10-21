package sysdig_test

import (
	"os"
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
				Config: trustedIdentityDatasourceAWS(),
			},
			{
				Config: trustedIdentityDatasourceGCP(),
			},
			{
				Config: trustedIdentityDatasourceAzure(),
			},
		},
	})
}

func trustedIdentityDatasourceAWS() string {
	return `
data "sysdig_secure_trusted_cloud_identity" "trusted_identity" {
	cloud_provider = "aws"
}
`
}

func trustedIdentityDatasourceGCP() string {
	return `
data "sysdig_secure_trusted_cloud_identity" "trusted_identity" {
	cloud_provider = "gcp"
}
`
}

func trustedIdentityDatasourceAzure() string {
	return `
data "sysdig_secure_trusted_cloud_identity" "trusted_identity" {
	cloud_provider = "azure"
}
`
}
