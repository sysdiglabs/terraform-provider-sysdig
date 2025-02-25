package sysdig_test

import (
	"os"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestRuleGuardDutyAppends(t *testing.T) {
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

func randomString() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

func runTest(steps []resource.TestStep, t *testing.T) {
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
