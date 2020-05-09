package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNotificationChannelDataSource(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
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
				Config: notificationChannelEmailWithNameAndDatasource(rText),
			},
		},
	})
}

func notificationChannelEmailWithNameAndDatasource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel" "sample_email" {
	name = "%s"
	enabled = true
	type = "EMAIL"
	recipients = "root@localhost.com"
	notify_when_ok = false
	notify_when_resolved = false
}

data "sysdig_secure_notification_channel" "sample_email" {
	name = sysdig_secure_notification_channel.sample_email.name
}
`, name)
}
