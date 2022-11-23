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

func TestAccAlertV2Metric(t *testing.T) {
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
				Config: alertV2MetricWithName(rText()),
			},
			{
				Config: alertV2MetricWithGroupBy(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_alert_v2_metric.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func alertV2MetricWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_metric" "sample" {

	name = "TERRAFORM TEST - METRICV2 %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	op = ">="
	threshold = 50
	warning_threshold = 20

	scope {
		label = "kube_cluster_name"
		op = "in"
		values = ["thom-cluster1", "demo-env-prom"]
	}

	trigger_after_minutes = 15

}

`, name)
}

func alertV2MetricWithGroupBy(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_alert_v2_metric" "sample" {

		name = "TERRAFORM TEST - METRICV2 %s"
		metric = "sysdig_container_cpu_used_percent"
		group_aggregation = "avg"
		time_aggregation = "avg"
		op = ">="
		threshold = 50
		warning_threshold = 20
		group_by = ["kube_cluster_name", "kube_pod_name", "cloud_provider_tag_Owner",]
	
		scope {
			label = "kube_cluster_name"
			op = "in"
			values = ["thom-cluster1", "demo-env-prom"]
		}
	
		trigger_after_minutes = 15
	
	}
	
	`, name)
}
