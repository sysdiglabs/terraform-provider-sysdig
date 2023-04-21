//go:build tf_acc_sysdig || tf_acc_monitor

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

func TestAccAlertV2Event(t *testing.T) {
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
				Config: alertV2Event(rText()),
			},
			{
				Config: alertV2EventWithSources(rText()),
			},
			{
				Config: alertV2EventWithScrambledSources(rText()),
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

	trigger_after_minutes = 15

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

	trigger_after_minutes = 15

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

	trigger_after_minutes = 15

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

	trigger_after_minutes = 15

}

`, name)
}
