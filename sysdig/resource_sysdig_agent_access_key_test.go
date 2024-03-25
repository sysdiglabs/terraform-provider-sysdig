//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure

package sysdig_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccAgentAccessKeyResource(t *testing.T) {

	limit := 1
	reservation := 0
	metadata := map[string]string{
		"test":            "yes",
		"acceptance_test": "true",
		"status":          "new",
	}

	updatedLimit := 10
	updatedMetadata := map[string]string{
		"test":            "yes",
		"acceptance_test": "true",
		"status":          "updated",
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: getAgentAccessKeyWithMetadata(limit, reservation, true, metadata),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "limit", strconv.Itoa(limit)),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "reservation", strconv.Itoa(reservation)),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "metadata.test", metadata["test"]),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "metadata.acceptance_test", metadata["acceptance_test"]),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "metadata.status", metadata["status"]),
				),
			},
			{
				Config: getAgentAccessKeyWithMetadata(updatedLimit, reservation, true, updatedMetadata),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "limit", strconv.Itoa(updatedLimit)),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "reservation", strconv.Itoa(reservation)),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "metadata.test", updatedMetadata["test"]),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "metadata.acceptance_test", updatedMetadata["acceptance_test"]),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "metadata.status", updatedMetadata["status"]),
				),
			},
			{
				Config: getAgentAccessKeyWithMetadata(updatedLimit, reservation, false, updatedMetadata),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "limit", strconv.Itoa(updatedLimit)),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "reservation", strconv.Itoa(reservation)),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "enabled", strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "metadata.test", updatedMetadata["test"]),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "metadata.acceptance_test", updatedMetadata["acceptance_test"]),
					resource.TestCheckResourceAttr("sysdig_agent_access_key.my_agent_access_key", "metadata.status", updatedMetadata["status"]),
				),
			},
			{
				ResourceName:      "sysdig_agent_access_key.my_agent_access_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getAgentAccessKeyWithMetadata(limit int, reservation int, enabled bool, metadata map[string]string) string {
	// Build the metadata string for Terraform configuration
	metadataStr := "metadata = {\n"
	for key, value := range metadata {
		metadataStr += fmt.Sprintf("    \"%s\" = \"%s\"\n", key, value)
	}
	metadataStr += "  }\n"

	// Return the full Terraform configuration
	return fmt.Sprintf(`
resource "sysdig_agent_access_key" "my_agent_access_key" {
  limit       = %d
  reservation = %d
  enabled     = %t
  %s
}

data "sysdig_agent_access_key" "data" {
  id = sysdig_agent_access_key.my_agent_access_key.id
}
`, limit, reservation, enabled, metadataStr)
}
