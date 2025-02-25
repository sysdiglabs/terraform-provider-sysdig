package sysdig_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
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
			Config: ruleStatefulAppend(rName()),
		},
	}
	runStatefulTest(steps, t)
}

func ruleStatefulAppend(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_secure_rule_stateful" "stateful_rule_append" {
	  name = "API Gateway Enumeration Detected"
	  source = "awscloudtrail_stateful"
	  ruletype = "STATEFUL_SEQUENCE"
	  append = true
	  exceptions {
      values = jsonencode([["abc", ["docker.io/library/busybox"]]])
      name = "tf_append_%s"
    }
	}`, name)
}

func rName() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

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
