//go:build tf_acc_sysdig_secure || tf_acc_onprem_secure || tf_acc_ibm_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccDataSourceSysdigSecureZone(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
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
						"data.sysdig_secure_zone.test",
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

func TestAccDataSourceSysdigSecureZone_ByName(t *testing.T) {
	zoneName := "Zone_DS_" + randomText(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureZoneByName(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.sysdig_secure_zone.test",
						"name",
						zoneName,
					),
					resource.TestCheckResourceAttr(
						"data.sysdig_secure_zone.test",
						"scope.0.target_type",
						"aws",
					),

					// v2 expressions
					resource.TestCheckResourceAttr(
						"data.sysdig_secure_zone.test",
						"scope.0.expression.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"data.sysdig_secure_zone.test",
						"scope.0.expression.0.field",
						"organization",
					),
				),
			},
		},
	})
}

func TestAccDataSourceSysdigSecureZone_ByID(t *testing.T) {
	zoneName := "Zone_DS_ID_" + randomText(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecureZoneByID(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.sysdig_secure_zone.test",
						"id",
					),
					resource.TestCheckResourceAttr(
						"data.sysdig_secure_zone.test",
						"scope.0.expression.#",
						"1",
					),
				),
			},
		},
	})
}

func testAccDataSourceSecureZoneByName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_zone" "test" {
  name        = "%s"
  description = "ds acceptance test"

  scope {
    target_type = "aws"

    expression {
      field    = "organization"
      operator = "in"
      values   = ["o1", "o2"]
    }
  }
}

data "sysdig_secure_zone" "test" {
  depends_on = [sysdig_secure_zone.test]
  name = "%s"
}
`, name, name)
}

func testAccDataSourceSecureZoneByID(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_zone" "test" {
  name        = "%s"
  description = "ds acceptance test"

  scope {
    target_type = "aws"

    expression {
      field    = "organization"
      operator = "in"
      values   = ["o1", "o2"]
    }
  }
}

data "sysdig_secure_zone" "test" {
  depends_on = [sysdig_secure_zone.test]
  id = sysdig_secure_zone.test.id
}
`, name)
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
  name       = "test-secure-zone"
}
	`
}
