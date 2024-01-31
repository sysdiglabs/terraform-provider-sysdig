//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccDriftPolicy(t *testing.T) {
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
				Config: driftPolicyWithName(rText()),
			},
		},
	})
}

func driftPolicyWithName(name string) string {
	return fmt.Sprintf(`
%s

resource "sysdig_secure_drift_policy" "sample" {
  name        = "Test Drift Policy %s"
  description = "Test Drift Policy Description"
  enabled     = true
  severity    = 4

  rule {
    description = "Test Drift Rule Description"

    enabled = true

    exceptions {
      items = ["/usr/bin/sh"]
    }
    prohibited_binaries {
      items = ["/usr/bin/curl"]
    }
  }

  actions {
    prevent_drift = true
  }

  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}

`, secureNotificationChannelEmailWithName(name), name)
}

// TODO: Specify only a single rule type!