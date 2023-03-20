//go:build sysdig

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

func TestAccNotificationChannelEmailDataSource(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
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
				Config: notificationChannelEmail(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_email.nc_email", "id", "sysdig_monitor_notification_channel_email.nc_email", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_email.nc_email", "name", "sysdig_monitor_notification_channel_email.nc_email", "name"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_monitor_notification_channel_email.nc_email", "recipients.*", "root@localhost.com"),
				),
			},
		},
	})
}

func notificationChannelEmail(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_email" "nc_email" {
	name = "%s"
	recipients = ["root@localhost.com"]
}

data "sysdig_monitor_notification_channel_email" "nc_email" {
	name = sysdig_monitor_notification_channel_email.nc_email.name
}
`, name)
}
