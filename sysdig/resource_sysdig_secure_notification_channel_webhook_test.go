//go:build tf_acc_sysdig_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelWebhook(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureNotificationChannelWebhookWithName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_webhook.sample-webhook",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureNotificationChannelWebhookWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_webhook" "sample-webhook" {
	name = "Example Channel %s - Webhook"
	enabled = true
	url = "http://1.1.1.1:8080"
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}
