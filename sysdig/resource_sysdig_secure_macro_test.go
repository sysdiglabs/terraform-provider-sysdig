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

func TestAccMacro(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	fixedRandomText := rText()

	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: macroWithName(rText()),
			},
			{
				Config: macroWithName(fixedRandomText),
			},
			{
				Config: macroUpdatedWithName(fixedRandomText),
			},
			{
				ResourceName:      "sysdig_secure_macro.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: macroAppendToDefault(),
			},
			{
				Config: macroWithMacroPart1(rText()),
			},
			{
				Config: macroWithMacroPart2(rText()),
			},
			{
				Config: listWithName(rText()),
			},
			{
				Config: macroWithMacroAndListPart1(rText()),
			},
			{
				Config: macroWithMacroAndListPart2(rText()),
			},
			{
				Config: macroWithMinimumEngineVersion(rText()),
			},
		},
	})
}

func macroWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_macro" "sample" {
  name = "terraform_test_%s"
  condition = "always_true"
}
`, name)
}

func macroUpdatedWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_macro" "sample" {
  name = "terraform_test_%s"
  condition = "never_true"
}
`, name)
}

func macroAppendToDefault() string {
	return `
resource "sysdig_secure_macro" "sample2" {
  name = "container"
  condition = "and always_true"
  append = true
}
`
}

func macroWithMacroPart1(name1 string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_macro" "sample3" {
  name = "terraform_test_%s"
  condition = "always_true"
}
`, name1)
}

func macroWithMacroPart2(name2 string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_macro" "sample4" {
  name = "terraform_test_%s"
  condition = "never_true and ${sysdig_secure_macro.sample3.name}"
}
`, name2)
}

func macroWithMacroAndListPart1(name1 string) string {
	return fmt.Sprintf(`

resource "sysdig_secure_macro" "sample5" {
  name = "terraform_test_%s"
  condition = "fd.name in (${sysdig_secure_list.sample.name})"
}

`, name1)
}

func macroWithMacroAndListPart2(name1 string) string {
	return fmt.Sprintf(`

resource "sysdig_secure_macro" "sample6" {
  name = "terraform_test_%s"
  condition = "never_true and ${sysdig_secure_macro.sample5.name}"
}
`, name1)
}

func macroWithMinimumEngineVersion(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_macro" "sample" {
	minimum_engine_version = 13
	name = "terraform_test_%s"
	condition = "always_true"
  }
`, name)
}
