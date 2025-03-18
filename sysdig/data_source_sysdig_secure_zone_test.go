package sysdig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccDataSourceSysdigSecureZone(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSysdigSecureZoneConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_zone.test", "name", "test-secure-zone"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_zone.test", "description"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_zone.test", "is_system"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_zone.test", "author"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_zone.test", "last_modified_by"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_zone.test", "last_updated"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"sysdig_secure_zone.test",
						"scope.*",
						map[string]string{
							"target_type": "aws",
							"rules":       "organization in (\"o1\", \"o2\") and account in (\"a1\", \"a2\")",
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceSysdigSecureZoneConfig() string {
	return `
resource "sysdig_secure_zone" "sample" {
  name        = "test-secure-zone"
  description = "Test secure zone"
  scope {
    target_type = "aws"
    rules       = "organization in (\"o1\", \"o2\") and account in (\"a1\", \"a2\")"
  }
}

data "sysdig_secure_zone" "test" {
  depends_on = ["sysdig_secure_zone.sample"]
  name       = sysdig_secure_zone.sample.name
}
	`
}
