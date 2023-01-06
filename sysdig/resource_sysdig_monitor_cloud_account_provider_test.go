package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

func TestCustomerProviderKeys(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accountId := rText()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_MONITOR_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: monitorCustomerProviderKey(accountId),
			},
			{
				Config: monitorCustomerProviderKeyAdditionalOptions(accountId, rText()),
			},
			{
				ResourceName:      "sysdig_monitor_cloud_account_provider.provider",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorCustomerProviderKey(rText string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_cloud_account_provider" "provider" {
  platform = "GCP"
  integration_type = "API"
  account_id = "sample-%s"
}
`, rText)
}

func monitorCustomerProviderKeyAdditionalOptions(accountId string, rText string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_cloud_account_provider" "provider" {
  platform = "Azure"
  integration_type = "Metrics Streams"
  account_id = "sample-%s"
  additional_options = "%s"
}
`, accountId, rText)
}
