//go:build tf_acc_sysdig || tf_acc_ibm

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

func TestAccAlertPromql(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_MONITOR_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: alertPromqlWithName(rText()),
			},
			{
				Config: alertPromqlWithGroupName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_alert_promql.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func alertPromqlWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_promql" "sample" {
	name = "TERRAFORM TEST - PROMQL %s"
	description = "TERRAFORM TEST - PROMQL %s"
	severity = 3

	promql = "(elasticsearch_jvm_memory_used_bytes{area=\"heap\"} / elasticsearch_jvm_memory_max_bytes{area=\"heap\"}) * 100 > 80"

	trigger_after_minutes = 10

	enabled = false
}
`, name, name)
}

func alertPromqlWithGroupName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_promql" "sample" {
	name = "TERRAFORM TEST - PROMQL %s"
	description = "TERRAFORM TEST - PROMQL %s"
	severity = 3
	group_name = "sample_group_name"
	promql = "(elasticsearch_jvm_memory_used_bytes{area=\"heap\"} / elasticsearch_jvm_memory_max_bytes{area=\"heap\"}) * 100 > 80"

	trigger_after_minutes = 10

	enabled = false
}
`, name, name)
}
