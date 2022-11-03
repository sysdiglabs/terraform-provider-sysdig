package sysdig_test

import (
	"context"
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/common"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func init() {
	resource.AddTestSweepers("sysdig_data_user", &resource.Sweeper{
		Name: "sysdig_data_user",

		F: func(region string) error {

			apiToken := os.Getenv("SYSDIG_MONITOR_API_TOKEN")
			monitorURL := os.Getenv("SYSDIG_MONITOR_URL")
			monitorTLS := os.Getenv("SYSDIG_MONITOR_INSECURE_TLS")
			isSecure := false
			var err error
			if len(monitorTLS) > 0 {
				isSecure, err = strconv.ParseBool(monitorTLS)
				if err != nil {
					return err
				}
			}
			commonClient := common.NewSysdigCommonClient(
				apiToken, monitorURL, isSecure)

			ctx := context.Background()
			user, err := commonClient.GetUserByEmail(ctx, "terraform-test+user@sysdig.com")

                        if err != nil {
                                return err
                        }
			if user == nil {
				return nil
			}

			err = commonClient.DeleteUser(ctx, user.ID)
			if err != nil {
				return nil
			}
			return nil

		},
	})
}

func TestAccDataUser(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			monitor := os.Getenv("SYSDIG_MONITOR_API_TOKEN")
			secure := os.Getenv("SYSDIG_SECURE_API_TOKEN")
			if monitor == "" && secure == "" {
				t.Fatal("either SYSDIG_MONITOR_API_TOKEN or SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: getUser(),
			},
		},
	})
}

func getUser() string {
	return `
resource "sysdig_user" "sample" {
  email = "terraform-test+user@sysdig.com"
}

data "sysdig_user" "me" {
	depends_on = ["sysdig_user.sample"]
	email = sysdig_user.sample.email
}
`
}
