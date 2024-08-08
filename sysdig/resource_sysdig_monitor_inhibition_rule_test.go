//go:build tf_acc_sysdig_monitor || tf_acc_ibm_monitor || tf_acc_onprem_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorInhibitionRule(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: sysdigOrIBMMonitorPreCheck(t),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: monitorInhibitionRuleBase(),
			},
			{
				ResourceName:      "sysdig_monitor_inhibition_rule.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorInhibitionRuleWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_inhibition_rule.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorInhibitionRuleWithDescription(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_inhibition_rule.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorInhibitionRuleWithEqual(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_inhibition_rule.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorInhibitionRuleWithEnabled(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_inhibition_rule.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorInhibitionRuleBase() string {
	return `
resource "sysdig_monitor_inhibition_rule" "sample" {
  source_matchers {
    label_name = "alertname"
    operator = "EQUALS"
    value = "networkAlert"
  }

  source_matchers {
    label_name = "device_type"
    operator = "EQUALS"
    value = "firewall"
  }

  target_matchers {
    label_name = "device_type"
    operator = "REGEXP_MATCHES"
    value = ".*server.*"
  }

	target_matchers {
    label_name = "l1"
    operator = "REGEXP_MATCHES"
    value = ".*l1val.*"
  }
}`
}

func monitorInhibitionRuleWithName(text string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_inhibition_rule" "sample" {
	name = "Example Inhibition Rule %s"
  source_matchers {
    label_name = "alertname"
    operator = "NOT_REGEXP_MATCHES"
    value = "networkAlert %s"
  }

  source_matchers {
    label_name = "device_type"
    operator = "EQUALS"
    value = "firewall"
  }

  target_matchers {
    label_name = "device_type"
    operator = "REGEXP_MATCHES"
    value = ".*server.*"
  }

	target_matchers {
    label_name = "l1"
    operator = "REGEXP_MATCHES"
    value = ".*l1val.*"
  }
}`, text, text)
}

func monitorInhibitionRuleWithDescription(text string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_inhibition_rule" "sample" {
	description = "Example Inhibition Rule descr %s"
  source_matchers {
    label_name = "alertname"
    operator = "EQUALS"
    value = "networkAlert"
  }

  target_matchers {
    label_name = "device_type"
    operator = "REGEXP_MATCHES"
    value = ".*server.*"
  }

}`, text)
}

func monitorInhibitionRuleWithEqual(text string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_inhibition_rule" "sample" {
  source_matchers {
    label_name = "alertname"
    operator = "EQUALS"
    value = "networkAlert %s"
  }

  source_matchers {
    label_name = "device_type_xx"
    operator = "EQUALS"
    value = "firewall"
  }

  target_matchers {
    label_name = "device_type"
    operator = "REGEXP_MATCHES"
    value = ".*server.*"
  }

	equal = ["l1", "l2"]
}`, text)
}

func monitorInhibitionRuleWithEnabled(text string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_inhibition_rule" "sample" {
	enabled = false
  source_matchers {
    label_name = "alertname"
    operator = "EQUALS"
    value = "networkAlert %s"
  }

  source_matchers {
    label_name = "device_type_xx"
    operator = "NOT_EQUALS"
    value = "firewall"
  }

  target_matchers {
    label_name = "device_type_yy"
    operator = "REGEXP_MATCHES"
    value = ".*server.*"
  }

}`, text)
}
