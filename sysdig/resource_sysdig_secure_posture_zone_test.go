package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSecurePostureZone(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: minimalSecurePostureZone(randomText(10)),
			},
			{
				Config: securePostureZoneWithScopes(randomText(10)),
			},
			{
				Config: securePostureZoneWithPolicies(randomText(10)),
			},
			{
				ResourceName:      "sysdig_secure_posture_zone.z1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func minimalSecurePostureZone(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_posture_zone" "z1" {
  name = "%s"
}`, name)
}

func securePostureZoneWithScopes(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_posture_zone" "z1" {
  name = "%s"
  scopes {
    scope {
      target_type = "aws"
      rules       = "organization in (\"o1\", \"o2\") and account in (\"a1\", \"a2\")"
    }

    scope {
      target_type = "azure"
      rules       = "organization contains \"o1\""
    }
  }
}`, name)
}

func securePostureZoneWithPolicies(name string) string {
	return fmt.Sprintf(`
data "sysdig_secure_posture_policies" "all" {}

resource "sysdig_secure_posture_zone" "z1" {
  name = "%s"
  policies = [data.sysdig_secure_posture_policies.all.policies[0].id]
}`, name)
}
