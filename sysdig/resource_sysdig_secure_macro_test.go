//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"os"
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
				Config: macroWithMacro(rText(), rText()),
			},
			{
				Config: macroWithMacroAndList(rText(), rText(), rText()),
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

// macroAppendToDefault exercises append mode, which the API only supports
// against a macro provided by Sysdig (not one created by Terraform in the
// same test) - see https://github.com/sysdiglabs/terraform-provider-sysdig/pull/749
// for why appending onto a freshly-created custom macro doesn't work.
func macroAppendToDefault() string {
	return `
resource "sysdig_secure_macro" "sample2" {
  name       = "container"
  condition  = "and always_true"
  append     = true
}
`
}

func macroWithMacro(name1, name2 string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_macro" "sample3" {
  name = "terraform_test_%s"
  condition = "always_true"
}

resource "sysdig_secure_macro" "sample4" {
  name = "terraform_test_%s"
  condition = "never_true and ${sysdig_secure_macro.sample3.name}"
}
`, name1, name2)
}

func macroWithMacroAndList(name1, name2, name3 string) string {
	return fmt.Sprintf(`
%s

resource "sysdig_secure_macro" "sample5" {
  name = "terraform_test_%s"
  condition = "fd.name in (${sysdig_secure_list.sample.name})"
}

resource "sysdig_secure_macro" "sample6" {
  name = "terraform_test_%s"
  condition = "never_true and ${sysdig_secure_macro.sample5.name}"
}
`, listWithName(name3), name1, name2)
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

// TestAccMacroVersionIsPersistedOnUpdate guards against the update path dropping
// the version returned by the API. The resource sends the version it holds in
// state on update; if the response version is not written back, state keeps the
// pre-update value while the backend has already moved on.
func TestAccMacroVersionIsPersistedOnUpdate(t *testing.T) {
	name := randomText(10)

	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: macroWithName(name),
				Check:  resource.TestCheckResourceAttr("sysdig_secure_macro.sample", "version", "1"),
			},
			{
				// An in-place update bumps the backend version to 2; state must agree.
				Config: macroUpdatedWithName(name),
				Check:  resource.TestCheckResourceAttr("sysdig_secure_macro.sample", "version", "2"),
			},
		},
	})
}

// TestAccMacroAppendToCustomMacro covers appending a macro the customer owns,
// rather than one Sysdig ships (which macroAppendToDefault already covers).
//
// This is only supported on backends that migrated to the v2 macro storage.
// Earlier backends reject it with "The field 'name' must not be the same as
// another Secure UI macro", so the test is opt-in via SYSDIG_SECURE_MACROS_V2.
func TestAccMacroAppendToCustomMacro(t *testing.T) {
	if os.Getenv("SYSDIG_SECURE_MACROS_V2") == "" {
		t.Skip("Skipping append-to-custom-macro test because SYSDIG_SECURE_MACROS_V2 is not set")
		return
	}

	name := randomText(10)

	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: macroAppendToCustomMacro(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_secure_macro.base", "append", "false"),
					resource.TestCheckResourceAttr("sysdig_secure_macro.extension", "append", "true"),
				),
			},
		},
	})
}

func macroAppendToCustomMacro(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_macro" "base" {
  name      = "terraform_test_%s"
  condition = "proc.name = foo"
}

resource "sysdig_secure_macro" "extension" {
  name      = sysdig_secure_macro.base.name
  condition = "or proc.name = bar"
  append    = true
}
`, name)
}
