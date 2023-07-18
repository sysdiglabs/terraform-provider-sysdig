//go:build tf_acc_sysdig_monitor

package sysdig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorCustomRolePermissionsDataSource(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: monitorCustomRolePermissions(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("data.sysdig_monitor_custom_role_permissions.dashboard_edit", "enriched_permissions.*", "dashboards.edit"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_monitor_custom_role_permissions.dashboard_edit", "enriched_permissions.*", "dashboards.read"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_monitor_custom_role_permissions.dashboard_edit", "enriched_permissions.*", "custom-events.read"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_monitor_custom_role_permissions.dashboard_edit", "enriched_permissions.*", "dashboard-metrics-data.read"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_monitor_custom_role_permissions.dashboard_edit", "enriched_permissions.*", "metrics-data.read"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_monitor_custom_role_permissions.dashboard_edit", "enriched_permissions.*", "alert-events.read"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_monitor_custom_role_permissions.dashboard_edit", "enriched_permissions.*", "api-token.read"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_monitor_custom_role_permissions.dashboard_edit", "enriched_permissions.*", "token.view"),

					resource.TestCheckResourceAttr("data.sysdig_monitor_custom_role_permissions.dashboard_edit", "enriched_permissions.#", "8"),
				),
			},
		},
	})
}

func monitorCustomRolePermissions() string {
	return `
data "sysdig_monitor_custom_role_permissions" "dashboard_edit" {
  requested_permissions = ["dashboards.edit", "token.view"]
}
`
}
