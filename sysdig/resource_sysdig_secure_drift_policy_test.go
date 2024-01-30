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
  name = "Test Malware Policy %s"
  description = "Test Malware Policy Description %s"
  enabled = true
  severity = 4
  
  %s
  
  actions {
    prevent_drift = true
  }
  
  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}
`, secureNotificationChannelEmailWithName(name), name, name, ruleDrift(name))
}

func ruleDrift(name string) string {
	return fmt.Sprintf(`
rules {
  description = "Test Malware Rule Description %s"
  tags = ["tag1", "tag2"]

  details {
    mode = "enabled"
    exceptions {
      items = ["304ef4cdda3463b24bf53f9cdd69ad3ecdab0842e7e70e2f3cfbb9f14e1c4ae6"]
      match_items = true
    }
    prohibited_binaries {
      items = ["304ef4cdda3463b24bf53f9cdd69ad3ecdab0842e7e70e2f3cfbb9f14e1c4ae6"]
      match_items = true
    }
  }
}
`, name)
}

// TODO: Specify only a single rule type!
