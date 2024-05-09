//go:build tf_acc_sysdig_secure || tf_acc_ibm_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCreatePosturePolicy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: minimalSecurePosturePolicy(randomText(10)),
			},
			{
				Config: securePosturePolicyWithGroups(randomText(10)),
			},
			{
				Config: securePosturePolicyWithGroupsReqsAndConbrols(randomText(10)),
			},
			{
				ResourceName:      "sysdig_secure_posture_policy.p1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func minimalSecurePosturePolicy(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_posture_policy" "p1" {
  name = "%s"
  description = "new description"
}`, name)
}

func securePosturePolicyWithGroups(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_posture_policy" "z1" {
  name = "%s"
  description = "new description"
  group {
	name = "group1"
	description = "new description"
  }
}`, name)
}

func securePosturePolicyWithGroupsReqsAndConbrols(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_posture_policy" "z1" {
  name = "%s"
  description = "new description"
  group {
	name = "g1"
	description = "new group"
	requirement {
	  name = "r1"
	  description = "r1"
	  control {
		name = "b"
	  }
	}
	group {
	  name = "g2"
	  description = "new group"
			requirement {
	  name = "r2"
	  description = "r1"
	  control {
		name = "b"
	  }
	}
	  group {
		name = "g3"
		description = "new group"
			  requirement {
	  name = "r3"
	  description = "r1"
	  control {
		name = "b"
	  }
	}
		group {
		  name = "g4"
		  description = "new group"
				requirement {
	  name = "r4"
	  description = "r1"
	  control {
		name = "b"
	  }
	}
		  group {
			name = "g5"
		  description = "new group"
				requirement {
	  name = "r5"
	  description = "r1"
	  control {
		name = "b"
	  }
	}
		  }
		}
	  }
	}
  }
}`, name)
}
