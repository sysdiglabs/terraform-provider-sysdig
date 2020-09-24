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

func TestAccSecureNotificationChannelEmail(t *testing.T) {
	//var ncBefore, ncAfter secure.NotificationChannel

	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureNotificationChannelEmailWithName(rText()),
			},
			{
				Config: secureNotificationChannelEmailWithNameInReverseOrder(rText()),
			},
		},
	})
}

func secureNotificationChannelEmailWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_email" "sample_email" {
	name = "%s"
	recipients = ["root@localhost.com", "bar@localhost.com"]
	enabled = true
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}

func secureNotificationChannelEmailWithNameInReverseOrder(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_email" "sample_email" {
	name = "%s"
	recipients = ["bar@localhost.com", "root@localhost.com"]
	enabled = false
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}
