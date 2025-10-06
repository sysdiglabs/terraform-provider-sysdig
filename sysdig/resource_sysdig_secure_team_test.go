//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common || tf_acc_ibm_secure || tf_acc_ibm_common || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/draios/terraform-provider-sysdig/buildinfo"

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
				Config: secureTeamWithAgentCliAndRapidResponse(randomText(10)),
			},
			{
				Config: secureTeamMinimumConfiguration(randomText(10)),
			},
			{
				Config: secureTeamWithPostureZones(randomText(10)),
				SkipFunc: func() (bool, error) {
					return buildinfo.OnpremSecure, nil
				},
			},
			{
				Config: secureTeamWithPostureZonesAndAllZones(randomText(10)),
				SkipFunc: func() (bool, error) {
					return buildinfo.OnpremSecure, nil
				},
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

func secureTeamWithAgentCliAndRapidResponse(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "sample" {
  name                   = "sample-%s"
  description            = "%s"
  scope_by               = "container"
  can_use_agent_cli      = false
	can_use_rapid_response = true
}
`, name, name)
}

func secureTeamMinimumConfiguration(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "sample" {
  name = "sample-%s"
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
