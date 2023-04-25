//go:build sysdig_secure || tf_acc_policies

package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccRuleNetwork(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: ruleNetworkWithName(rText()),
			},
			{
				Config: ruleNetworkWithoutTags(rText()),
			},
			{
				Config: ruleNetworkWithTCP(rText()),
			},
			{
				Config: ruleNetworkWithUDP(rText()),
			},
			{
				ResourceName:      "sysdig_secure_rule_network.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: ruleNetworkWithMinimalConfig(rText()),
			},
		},
	})
}

func ruleNetworkWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_network" "foo" {
  name = "TERRAFORM TEST %s" // ID
  description = "TERRAFORM TEST %s"
  tags = ["network", "cis"]

  block_inbound = true
  block_outbound = true

  tcp {
    matching = true // default
    ports = [80, 443]
  }

  udp {
    matching = true // default
    ports = [80, 443]
  }
}`, name, name)
}

func ruleNetworkWithoutTags(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_network" "foo" {
  name = "TERRAFORM TEST %s" // ID
  description = "TERRAFORM TEST %s"

  block_inbound = true
  block_outbound = true

  tcp {
    matching = true // default
    ports = [80, 443]
  }

  udp {
    matching = true // default
    ports = [80, 443]
  }
}`, name, name)
}

func ruleNetworkWithTCP(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_network" "foo" {
  name = "TERRAFORM TEST %s" // ID
  description = "TERRAFORM TEST %s"

  block_inbound = true
  block_outbound = true

  tcp {
    matching = true // default
    ports = [80, 443]
  }
}`, name, name)
}

func ruleNetworkWithUDP(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_network" "foo" {
  name = "TERRAFORM TEST %s" // ID
  description = "TERRAFORM TEST %s"

  block_inbound = true
  block_outbound = true

  udp {
    matching = true // default
    ports = [80, 443]
  }
}`, name, name)
}

func ruleNetworkWithMinimalConfig(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_network" "foo-minimal" {
  name = "TERRAFORM TEST %s" // ID
  description = "TERRAFORM TEST %s"

  block_inbound = true
  block_outbound = true
}`, name, name)
}
