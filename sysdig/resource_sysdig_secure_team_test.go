//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common || tf_acc_ibm_secure || tf_acc_ibm_common

package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/buildinfo"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureTeam(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureTeamWithName(randomText(10)),
			},
			{
				Config: secureTeamMinimumConfiguration(randomText(10)),
			},
			{
				Config: secureTeamWithPlatformMetricsIBM(randomText(10)),
				SkipFunc: func() (bool, error) {
					return !buildinfo.IBMSecure, nil
				},
			},
			{
				Config: secureTeamWithPostureZones(randomText(10)),
			},
			{
				Config: secureTeamWithPostureZonesAndAllZones(randomText(10)),
				ExpectError: regexp.MustCompile(
					fmt.Sprintf("if %s is enabled, %s must be omitted",
						sysdig.SchemaAllZones,
						sysdig.SchemaZonesIDsKey,
					),
				),
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

func secureTeamWithPostureZones(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_posture_zone" "z1" {
  name = "Zone-%[1]s"
}

resource "sysdig_secure_team" "sample" {
  name     = "sample-%[1]s"
  zone_ids = [sysdig_secure_posture_zone.z1.id]
}`, name)
}

func secureTeamWithPostureZonesAndAllZones(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_posture_zone" "z1" {
  name = "Zone-%[1]s"
}

resource "sysdig_secure_team" "sample" {
  name      = "sample-%[1]s"
  zone_ids  = [sysdig_secure_posture_zone.z1.id]
  all_zones = true
}`, name)
}
