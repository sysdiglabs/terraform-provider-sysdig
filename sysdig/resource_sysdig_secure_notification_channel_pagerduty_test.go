//go:build tf_acc_sysdig || tf_acc_sysdig_secure

package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelPagerduty(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
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
				Config: secureNotificationChannelPagerdutyWithName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_pagerduty.sample-pagerduty",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureNotificationChannelPagerdutyWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_pagerduty" "sample-pagerduty" {
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
