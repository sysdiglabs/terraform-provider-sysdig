// This test changes the tenant-wide Monitor password login setting, so it only builds into the Monitor suites:
// running it from the Secure suites too would let two processes race on the same setting.
//go:build tf_acc_sysdig_monitor || tf_acc_onprem_monitor

package sysdig_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSSOGlobalSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			preCheckAnyEnv(t, SysdigMonitorApiTokenEnv)()
			restoreSSOGlobalSettingsAfterTest(t, "monitor")
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: ssoGlobalSettingsConfig("monitor", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_global_settings.test",
						"product",
						"monitor",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_global_settings.test",
						"is_password_login_enabled",
						"true",
					),
				),
			},
			{
				Config: ssoGlobalSettingsConfig("monitor", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_global_settings.test",
						"is_password_login_enabled",
						"false",
					),
				),
			},
			{
				ResourceName:      "sysdig_sso_global_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func ssoGlobalSettingsConfig(product string, isPasswordLoginEnabled bool) string {
	return fmt.Sprintf(`
resource "sysdig_sso_global_settings" "test" {
  product                    = "%s"
  is_password_login_enabled  = %t
}
`, product, isPasswordLoginEnabled)
}

// restoreSSOGlobalSettingsAfterTest records the current global SSO settings of the product and registers a
// cleanup that puts them back. The resource's Delete is a no-op, so without this the test leaves password
// login disabled for the whole customer. t.Cleanup runs even when the test fails or aborts.
// ssoGlobalSettingsRestoreTimeout bounds the snapshot and restore requests, as the API client has no timeout of its own.
const ssoGlobalSettingsRestoreTimeout = 5 * time.Minute

func restoreSSOGlobalSettingsAfterTest(t *testing.T, product string) {
	t.Helper()

	client := ssoGlobalSettingsTestClient()
	ctx, cancel := context.WithTimeout(context.Background(), ssoGlobalSettingsRestoreTimeout)
	defer cancel()
	original, err := client.GetSSOGlobalSettings(ctx, product)
	if err != nil {
		t.Fatalf("cannot read the %s global SSO settings to restore them after the test: %v", product, err)
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), ssoGlobalSettingsRestoreTimeout)
		defer cancel()
		if _, err := client.UpdateSSOGlobalSettings(ctx, product, original); err != nil {
			t.Errorf("cannot restore the %s global SSO settings (is_password_login_enabled=%t): %v",
				product, original.IsPasswordLoginEnabled, err)
		}
	})
}

// ssoGlobalSettingsTestClient builds an API client from the same environment the provider reads.
func ssoGlobalSettingsTestClient() v2.SSOGlobalSettingsInterface {
	// Parsed like the provider's boolean attributes, so 1 or TRUE also enable it. An invalid value counts as false.
	insecure, _ := strconv.ParseBool(os.Getenv("SYSDIG_MONITOR_INSECURE_TLS"))
	return v2.NewSysdigMonitor(
		v2.WithToken(os.Getenv(SysdigMonitorApiTokenEnv)),
		v2.WithURL(envOrDefault("SYSDIG_MONITOR_URL", "https://app.sysdigcloud.com")),
		v2.WithInsecure(insecure),
	)
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
