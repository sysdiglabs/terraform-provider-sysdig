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

func TestAccDashboard(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_MONITOR_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: minimumDashboard(rText()),
			},
			{
				Config: minimumNumberDashboard(rText()),
			},
			{
				Config: multiplePanelsDashboard(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_dashboard.dashboard",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: multipleUpdatedPanelsDashboard(rText()),
			},
		},
	})
}

func minimumDashboard(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_dashboard" "dashboard" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"

	panel {
		pos_x = 0
		pos_y = 0
		width = 12 # Maximum size: 24
		height = 6
		type = "timechart"
		name = "example panel"
		description = "description"

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "percent"
		}
	}
}
`, name, name)
}

func minimumNumberDashboard(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_dashboard" "dashboard_2" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"

	panel {
		pos_x = 0
		pos_y = 0
		width = 12 # Maximum size: 24
		height = 6
		type = "number"
		name = "example panel"
		description = "description"

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "percent"
		}
	}
}
`, name, name)
}

func multiplePanelsDashboard(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_dashboard" "dashboard" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"

	panel {
		pos_x = 0
		pos_y = 0
		width = 12 # Maximum size: 24
		height = 6
		type = "timechart"
		name = "example panel"
		description = "description"

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "percent"
		}
		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "number"
		}
	}

	panel {
		pos_x = 12
		pos_y = 0
		width = 12
		height = 6
		type = "number"
		name = "example panel - 2"
		description = "description of panel 2"

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "time"
		}
	}

	panel {
		pos_x = 12
		pos_y = 12
		width = 12
		height = 6
		type = "text"
		name = "example panel - 2"
		content = "description of panel 2"
		visible_title = true
		autosize_text = true
		transparent_background = true
	}
}
`, name, name)
}


func multipleUpdatedPanelsDashboard(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_dashboard" "dashboard" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"

	panel {
		pos_x = 0
		pos_y = 0
		width = 12 # Maximum size: 24
		height = 6
		type = "timechart"
		name = "example panel"
		description = "description"

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "percent"
		}
		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "number"
		}
	}

	panel {
		pos_x = 12
		pos_y = 0
		width = 12
		height = 12
		type = "number"
		name = "example panel - 2"
		description = "description of panel 2"

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "time"
		}
	}

	panel {
		pos_x = 12
		pos_y = 12
		width = 12
		height = 6
		type = "text"
		name = "example panel - 2"
		content = "description of panel 2"
		visible_title = true
		autosize_text = true
		transparent_background = true
	}
}
`, name, name)
}
