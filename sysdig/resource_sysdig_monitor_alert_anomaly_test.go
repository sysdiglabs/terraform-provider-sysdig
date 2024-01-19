//go:build tf_acc_sysdig_monitor || tf_acc_ibm_monitor || tf_acc_onprem_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccAlertAnomaly(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigIBMMonitorAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: alertAnomalyWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_alert_anomaly.sample",
				ImportState:       true,
				ImportStateVerify: true,
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
		filename = "TERRAFORM_TEST.scap"
		duration = 15
	}

	custom_notification {
		title = "{{__alert_name__}} is {{__alert_status__}}"
		prepend = "{{kubernetes.deployment.name}}"
		append = "{{kubernetes.deployment.name}}"
	}
}
`, name, name)
}
