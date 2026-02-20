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

func TestAccAlertV2Event(t *testing.T) {
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
				Config: alertV2Event(rText()),
			},
			{
				Config: alertV2EventWithTriggerAfterMinutes(rText()),
			},
			{
				Config: alertV2EventWithSources(rText()),
			},
			{
				Config: alertV2EventWithScrambledSources(rText()),
			},
			{
				Config: alertV2EventWithNotEqualOperator(rText()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_monitor_alert_v2_event.sample", "operator", "!="),
				),
			},
			{
				Config: alertV2EventWithEqualOperator(rText()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_monitor_alert_v2_event.sample", "operator", "="),
				),
			},
			{
				Config: alertV2EventWithWarningThreshold(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_alert_v2_event.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func alertV2Event(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_event" "sample" {

	name = "TERRAFORM TEST - EVENTV2 %s"
	filter = "xxx"
	operator = ">="
	threshold = 50

	scope {
		label = "kube_cluster_name"
		operator = "in"
		values = ["thom-cluster1", "demo-env-prom"]
	}

	range_seconds = 600

}

`, name)
}

func alertV2EventWithTriggerAfterMinutes(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_event" "sample" {

	name = "TERRAFORM TEST - EVENTV2 %s"
	filter = "xxx"
	operator = ">="
	threshold = 50

	scope {
		label = "kube_cluster_name"
		operator = "in"
		values = ["thom-cluster1", "demo-env-prom"]
	}

	range_seconds = 900

}

`, name)
}

func alertV2EventWithSources(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_event" "sample" {

	name = "TERRAFORM TEST - EVENTV2 %s"
	filter = "xxx"
	sources = ["kubernetes", "zup", "fix"]
	operator = ">="
	threshold = 50

	scope {
		label = "kube_cluster_name"
		operator = "in"
		values = ["thom-cluster1", "demo-env-prom"]
	}

	range_seconds = 600

}

`, name)
}

func alertV2EventWithScrambledSources(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_event" "sample" {

	name = "TERRAFORM TEST - EVENTV2 %s"
	filter = "xxx"
	sources = ["kubernetes", "fix", "zup"]
	operator = ">="
	threshold = 50

	scope {
		label = "kube_cluster_name"
		operator = "in"
		values = ["thom-cluster1", "demo-env-prom"]
	}

	range_seconds = 600

}

`, name)
}

func alertV2EventWithNotEqualOperator(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_event" "sample" {

	name = "TERRAFORM TEST - EVENTV2 %s"
	filter = "xxx"
	operator = "!="
	threshold = 50

	scope {
		label = "kube_cluster_name"
		operator = "in"
		values = ["thom-cluster1", "demo-env-prom"]
	}

	range_seconds = 600

}

`, name)
}

func alertV2EventWithEqualOperator(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_event" "sample" {

	name = "TERRAFORM TEST - EVENTV2 %s"
	filter = "xxx"
	operator = "="
	threshold = 50

	scope {
		label = "kube_cluster_name"
		operator = "in"
		values = ["thom-cluster1", "demo-env-prom"]
	}

	range_seconds = 600

}

`, name)
}

func alertV2EventWithWarningThreshold(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_event" "sample" {

	name = "TERRAFORM TEST - EVENTV2 %s"
	filter = "xxx"
	operator = ">="
	threshold = 50
	warning_threshold = 20

	scope {
		label = "kube_cluster_name"
		operator = "in"
		values = ["thom-cluster1", "demo-env-prom"]
	}

	range_seconds = 600

}

`, name)
}
