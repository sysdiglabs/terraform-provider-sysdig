//go:build tf_acc_sysdig_monitor || tf_acc_ibm_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccAlertV2Downtime(t *testing.T) {
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
				Config: alertV2DowntimeWithName(rText()),
			},
			{
				Config: alertV2DowntimeWithGroupBy(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_alert_v2_downtime.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func alertV2DowntimeWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_downtime" "sample" {

	name = "TERRAFORM TEST - DOWNTIMEV2 %s"
  metric = "sysdig_container_up"
  threshold = 75

	scope {
		label = "kube_cluster_name"
		operator = "in"
		values = ["thom-cluster1", "demo-env-prom"]
	}

	trigger_after_minutes = 15

}

`, name)
}

func alertV2DowntimeWithGroupBy(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_alert_v2_downtime" "sample" {

		name = "TERRAFORM TEST - DOWNTIMEV2 %s"
		metric = "sysdig_container_up"
		threshold = 75
		group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]

		scope {
			label = "kube_cluster_name"
			operator = "in"
			values = ["thom-cluster1", "demo-env-prom"]
		}

		trigger_after_minutes = 15

	}

	`, name)
}
