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

func TestAccList(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	fixedRandomText := rText()

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
				Config: listWithName(rText()),
			},
			{
				Config: listWithName(fixedRandomText),
			},
			{
				Config: listUpdatedWithName(fixedRandomText),
			},
			{
				Config: listAppendToDefault(),
			},
			{
				Config: listWithList(rText(), rText()),
			},
		},
	})
}

func listWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_list" "sample" {
  name = "terraform_test_%s"
  items = ["foo", "bar"]
}
`, name)
}

func listUpdatedWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_list" "sample" {
  name = "terraform_test_%s"
  items = ["foo", "bar", "baz"]
}
`, name)
}

func listAppendToDefault() string {
	return fmt.Sprintf(`
resource "sysdig_secure_list" "sample2" {
  name = "allowed_k8s_nodes"
  items = ["foo", "bar"]
  append = true
}
`)
}

func listWithList(name1, name2 string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_list" "sample3" {
  name = "terraform_test_%s"
  items = ["foo", "bar"]
}

resource "sysdig_secure_list" "sample4" {
  name = "terraform_test_%s"
  items = [sysdig_secure_list.sample3.name]
}
`, name1, name2)
}
