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

func TestAccMonitorNotificationChannelGoogleChatDataSource(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: sysdigOrIBMMonitorPreCheck(t),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: monitorNotificationChannelGoogleChat(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_google_chat.nc_google_chat", "id", "sysdig_monitor_notification_channel_google_chat.nc_google_chat", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_google_chat.nc_google_chat", "name", "sysdig_monitor_notification_channel_google_chat.nc_google_chat", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_google_chat.nc_google_chat", "url", "sysdig_monitor_notification_channel_google_chat.nc_google_chat", "url"),
				),
			},
		},
	})
}

func monitorNotificationChannelGoogleChat(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_google_chat" "nc_google_chat" {
	name = "Example Channel %s - google chat"
	url = "https://chat.googleapis.com/v1/spaces/XXXXXX/messages?key=XXXXXXXXXXXXXXXXX"
}

data "sysdig_monitor_notification_channel_google_chat" "nc_google_chat" {
	name = sysdig_monitor_notification_channel_google_chat.nc_google_chat.name
}
`, name)
}
