package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelWebhook(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		Providers: map[string]terraform.ResourceProvider{
			"sysdig": sysdig.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: secureNotificationChannelWebhookWithName(rText()),
			},
		},
	})
}

func secureNotificationChannelWebhookWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_webhook" "sample-webhook" {
	name = "Example Channel %s - Webhook"
	enabled = true
	url = "localhost:8080"
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}
