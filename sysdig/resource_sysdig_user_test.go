package sysdig_test

import (
	"context"
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func init() {
	resource.AddTestSweepers("sysdig_user", &resource.Sweeper{
		Name: "sysdig_user",

		F: func(region string) error {
			apiToken := os.Getenv("SYSDIG_MONITOR_API_TOKEN")
			monitorURL := os.Getenv("SYSDIG_MONITOR_URL")
			monitorTLS := os.Getenv("SYSDIG_MONITOR_INSECURE_TLS")
			isSecure, err := strconv.ParseBool(monitorTLS)
			if err != nil {
				return err
			}
			commonClient := common.NewSysdigCommonClient(
				apiToken, monitorURL, isSecure)

			ctx := context.Background()
			user, err := commonClient.GetUserByEmail(ctx, "terraform-test+user@sysdig.com")

			if err != nil {
				return err
			}

			err = commonClient.DeleteUser(ctx, user.ID)
			if err != nil {
				return err
			}
			return nil

		},
	})
}

func TestAccUser(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

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
				Config: userWithSystemRole(rText()),
			},
			{
				Config: userWithName(rText()),
			},
			{
				Config: userWithoutSystemRole(rText()),
			},
			{
				Config: userMinimumConfiguration(),
			},
			{
				ResourceName:      "sysdig_user.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func userWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample" {
  email      = "terraform-test+user@sysdig.com"
  system_role = "ROLE_USER"
  first_name = "%s"
  last_name  = "%s"
}`, name, name)
}

func userWithSystemRole(name string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample" {
  email      = "terraform-test+user@sysdig.com"
  system_role = "ROLE_CUSTOMER"
  first_name = "%s"
  last_name  = "%s"
}`, name, name)
}

func userWithoutSystemRole(name string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample" {
  email      = "terraform-test+user@sysdig.com"
  first_name = "%s"
  last_name  = "%s"
}`, name, name)
}

func userMinimumConfiguration() string {
	return `
resource "sysdig_user" "sample" {
  email      = "terraform-test+user@sysdig.com"
}`
}

func testAccUserResourceDestroy(s *terraform.State) error {
	// retrieve the connection established in Provider configuration

	conn := testAccProvider.Meta().(*ExampleClient)

	// loop through the resources in state, verifying each widget
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "example_widget" {
			continue
		}

		// Retrieve our widget by referencing it's state ID for API lookup
		request := &example.DescribeWidgets{
			IDs: []string{rs.Primary.ID},
		}

		response, err := conn.DescribeWidgets(request)
		if err == nil {
			if len(response.Widgets) > 0 && *response.Widgets[0].ID == rs.Primary.ID {
				return fmt.Errorf("Widget (%s) still exists.", rs.Primary.ID)
			}

			return nil
		}

		// If the error is equivalent to 404 not found, the widget is destroyed.
		// Otherwise return the error
		if !strings.Contains(err.Error(), "Widget not found") {
			return err
		}
	}

	return nil
}