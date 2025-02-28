//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

package sysdig_test

import (
	"os"
	"strings"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestRuleStatefulAppends(t *testing.T) {
	if strings.HasSuffix(os.Getenv("SYSDIG_SECURE_URL"), "ibm.com") {
		t.Skip("Skipping stateful tests for IBM Cloud")
		return
	}
	steps := []resource.TestStep{
		{
			Config: ruleStatefulAppend(),
		},
	}
	runStatefulTest(steps, t)
}

func ruleStatefulAppend() string {
	return `
	resource "sysdig_secure_rule_stateful" "stateful_rule_append" {
	  name = "API Gateway Enumeration Detected"
	  source = "awscloudtrail_stateful"
	  ruletype = "STATEFUL_SEQUENCE"
	  append = true
	  exceptions {
      values = jsonencode([["user_abc", ["12345"]]])
      name = "user_accountid"
    }
	}`
}

func runStatefulTest(steps []resource.TestStep, t *testing.T) {
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
		Steps: steps,
	})

}
