package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureBenchmarkTask(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
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
				Config: secureBenchmarkTaskWithName(rText()),
			},
			{
				Config: multiRegionSecureBenchmarkTaskWithName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_benchmark_task.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureBenchmarkTaskWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_benchmark_task" "sample" {
  name     = "%s"
  schedule = "0 6 * * *"
  schema   = "aws_foundations_bench-1.3.0"
  scope    = "aws.accountId = \"123456789012\" and aws.region = \"us-west-2\""
  enabled  = true
}
`, name)
}

func multiRegionSecureBenchmarkTaskWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_benchmark_task" "sample2" {
  name     = "%s"
  schedule = "0 6 * * *"
  schema   = "aws_foundations_bench-1.3.0"
  scope    = "aws.accountId = \"123456789012\" and aws.region in (\"us-east-1\", \"us-west-2\", \"eu-central-1\")"
  enabled  = true
}
`, name)
}
