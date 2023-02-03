package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccScanningPolicyAssignment(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	rText2 := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

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
				Config: scanningPolicyAssignmentWithoutName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_scanning_policy_assignment.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: scanningPolicyAssignmentWithName(rText2()),
			},
			// {
			// 	Config: scanningPolicyAssignmentWithWhitelistIDs(rText()),
			// },
		},
	})
}

func scanningPolicyAssignmentWithoutName(name string) string {
	return fmt.Sprintf(`
%s
resource "sysdig_secure_scanning_policy_assignment" "sample" {
  items {
    image {
      type = "tag"
      value = "*"
    }
    registry = "*"
    repository = "*"

    policy_ids = [sysdig_secure_scanning_policy.sample.id]
  }
}
`, scanningPolicyWithName(name))
}

func scanningPolicyAssignmentWithName(name string) string {
	return fmt.Sprintf(`
%s
resource "sysdig_secure_scanning_policy_assignment" "sample_2" {
  items {
    name = "example %s"
    image {
      type = "tag"
      value = "latest"
    }
    registry = "icr.io"
    repository = "example"

    policy_ids = [sysdig_secure_scanning_policy.sample.id]
  }

  items {
    name = ""
    image {
      type = "tag"
      value = "*"
    }
    registry = "*"
    repository = "*"

    policy_ids = ["default"]
  }
}
`, scanningPolicyWithName(name), name)
}

func scanningPolicyAssignmentWithWhitelistIDs(name string) string {
	return fmt.Sprintf(`
%s
%s
resource "sysdig_secure_scanning_policy_assignment" "sample_3" {
  items {
    name = "example %s"
    image {
      type = "tag"
      value = "latest"
    }
    registry = "icr.io"
    repository = "example"

    policy_ids = [sysdig_secure_scanning_policy.sample.id]
  }

  items {
    name = ""
    image {
      type = "tag"
      value = "*"
    }
    registry = "*"
    repository = "*"

    policy_ids = ["default"]
	  whitelist_ids = [sysdig_secure_vulnerability_exception_list.sample.id]
  }
}
`, scanningPolicyWithName(name), vulnerabilityException(name), name)
}
