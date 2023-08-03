//go:build tf_acc_ibm_monitor || tf_acc_ibm_common

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

func TestAccMonitorNotificationChannelIBMEventNotification(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

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
				Config: monitorNotificationChannelIBMEventNotificationWithName(rText(), ibmEventNotificationInstanceId),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_ibm_event_notification.sample1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelIBMEventNotificationSharedWithCurrentTeam(rText(), ibmEventNotificationInstanceId),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_ibm_event_notification.sample2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorNotificationChannelIBMEventNotificationWithName(name, ibmEventNotificationInstanceId string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_notification_channel_ibm_event_notification" "sample1" {
		name = "Example Channel %s - IBM Event Notification"
		enabled = true
		instance_id = "%s"
		notify_when_ok = true
		notify_when_resolved = true
}`, name, ibmEventNotificationInstanceId)
}

func monitorNotificationChannelIBMEventNotificationSharedWithCurrentTeam(name, ibmEventNotificationInstanceId string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_ibm_event_notification" "sample2" {
	name = "Example Channel %s - IBM Event Notification"
	share_with_current_team = true
	enabled = true
	instance_id = "%s"
	notify_when_ok = true
	notify_when_resolved = true
}`, name, ibmEventNotificationInstanceId)
}
