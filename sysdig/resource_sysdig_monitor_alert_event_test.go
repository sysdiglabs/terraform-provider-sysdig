//go:build sysdig_monitor

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

func TestAccAlertEvent(t *testing.T) {
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
				Config: alertEventWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_alert_event.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func alertEventWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_event" "sample" {
	name = "TERRAFORM TEST - EVENT %s"
	description = "TERRAFORM TEST - EVENT %s"
	severity = 4

	event_name = "deployment"
	source = "kubernetes"
	event_rel = ">"
	event_count = 2

	multiple_alerts_by = ["kubernetes.deployment.name"]
	scope = "kubernetes.cluster.name in (\"pulsar\")"
	
	trigger_after_minutes = 10

	enabled = false
}
`, name, name)
}
