//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common || tf_acc_ibm_secure || tf_acc_ibm_common

package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/buildinfo"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureTeam(t *testing.T) {
	t.Cleanup(func() {
		handleReport(t)
	})

	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureTeamWithName(rText()),
			},
			{
				Config: secureTeamMinimumConfiguration(rText()),
			},
			{
				Config: secureTeamWithPlatformMetricsIBM(rText()),
				SkipFunc: func() (bool, error) {
					return !buildinfo.IBMSecure, nil
				},
			},
			{
				ResourceName:      "sysdig_secure_team.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureTeamWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "sample" {
  name               = "sample-%s"
  description        = "%s"
  scope_by           = "container"
  filter             = "container.image.repo = \"sysdig/agent\""
}
`, name, name)
}

func secureTeamMinimumConfiguration(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "sample" {
  name      = "sample-%s"
}`, name)
}

func secureTeamWithPlatformMetricsIBM(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "sample" {
  name = "sample-%s"
  enable_ibm_platform_metrics = true
  ibm_platform_metrics = "foo in (\"0\") and bar in (\"3\")"
}`, name)
}
