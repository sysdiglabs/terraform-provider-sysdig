package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccAlertMetric(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_MONITOR_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN must be set for acceptance tests")
			}
		},
		Providers: map[string]terraform.ResourceProvider{
			"sysdig": sysdig.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: alertMetricWithName(rText()),
			},
			{
				Config: alertMetricWithoutScopeWithName(rText()),
			},
		},
	})
}

func alertMetricWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_metric" "sample" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"
	severity = 3

	metric = "avg(avg(cpu.used.percent)) > 50"
	scope = "kubernetes.cluster.name in (\"pulsar\")"
	
	trigger_after_minutes = 10

	enabled = false

	multiple_alerts_by = ["kubernetes.deployment.name"]

	capture {
		filename = "TERRAFORM_TEST"
		duration = 15
	}
}
`, name, name)
}

func alertMetricWithoutScopeWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_metric" "sample2" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"
	severity = 3

	metric = "avg(avg(cpu.used.percent)) > 50"
	
	trigger_after_minutes = 10

	enabled = false

	multiple_alerts_by = ["kubernetes.deployment.name"]

	capture {
		filename = "TERRAFORM_TEST"
		duration = 15
	}
}
`, name, name)
}
