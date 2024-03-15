//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_common || tf_acc_ibm_monitor || tf_acc_ibm_common || tf_acc_onprem_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorNotificationChannelGoogleChat(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: sysdigOrIBMMonitorPreCheck(t),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: monitorNotificationChannelGoogleChatWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_google_chat.sample_google_chat1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelGoogleChatSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_google_chat.sample_google_chat2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorNotificationChannelGoogleChatWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_google_chat" "sample_google_chat1" {
	name = "Example Channel %s - google chat"
	enabled = true
	url = "https://chat.googleapis.com/v1/spaces/XXXXXX/messages?key=XXXXXXXXXXXXXXXXX"
	notify_when_ok = true
	notify_when_resolved = true
}`, name)
}

func monitorNotificationChannelGoogleChatSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_google_chat" "sample_google_chat2" {
	name = "Example Channel %s - google chat"
	enabled = true
	url = "https://chat.googleapis.com/v1/spaces/XXXXXX/messages?key=XXXXXXXXXXXXXXXXX"
	notify_when_ok = true
	notify_when_resolved = true
	share_with_current_team = true
}`, name)
}
