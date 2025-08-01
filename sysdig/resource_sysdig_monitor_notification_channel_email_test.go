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

func TestAccMonitorNotificationChannelEmail(t *testing.T) {
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
				Config: monitorNotificationChannelEmailWithName(rText()),
			},
			{
				Config: monitorNotificationChannelEmailWithNameInReverseOrder(rText()),
			},
			{
				Config: monitorNotificationChannelEmailSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:            "sysdig_monitor_notification_channel_email.sample_email",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
		},
	})
}

func monitorNotificationChannelEmailWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_email" "sample_email" {
	name = "%s"
	recipients = ["root@localhost.com", "bar@localhost.com"]
	enabled = true
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}

func monitorNotificationChannelEmailWithNameInReverseOrder(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_email" "sample_email" {
	name = "%s"
	recipients = ["bar@localhost.com", "root@localhost.com"]
	enabled = false
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}

func monitorNotificationChannelEmailSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_email" "sample_email" {
	name = "%s"
	share_with_current_team = true
	recipients = ["bar@localhost.com", "root@localhost.com"]
	enabled = false
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}
