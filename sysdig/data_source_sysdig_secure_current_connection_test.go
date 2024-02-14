//go:build tf_acc_sysdig_secure || tf_acc_onprem_secure

package sysdig_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureConnection(t *testing.T) {
	dataSourceResourceName := "data.sysdig_secure_connection.current"

	apiToken := os.Getenv("SYSDIG_SECURE_API_TOKEN")
	secureUrl := os.Getenv("SYSDIG_SECURE_URL")
	if secureUrl == "" {
		secureUrl = "https://secure.sysdig.com"
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if apiToken == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN and must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: getSysdigSecureCurrentConnection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceResourceName, "secure_url", secureUrl),
					resource.TestCheckResourceAttr(dataSourceResourceName, "secure_api_token", apiToken),
				),
			},
		},
	})
}

func getSysdigSecureCurrentConnection() string {
	return `
data "sysdig_secure_connection" "current" {
}
`
}
