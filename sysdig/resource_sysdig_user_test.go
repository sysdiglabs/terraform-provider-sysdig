//go:build tf_acc_sysdig || tf_acc_sysdig_monitor || tf_acc_sysdig_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccUser(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	name := rText()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: userWithSystemRole(name),
			},
			{
				Config: userWithName(name),
			},
			{
				Config: userWithoutSystemRole(name),
			},
			{
				Config: userMinimumConfiguration(name),
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
  email      = "terraform-test+user%[1]s@sysdig.com"
  system_role = "ROLE_USER"
  first_name = "%[1]s"
  last_name  = "%[1]s"
}`, name)
}

func userWithSystemRole(name string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample" {
  email      = "terraform-test+user%[1]s@sysdig.com"
  system_role = "ROLE_USER"
  first_name = "%[1]s"
  last_name  = "%[1]s"
}`, name)
}

func userWithoutSystemRole(name string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample" {
  email      = "terraform-test+user%[1]s@sysdig.com"
  first_name = "%[1]s"
  last_name  = "%[1]s"
}`, name)
}

func userMinimumConfiguration(name string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample" {
  email      = "terraform-test+user%s@sysdig.com"
}`, name)
}
