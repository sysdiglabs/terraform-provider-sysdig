package sysdig_test

import (
	"context"
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
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
	resource.AddTestSweepers("sysdig_macros", &resource.Sweeper{
		Name: "sysdig_macros",

		F: func(region string) error {
			apiToken := os.Getenv("SYSDIG_SECURE_API_TOKEN")
			secureURL := os.Getenv("SYSDIG_SECURE_URL")
			secureTLS := os.Getenv("SYSDIG_SECURE_INSECURE_TLS")
			isSecure, err := strconv.ParseBool(secureTLS)
			if err != nil {
				return err
			}
			secureClient := secure.NewSysdigSecureClient(
				apiToken, secureURL, isSecure)

			ctx := context.Background()
			summaries, err := secureClient.GetMacroSummaries(ctx)

			if err != nil {
				return err
			}

			if summaries != nil {
				for _, summary := range *summaries {
					if strings.Contains(summary.Name, "terraform_test_") ||
						strings.Contains(summary.Name, "container") {
						for _, id := range summary.Ids {
							err := secureClient.DeleteMacro(ctx, id)
							_ = err
						}
					}

				}
			}
			return nil

		},
	})
}

func TestAccMacro(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	fixedRandomText := rText()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
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
