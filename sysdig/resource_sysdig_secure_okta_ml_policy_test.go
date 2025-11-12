//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_policies_okta

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccOktaMLPolicy(t *testing.T) {
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
				Config: oktaMLPolicyWithName(rText()),
			},
			{
				Config: oktaMLPolicyWithoutNotificationChannel(rText()),
			},
		},
	})
}

func oktaMLPolicyWithName(name string) string {
	return fmt.Sprintf(`
%s

resource "sysdig_secure_okta_ml_policy" "sample" {
  name        = "Test Okta ML Policy %s"
  description = "Test Okta ML Policy Description"
  enabled     = true
  severity    = 4

  rule {
    description = "Test Okta ML Rule Description"

    anomalous_console_login {
      enabled   = true
      threshold = 2
    }
  }

  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}

`, secureNotificationChannelEmailWithName(name), name)
}

func oktaMLPolicyWithoutNotificationChannel(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_okta_ml_policy" "sample" {
  name        = "Test Okta ML Policy %s"
  description = "Test Okta ML Policy Description"
  enabled     = true
  severity    = 4

  rule {
    description = "Test Okta ML Rule Description"

    anomalous_console_login {
      enabled   = true
      threshold = 2
    }
  }

}

`, name)
}
