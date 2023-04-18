//go:build tf_acc_sysdig || tf_acc_ibm

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorNotificationChannelVictorOps(t *testing.T) {
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
				Config: monitorNotificationChannelVictorOpsWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_victorops.sample-victorops",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorNotificationChannelVictorOpsWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_victorops" "sample-victorops" {
	name = "Example Channel %s - VictorOps"
	enabled = true
	api_key = "1234342-4234243-4234-2"
	routing_key = "My team"
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}
