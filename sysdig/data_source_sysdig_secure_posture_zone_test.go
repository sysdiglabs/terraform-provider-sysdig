//go:build tf_acc_sysdig_secure || tf_acc_ibm_secure || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceSysdigSecurePostureZones(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSysdigSecurePostureZonesWithMultipleResourcesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSysdigSecurePostureZonesExists("data.sysdig_secure_posture_zone.test_posture_zone"),
					resource.TestCheckResourceAttr("data.sysdig_secure_posture_zone.test_posture_zone", "name", "test-zone-1"),
					resource.TestCheckResourceAttr("data.sysdig_secure_posture_zone.test_posture_zone", "description", "Test description 1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"data.sysdig_secure_posture_zone.test_posture_zone",
						"scopes.*",
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

func testAccDataSourceSysdigSecurePostureZonesWithMultipleResourcesConfig() string {
	return `
	resource "sysdig_secure_posture_zone" "test_posture_zone" {
		name        = "test-zone-1"
		description = "Test description 1"

		scopes {
    			scope {
    			  target_type = "aws"
                  rules       = "organization in (\"o1\", \"o2\") and account in (\"a1\", \"a2\")"
    			}
  			}
	}

	data "sysdig_secure_posture_zone" "test_posture_zone" {
		id = sysdig_secure_posture_zone.test_posture_zone.id
	}
	`
}

func testAccCheckDataSourceSysdigSecurePostureZonesExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}
