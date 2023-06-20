//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common || tf_acc_ibm_secure || tf_acc_ibm_common

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccNotificationChannelDataSource(t *testing.T) {
	t.Cleanup(func() {
		handleReport(t)
	})

	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
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
resource "sysdig_secure_notification_channel_email" "sample_email" {
	name = "%s"
	enabled = true
	recipients = ["root@localhost.com"]
	notify_when_ok = false
	notify_when_resolved = false
}

data "sysdig_secure_notification_channel" "sample_email" {
	depends_on = [sysdig_secure_notification_channel_email.sample_email]
	name = sysdig_secure_notification_channel_email.sample_email.name
}
`, name)
}
