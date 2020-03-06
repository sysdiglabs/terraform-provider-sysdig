package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccUser(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		Providers: map[string]terraform.ResourceProvider{
			"sysdig": sysdig.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: userWithName(rText()),
			},
			{
				Config: userWithSystemRole(rText()),
			},
			{
				Config: userWithoutSystemRole(rText()),
			},
			{
				Config: userMinimumConfiguration(rText()),
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

func userMinimumConfiguration(name string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample" {
  email      = "terraform-test+user@sysdig.com"
}`)
}
