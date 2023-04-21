//go:build tf_acc_sysdig || tf_acc_monitor

package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorNotificationChannelOpsGenie(t *testing.T) {
	//var ncBefore, ncAfter monitor.NotificationChannel

	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			monitor := os.Getenv("SYSDIG_MONITOR_API_TOKEN")
			ibmMonitor := os.Getenv("SYSDIG_IBM_MONITOR_API_KEY")
			if monitor == "" && ibmMonitor == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN or SYSDIG_IBM_MONITOR_API_KEY must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: monitorNotificationChannelOpsGenieWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_opsgenie.sample-opsgenie",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelOpsGenieWithNameAndRegion(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_opsgenie.sample-opsgenie-2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorNotificationChannelOpsGenieWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_opsgenie" "sample-opsgenie" {
	name = "Example Channel %s - OpsGenie"
	enabled = true
	api_key = "2349324-342354353-5324-23"
	notify_when_ok = false
	notify_when_resolved = false
}`, name)
}

func monitorNotificationChannelOpsGenieWithNameAndRegion(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_opsgenie" "sample-opsgenie-2" {
	name = "Example Channel %s - OpsGenie - 2"
	enabled = true
	api_key = "2349324-342354353-5324-23"
	notify_when_ok = false
	notify_when_resolved = false
	region = "EU"
}`, name)
}
