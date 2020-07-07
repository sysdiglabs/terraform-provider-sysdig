package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccMonitorNotificationChannelEmail(t *testing.T) {
	//var ncBefore, ncAfter monitor.NotificationChannel

	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_MONITOR_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN must be set for acceptance tests")
			}
		},
		Providers: map[string]terraform.ResourceProvider{
			"sysdig": sysdig.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: monitorNotificationChannelEmailWithName(rText()),
			},
			{
				Config: monitorNotificationChannelEmailWithNameInReverseOrder(rText()),
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
