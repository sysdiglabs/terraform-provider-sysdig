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

func TestAccCompositePolicy(t *testing.T) {
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
				Config: compositePolicyWithName(rText()),
			},
		},
	})
}

func compositePolicyWithName(name string) string {
	return fmt.Sprintf(`
%s

resource "sysdig_secure_composite_policy" "sample" {
  name = "Test Malware Policy %s"
  description = "Test Malware Policy Description %s"
  enabled = true
  severity = 4
  
  %s
  
  actions {
    prevent_malware = true
  }
  
  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}
`, secureNotificationChannelEmailWithName(name), name, name, ruleMalware(name))
}

func ruleMalware(name string) string {
	return fmt.Sprintf(`
rules {
  enabled = true
  description = "Test Malware Rule Description %s"
  tags = ["tag1", "tag2"]

  details {
    use_managed_hashes = true
    additionals_hashes = {
      "304ef4cdda3463b24bf53f9cdd69ad3ecdab0842e7e70e2f3cfbb9f14e1c4ae6" = []
      "f953f70b9132340e2782cba8feef678642693d5fa6acccaebfbc452c5bb358a5" = ["hash_name"]
    }
    ignore_hashes = {
      "6ac3c336e4094835293a3fed8a4b5fedde1b5e2626d9838fed50693bba00af0e" = ["ignore_hashes"]
    }
  }
}
`, name)
}

// TODO: Specify only a single rule type!
