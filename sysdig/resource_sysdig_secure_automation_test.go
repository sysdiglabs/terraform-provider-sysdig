//go:build tf_acc_sysdig_secure

package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureAutomation(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

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
				Config: secureAutomationBasic(rText),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_secure_automation.test", "name", fmt.Sprintf("TERRAFORM TEST %s", rText)),
					resource.TestCheckResourceAttr("sysdig_secure_automation.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("sysdig_secure_automation.test", "automation_id"),
					resource.TestCheckResourceAttrSet("sysdig_secure_automation.test", "created_at"),
					resource.TestCheckResourceAttrSet("sysdig_secure_automation.test", "customer_id"),
					resource.TestCheckResourceAttrSet("sysdig_secure_automation.test", "team_id"),
				),
			},
			{
				Config: secureAutomationUpdated(rText),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_secure_automation.test", "name", fmt.Sprintf("TERRAFORM TEST UPDATED %s", rText)),
					resource.TestCheckResourceAttr("sysdig_secure_automation.test", "enabled", "true"),
				),
			},
			{
				ResourceName:      "sysdig_secure_automation.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"automation_json", // JSON format might differ after API round-trip
					"name",            // Name is in the JSON, not separately stored
				},
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources["sysdig_secure_automation.test"]
					if !ok {
						return "", fmt.Errorf("not found: sysdig_secure_automation.test")
					}
					return rs.Primary.Attributes["automation_id"], nil
				},
			},
		},
	})
}

func secureAutomationBasic(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_automation" "test" {
  name = "TERRAFORM TEST %s"
  
  automation_json = jsonencode({
    automation = {
      name    = "Original Name From UI"
      enabled = false
      version = "v1"
      
      nodes = {
        Send_Email_1 = {
          action = {
            type = "email"
            inputs = {
              channelId = "55179"
            }
          }
          outboundEdges = []
          onError       = []
        }
      }
      
      trigger = {
        on   = "new_findings"
        when = "finding.severity in (0, 1, 2, 3)"
        outboundEdges = [
          {
            node = "Send_Email_1"
          }
        ]
      }
    }
  })
}
`, name)
}

func secureAutomationUpdated(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_automation" "test" {
  name = "TERRAFORM TEST UPDATED %s"
  
  automation_json = jsonencode({
    automation = {
      name    = "Original Name From UI"
      enabled = true
      version = "v1"
      
      nodes = {
        Send_Email_1 = {
          action = {
            type = "email"
            inputs = {
              channelId = "55179"
            }
          }
          outboundEdges = []
          onError       = []
        }
      }
      
      trigger = {
        on   = "new_findings"
        when = "finding.severity in (0, 1, 2, 3)"
        outboundEdges = [
          {
            node = "Send_Email_1"
          }
        ]
      }
    }
  })
}
`, name)
}
