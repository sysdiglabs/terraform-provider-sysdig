//go:build tf_acc_sysdig_secure || tf_acc_ibm_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSecureZonePosturePolicyAssignment_basic(t *testing.T) {
	zoneName := "ZonePolicyAssign_TF_" + randomText(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			// Step 1: Create zone + assignment with 1 policy
			{
				Config: testAccZonePolicyAssignmentWith1Policy(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("sysdig_secure_zone_posture_policy_assignment.test", "zone_id"),
					resource.TestCheckResourceAttr("sysdig_secure_zone_posture_policy_assignment.test", "policy_ids.#", "1"),
				),
			},
			// Step 2: Update to 2 policies
			{
				Config: testAccZonePolicyAssignmentWith2Policies(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_secure_zone_posture_policy_assignment.test", "policy_ids.#", "2"),
				),
			},
			// Step 3: Import by zone_id
			{
				ResourceName:      "sysdig_secure_zone_posture_policy_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccZonePolicyAssignmentWith1Policy(zoneName string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_zone" "test" {
  name = "%s"
  scope {
    target_type = "aws"
    rules       = "account in (\"111111111111\")"
  }
}

data "sysdig_secure_posture_policy" "p1" {
  name = "Sysdig Kubernetes"
}

resource "sysdig_secure_zone_posture_policy_assignment" "test" {
  zone_id    = sysdig_secure_zone.test.id
  policy_ids = [data.sysdig_secure_posture_policy.p1.id]
}
`, zoneName)
}

func testAccZonePolicyAssignmentWith2Policies(zoneName string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_zone" "test" {
  name = "%s"
  scope {
    target_type = "aws"
    rules       = "account in (\"111111111111\")"
  }
}

data "sysdig_secure_posture_policy" "p1" {
  name = "Sysdig Kubernetes"
}

data "sysdig_secure_posture_policy" "p2" {
  id = 1
}

resource "sysdig_secure_zone_posture_policy_assignment" "test" {
  zone_id    = sysdig_secure_zone.test.id
  policy_ids = [
    data.sysdig_secure_posture_policy.p1.id,
    data.sysdig_secure_posture_policy.p2.id,
  ]
}
`, zoneName)
}
