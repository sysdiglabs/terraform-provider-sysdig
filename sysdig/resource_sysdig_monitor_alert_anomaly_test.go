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

func TestAccAlertAnomaly(t *testing.T) {
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
				Config: alertAnomalyWithName(rText()),
			},
		},
	})
}

func alertAnomalyWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_anomaly" "sample" {
	name = "TERRAFORM TEST - ANOMALY %s"
	description = "TERRAFORM TEST - ANOMALY %s"
	severity = 0

	monitor = ["cpu.cores.used", "cpu.cores.used.percent", "cpu.stolen.percent", "cpu.used.percent"]

	multiple_alerts_by = ["kubernetes.deployment.name"]
	scope = "kubernetes.cluster.name in (\"pulsar\")"
	
	trigger_after_minutes = 10

	enabled = false

	capture {
		filename = "TERRAFORM_TEST"
		duration = 15
	}
}
`, name, name)
}
