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

func TestAccAlertGroupOutlier(t *testing.T) {
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
				Config: alertGroupOutlierWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_alert_group_outlier.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func alertGroupOutlierWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_group_outlier" "sample" {
	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	description = "TERRAFORM TEST - GROUP OUTLIER %s"
	severity = 6

	monitor = ["cpu.cores.used", "cpu.cores.used.percent", "cpu.stolen.percent", "cpu.used.percent"]

	scope = "kubernetes.cluster.name in (\"pulsar\")"
	
	trigger_after_minutes = 10

	enabled = false

	capture {
		filename = "TERRAFORM_TEST.scap"
		duration = 15
	}
}
`, name, name)
}
