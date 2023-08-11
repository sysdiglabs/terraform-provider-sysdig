//go:build tf_acc_sysdig_monitor || tf_acc_ibm_monitor

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

func TestAccMonitorNotificationChannelIBMEventNotificationDataSource(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	ibmEventNotificationInstanceId := os.Getenv("IBM_EVENT_NOTIFICATION_INSTANCE_ID")
	if ibmEventNotificationInstanceId == "" {
		t.Skip("Skipping tests on sysdig_monitor_notification_channel_ibm_event_notification resource because IBM_EVENT_NOTIFICATION_INSTANCE_ID is not set")
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigIBMMonitorAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: monitorNotificationChannelIBMEventNotification(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_ibm_event_notification.nc_ibm_event_notification", "id", "sysdig_monitor_notification_channel_ibm_event_notification.nc_ibm_event_notification", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_ibm_event_notification.nc_ibm_event_notification", "name", "sysdig_monitor_notification_channel_ibm_event_notification.nc_ibm_event_notification", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_ibm_event_notification.nc_ibm_event_notification", "instance_id", "sysdig_monitor_notification_channel_ibm_event_notification.nc_ibm_event_notification", "instance_id"),
				),
			},
		},
	})
}

func monitorNotificationChannelIBMEventNotification(name, ibmEventNotificationInstanceId string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_ibm_event_notification" "nc_ibm_event_notification" {
	name = "Example Channel %s - IBM Event Notification"
	instance_id = "%s"
}

data "sysdig_monitor_notification_channel_ibm_event_notification" "nc_ibm_event_notification" {
	name = sysdig_monitor_notification_channel_ibm_event_notification.nc_ibm_event_notification.name
}
`, name, ibmEventNotificationInstanceId)
}
