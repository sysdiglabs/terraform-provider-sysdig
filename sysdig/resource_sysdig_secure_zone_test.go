//go:build tf_acc_sysdig_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSysdigZone_basic(t *testing.T) {
	zoneName := "Zone_TF_" + randomText(5)
	zoneDescription := "Test Zone Description"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: zoneConfig(zoneName, zoneDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_secure_zone.test", "name", zoneName),
					resource.TestCheckResourceAttr("sysdig_secure_zone.test", "description", zoneDescription),
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
			{
				ResourceName:      "sysdig_secure_zone.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: zoneConfig(zoneName, "Updated Description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_secure_zone.test", "description", "Updated Description"),
				),
			},
		},
	})
}

func zoneConfig(name, description string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_zone" "test" {
  name        = "%s"
  description = "%s"
  scope {
    target_type = "aws"
    rules       = "organization in (\"o1\", \"o2\") and account in (\"a1\", \"a2\")"
  }
}
`, name, description)
}
