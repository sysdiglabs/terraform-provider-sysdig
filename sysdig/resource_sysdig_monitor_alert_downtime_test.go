//go:build tf_acc_sysdig || tf_acc_sysdig_monitor || tf_acc_ibm || tf_acc_ibm_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccAlertDowntime(t *testing.T) {
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
				Config: alertDowntimeWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_alert_downtime.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func alertDowntimeWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_downtime" "sample" {
	name = "TERRAFORM TEST - DOWNTIME %s"
	description = "TERRAFORM TEST - DOWNTIME %s"
	severity = 2

	entities_to_monitor = ["host.hostName", "host.mac"]
	scope = "kubernetes.cluster.name in (\"pulsar\")"
	
	trigger_after_minutes = 10
	trigger_after_pct = 99

	enabled = false

	capture {
		filename = "TERRAFORM_TEST.scap"
		duration = 15
	}
}
`, name, name)
}
