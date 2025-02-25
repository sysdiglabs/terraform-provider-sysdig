package sysdig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestRuleStatefulAppends(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleStatefulAppend(randomString()),
		},
	}
	runTest(steps, t)
}

func ruleStatefulAppend(name string) string {
	return `
	resource "sysdig_secure_rule_stateful" "stateful_rule_append" {
	  name = "API Gateway Enumeration Detected"
	  source = "awscloudtrail_stateful"
	  ruletype = "STATEFUL_SEQUENCE"
	  append = true
	  exceptions {
      values = jsonencode([["abc", ["docker.io/library/busybox"]]])
      name = "tf_append_%s"
    }
	}`
}
