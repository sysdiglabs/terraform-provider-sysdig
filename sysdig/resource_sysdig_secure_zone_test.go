//go:build tf_acc_sysdig_secure || tf_acc_onprem_secure || tf_acc_ibm_secure

package sysdig_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ImportStateVerify is intentionally disabled here because
// legacy rules are normalized into expression blocks during Read.
// Structural equality is not preserved, but semantic equivalence is.
func TestAccSysdigZone_basic(t *testing.T) {
	zoneName := "Zone_TF_" + randomText(5)
	zoneDescription := "Test Zone Description"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
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
				ResourceName: "sysdig_secure_zone.test",
				ImportState:  true,
			},
			{
				Config:   zoneConfig(zoneName, zoneDescription),
				PlanOnly: true,
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

func TestAccSysdigSecureZone_LegacyRules(t *testing.T) {
	resourceName := "sysdig_secure_zone.legacy"

	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSecureZoneLegacy(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc-legacy"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.target_type", "kubernetes"),
				),
			},
			{
				// refresh only
				PlanOnly: true,
				Config:   testAccSecureZoneLegacy(),
			},
		},
	})
}

func TestAccSysdigSecureZone_ExpressionOnly(t *testing.T) {
	resourceName := "sysdig_secure_zone.expr"

	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSecureZoneExpression(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc-expr"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.target_type", "kubernetes"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.expression.#", "2"),
					// In SDK v2, optional attributes in nested TypeSet elements are always
					// materialized in state (as empty string). We verify rules is empty, not absent.
					resource.TestCheckResourceAttr(resourceName, "scope.0.rules", ""),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccSecureZoneExpression(),
			},
		},
	})
}

func TestAccSysdigSecureZone_MigrateRulesToExpression(t *testing.T) {
	resourceName := "sysdig_secure_zone.migrate"

	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSecureZoneLegacyMigration(),
			},
			{
				Config: testAccSecureZoneExpressionMigration(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "migrated"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.expression.#", "2"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccSecureZoneExpressionMigration(),
			},
		},
	})
}

func TestAccSysdigSecureZone_V2RulesOnly(t *testing.T) {
	resourceName := "sysdig_secure_zone.v2rules"

	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSecureZoneV2Rules(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc-v2rules"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.target_type", "kubernetes"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.expression.#", "0"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccSecureZoneV2Rules(),
			},
		},
	})
}

func TestAccSysdigSecureZone_InvalidRulesAndExpression(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccSecureZoneInvalid(),
				ExpectError: regexp.MustCompile("cannot be used together with"),
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

func testAccSecureZoneLegacy() string {
	return `
resource "sysdig_secure_zone" "legacy" {
  name        = "acc-legacy"
  description = "legacy rules"

  scope {
    target_type = "kubernetes"
    rules = "agentTags != \"key: value\" and not agentTags contains \"key2: value2\""
  }
}
`
}

func testAccSecureZoneExpression() string {
	return `
resource "sysdig_secure_zone" "expr" {
  name        = "acc-expr"
  description = "expression test"

  scope {
    target_type = "kubernetes"

    expression {
      field    = "agent.tag.key"
      operator = "is_not"
      value    = "value"
    }

    expression {
      field    = "agent.tag.key2"
      operator = "not_contains"
      value    = "value2"
    }
  }
}
`
}

func testAccSecureZoneLegacyMigration() string {
	return `
resource "sysdig_secure_zone" "migrate" {
  name        = "acc-migrate"
  description = "legacy"

  scope {
    target_type = "kubernetes"
    rules = "agentTags != \"key: value\" and not agentTags contains \"key2: value2\""
  }
}
`
}

func testAccSecureZoneExpressionMigration() string {
	return `
resource "sysdig_secure_zone" "migrate" {
  name        = "acc-migrate"
  description = "migrated"

  scope {
    target_type = "kubernetes"

    expression {
      field    = "agent.tag.key"
      operator = "is_not"
      value    = "value"
    }

    expression {
      field    = "agent.tag.key2"
      operator = "not_contains"
      value    = "value2"
    }
  }
}
`
}

func testAccSecureZoneV2Rules() string {
	return `
resource "sysdig_secure_zone" "v2rules" {
  name        = "acc-v2rules"
  description = "v2 rules test"

  scope {
    target_type = "kubernetes"
    rules = "agent.tag.key != \"value\" and not agent.tag.key2 contains \"value2\""
  }
}
`
}

func testAccSecureZoneInvalid() string {
	return `
resource "sysdig_secure_zone" "invalid" {
  name = "acc-invalid"

  scope {
    target_type = "kubernetes"
    rules = "agentTags != \"key: value\""

    expression {
      field    = "agent.tag.key"
      operator = "is_not"
      value    = "value"
    }
  }
}
`
}
