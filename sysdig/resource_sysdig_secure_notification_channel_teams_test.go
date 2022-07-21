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

func TestAccSecureNotificationChannelTeams(t *testing.T) {
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
				Config: secureNotificationChannelTeamsWithName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_teams.sample-teams",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureNotificationChannelTeamsWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_teams" "sample-teams" {
	name = "Example Channel %s - Teams"
	enabled = true
	url = "https://webhook.office.com/webhookb2/XXXXXXXX/IncomingWebhook/XXXXXXXX/XXXXXXXXXX"
	channel = "Example Channel"
	notify_when_ok = true
	notify_when_resolved = true
}`, name)
}
