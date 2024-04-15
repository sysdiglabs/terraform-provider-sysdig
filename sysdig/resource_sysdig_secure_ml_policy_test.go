//go:build tf_acc_sysdig_secure || tf_acc_policies_aws || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMLPolicy(t *testing.T) {
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
				Config: mlPolicyWithName(rText()),
			},
			{
				Config: mlPolicyWithoutNotificationChannel(rText()),
			},
		},
	})
}

func mlPolicyWithName(name string) string {
	return fmt.Sprintf(`
%s

resource "sysdig_secure_ml_policy" "sample" {
  name        = "Test ML Policy %s"
  description = "Test ML Policy Description"
  enabled     = true
  severity    = 4

  rule {
    description = "Test ML Rule Description"

    cryptomining_trigger {
      enabled   = true
      threshold = 2
    }
  }

  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}

`, secureNotificationChannelEmailWithName(name), name)
}

func mlPolicyWithoutNotificationChannel(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_ml_policy" "sample" {
  name        = "Test ML Policy %s"
  description = "Test ML Policy Description"
  enabled     = true
  severity    = 4

  rule {
    description = "Test ML Rule Description"

    cryptomining_trigger {
      enabled   = true
      threshold = 2
    }
  }

}

`, name)
}
