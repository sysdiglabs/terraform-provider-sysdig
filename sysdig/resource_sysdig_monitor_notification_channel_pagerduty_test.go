//go:build tf_acc_sysdig || tf_acc_ibm

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorNotificationChannelPagerduty(t *testing.T) {
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
				Config: monitorNotificationChannelPagerdutyWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_pagerduty.sample-pagerduty",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorNotificationChannelPagerdutyWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_pagerduty" "sample-pagerduty" {
	name = "Example Channel %s - Pagerduty"
	enabled = true
	account = "account"
	service_key = "XXXXXXXXXX"
	service_name = "sysdig"
	notify_when_ok = true
	notify_when_resolved = true
	send_test_notification = false
}`, name)
}
