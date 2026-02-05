//go:build tf_acc_sysdig || tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccRuleStatefulCountDataSource(t *testing.T) {
	if strings.HasSuffix(os.Getenv("SYSDIG_SECURE_URL"), "ibm.com") {
		t.Skip("Skipping stateful tests for IBM Cloud")
		return
	}

	resource.Test(t, resource.TestCase{
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
				Config: ruleStatefulCountDataSource(),
				Check: resource.ComposeTestCheckFunc(
					testCheckRuleCountAtLeast("data.sysdig_secure_rule_stateful_count.data_stateful_rule_append", 2),
				),
			},
		},
	})
}

func testCheckRuleCountAtLeast(resourceName string, minCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		countStr := rs.Primary.Attributes["rule_count"]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("rule_count is not a valid integer: %s", countStr)
		}
		if count < minCount {
			return fmt.Errorf("rule_count expected >= %d, got %d", minCount, count)
		}
		return nil
	}
}

func ruleStatefulCountDataSource() string {
	return fmt.Sprintf(`
%s

data "sysdig_secure_rule_stateful_count" "data_stateful_rule_append" {
  name = "API Gateway Enumeration Detected"
  source = "awscloudtrail_stateful"
  depends_on = [ sysdig_secure_rule_stateful.stateful_rule_append ]
}
`, ruleStatefulAppend())
}
