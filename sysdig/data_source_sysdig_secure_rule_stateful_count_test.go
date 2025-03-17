//go:build tf_acc_sysdig || tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
					resource.TestCheckResourceAttr("data.sysdig_secure_rule_stateful_count.data_stateful_rule_append", "rule_count", "2"),
				),
			},
		},
	})
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
