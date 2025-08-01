//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common || tf_acc_ibm_secure || tf_acc_ibm_common || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelPrometheusAlertManager(t *testing.T) {
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
				Config: secureNotificationChannelPrometheusAlertManagerWithName(rText()),
			},
			{
				ResourceName:            "sysdig_secure_notification_channel_prometheus_alert_manager.sample-channel1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
			{
				Config: secureNotificationChannelPrometheusAlertManagerWithNameWithAdditionalheaders(rText()),
			},
			{
				ResourceName:            "sysdig_secure_notification_channel_prometheus_alert_manager.sample-channel2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
			{
				Config: secureNotificationChannelPrometheusAlertManagerWithNameWithAllowInsecureConnections(rText()),
			},
			{
				ResourceName:            "sysdig_secure_notification_channel_prometheus_alert_manager.sample-channel3",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
			{
				Config: secureNotificationChannelPrometheusAlertManagerSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:            "sysdig_secure_notification_channel_prometheus_alert_manager.sample-channel4",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
		},
	})
}

func secureNotificationChannelPrometheusAlertManagerWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_prometheus_alert_manager" "sample-channel1" {
	name = "Example Channel %s - prometheus alert manager"
	enabled = true
	url = "https://testurl.com/xxx"
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}

func secureNotificationChannelPrometheusAlertManagerWithNameWithAdditionalheaders(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_secure_notification_channel_prometheus_alert_manager" "sample-channel2" {
		name = "Example Channel %s - prometheus alert manager With Additional Headers"
		enabled = true
		url = "https://testurl.com/xxx"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
		additional_headers = {
			"custom-Header": "TestHeader"
		}
	}`, name)
}

func secureNotificationChannelPrometheusAlertManagerWithNameWithAllowInsecureConnections(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_secure_notification_channel_prometheus_alert_manager" "sample-channel3" {
		name = "Example Channel %s - prometheus alert manager with insecure connections"
		enabled = true
		url = "https://testurl.com/xxx"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
		allow_insecure_connections = true
	}`, name)
}

func secureNotificationChannelPrometheusAlertManagerSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_secure_notification_channel_prometheus_alert_manager" "sample-channel4" {
		name = "Example Channel %s - prometheus alert manager with share with current team"
		enabled = true
		url = "https://testurl.com/xxx"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
		share_with_current_team = true
	}`, name)
}
