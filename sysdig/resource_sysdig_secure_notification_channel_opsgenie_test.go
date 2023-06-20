//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common || tf_acc_ibm_secure || tf_acc_ibm_common

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelOpsGenie(t *testing.T) {
	t.Cleanup(func() {
		handleReport(t)
	})

	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureNotificationChannelOpsGenieWithName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_opsgenie.sample-opsgenie",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: secureNotificationChannelOpsGenieWithNameAndRegion(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_opsgenie.sample-opsgenie-2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: secureNotificationChannelOpsGenieSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_opsgenie.sample-opsgenie-3",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureNotificationChannelOpsGenieWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_opsgenie" "sample-opsgenie" {
	name = "Example Channel %s - OpsGenie"
	enabled = true
	api_key = "2349324-342354353-5324-23"
	notify_when_ok = false
	notify_when_resolved = false
}`, name)
}

func secureNotificationChannelOpsGenieWithNameAndRegion(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_opsgenie" "sample-opsgenie-2" {
	name = "Example Channel %s - OpsGenie - 2"
	enabled = true
	api_key = "2349324-342354353-5324-23"
	notify_when_ok = false
	notify_when_resolved = false
	region = "EU"
}`, name)
}

func secureNotificationChannelOpsGenieSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_opsgenie" "sample-opsgenie-3" {
	name = "Example Channel %s - OpsGenie - 3"
	enabled = true
    share_with_current_team = true
	api_key = "2349324-342354353-5324-23"
	notify_when_ok = false
	notify_when_resolved = false
	region = "EU"
}`, name)
}
