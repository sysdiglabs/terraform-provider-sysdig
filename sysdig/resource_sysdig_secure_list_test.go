//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccList(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	fixedRandomText := rText()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
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
				ResourceName:      "sysdig_secure_list.sample",
				ImportState:       true,
				ImportStateVerify: true,
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
	return `
resource "sysdig_secure_list" "sample2" {
  name = "allowed_k8s_nodes"
  items = ["foo", "bar"]
  append = true
}
`
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
