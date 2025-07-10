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
			{
				Config: driftPolicyWithAllActions(rText()),
			},
			{
				Config: driftPolicyWithoutActions(rText()),
			},
			{
				Config: driftPolicyWithoutNotificationChannel(rText()),
			},
			{
				Config: driftPolicyWithoutExceptions(rText()),
			},
			{
				Config: driftPolicyWithMountedVolumeDriftEnabled(rText()),
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

func driftPolicyWithAllActions(name string) string {
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
    use_regex = true

    exceptions {
      items = ["/usr/bin/sh"]
    }
    prohibited_binaries {
      items = ["/usr/bin/curl"]
    }
    process_based_exceptions {
      items = ["/usr/bin/curl"]
    } 
    process_based_prohibited_binaries {
      items = ["/usr/bin/sh"]
    }
  }

  actions {
    prevent_drift = true
    container = "stop"
    capture {
      seconds_before_event = 5
      seconds_after_event = 10
      name = "testcapture"
    }
  }

  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}

`, secureNotificationChannelEmailWithName(name), name)
}

func driftPolicyWithoutActions(name string) string {
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
    use_regex = true

    exceptions {
      items = ["/usr/bin/sh"]
    }
    prohibited_binaries {
      items = ["/usr/bin/curl"]
    }
    process_based_exceptions {
      items = ["/usr/bin/curl"]
    } 
  }

  actions {}

  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}

`, secureNotificationChannelEmailWithName(name), name)
}

func driftPolicyWithoutNotificationChannel(name string) string {
	return fmt.Sprintf(`
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
    process_based_exceptions {
      items = ["/usr/bin/curl"]
    }
    process_based_prohibited_binaries {
      items = ["/usr/bin/sh"]
    }
  }

  actions {
    prevent_drift = true
  }
}

`, name)
}

func driftPolicyWithoutExceptions(name string) string {
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
  }

  actions {
    prevent_drift = true
  }

  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}

`, secureNotificationChannelEmailWithName(name), name)
}

func driftPolicyWithMountedVolumeDriftEnabled(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_drift_policy" "sample" {

  name        = "Test Drift Policy %s"
  description = "Test Drift Policy Description"
  enabled     = true
  severity    = 4

  rule {
    description = "Test Drift Rule Description"

    enabled = true
    mounted_volume_drift_enabled = true

    exceptions {
      items = ["/usr/bin/sh"]
    }
    prohibited_binaries {
      items = ["/usr/bin/curl"]
    }
    process_based_exceptions {
      items = ["/usr/bin/curl"]
    }
    process_based_prohibited_binaries {
      items = ["/usr/bin/sh"]
    }
  }
}
  `, name)
}
