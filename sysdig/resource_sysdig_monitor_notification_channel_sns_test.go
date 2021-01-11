package sysdig_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonitorNotificationChannelSNS(t *testing.T) {
	//var ncBefore, ncAfter monitor.NotificationChannel

	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.Test(t, resource.TestCase{
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
				Config: monitorNotificationChannelAmazonSNSWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_sns.sample-amazon-sns",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorNotificationChannelAmazonSNSWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_sns" "sample-amazon-sns" {
	name = "Example Channel %s - Amazon SNS"
	enabled = true
	topics = ["arn:aws:sns:us-east-1:273489009834:my-alerts2", "arn:aws:sns:us-east-1:279948934544:my-alerts"]
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}
