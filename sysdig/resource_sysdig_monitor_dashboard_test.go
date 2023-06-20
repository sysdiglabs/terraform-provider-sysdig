//go:build tf_acc_sysdig_monitor

package sysdig_test

import (
	"fmt"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccDashboard(t *testing.T) {
	t.Cleanup(func() {
		handleReport(t)
	})

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
			{
				Config: sharedDashboard(rText()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_monitor_dashboard.dashboard", "share.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("sysdig_monitor_dashboard.dashboard", "share.*", map[string]string{
						"member.#":      "1",
						"role":          "ROLE_RESOURCE_EDIT",
						"member.0.type": "TEAM",
					}),
				),
			},
			{
				Config: multiplePanelsDashboardWithDisplayInfo(rText()),
			},
			{
				Config: timeChartDashboardWithLegend(
					rText(),
					"true",
					"true",
					"bottom",
					"inline",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard",
						"panel.0.legend.0.enabled",
						"true",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard",
						"panel.0.legend.0.show_current",
						"true",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard",
						"panel.0.legend.0.layout",
						"inline",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard",
						"panel.0.legend.0.position",
						"bottom",
					),
				),
			},
			{
				Config: timeChartDashboardWithLegend(
					rText(),
					"false",
					"false",
					"right",
					"table",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard",
						"panel.0.legend.0.enabled",
						"false",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard",
						"panel.0.legend.0.show_current",
						"false",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard",
						"panel.0.legend.0.layout",
						"table",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard",
						"panel.0.legend.0.position",
						"right",
					),
				),
			},
			{
				Config: dashboardWithQueryFormat(rText(), *v2.NewFormat(
					"number",
					"K",
					"K",
					1,
					"left",
					"10s",
					"nullGap",
				)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.unit", string(v2.FormatUnitNumber),
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.input_format", "K",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.display_format", "K",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.decimals", "1",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.y_axis", "left",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.min_interval", "10s",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.null_value_display_mode", "nullGap",
					),
				),
			},
			{
				Config: dashboardWithQueryFormat(rText(), *v2.NewFormat(
					"percent",
					"0-100",
					"0-100",
					0,
					"right",
					"20s",
					"connectSolid",
				)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.unit", "percent",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.input_format", "0-100",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.display_format", "0-100",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.decimals", "0",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.y_axis", "right",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.min_interval", "20s",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.null_value_display_mode", "connectSolid",
					),
				),
			},
			{
				Config: dashboardWithQueryFormat(rText(), *v2.NewFormat(
					"data",
					"B",
					"B",
					2,
					"left",
					"1d",
					"nullZero",
				)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.unit", "data",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.input_format", "B",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.display_format", "B",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.decimals", "2",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.y_axis", "left",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.min_interval", "1d",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.null_value_display_mode", "nullZero",
					),
				),
			},
			{
				Config: dashboardWithQueryFormat(rText(), *v2.NewFormat(
					"data rate",
					"B/s",
					"B/s",
					2,
					"left",
					"1d",
					"nullGap",
				)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.unit", "data rate",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.input_format", "B/s",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.display_format", "B/s",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.decimals", "2",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.y_axis", "left",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.min_interval", "1d",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.null_value_display_mode", "nullGap",
					),
				),
			},
			{
				Config: dashboardWithQueryFormat(rText(), *v2.NewFormat(
					"number rate",
					"/s",
					"/s",
					1,
					"right",
					"10m",
					"connectDotted",
				)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.unit", "number rate",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.input_format", "/s",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.display_format", "/s",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.decimals", "1",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.y_axis", "right",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.min_interval", "10m",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.null_value_display_mode", "connectDotted",
					),
				),
			},
			{
				Config: dashboardWithQueryFormat(rText(), *v2.NewFormat(
					"time",
					"ns",
					"ms",
					1,
					"right",
					"10m",
					"connectDotted",
				)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.unit", "time",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.input_format", "ns",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.display_format", "ms",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.decimals", "1",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.y_axis", "right",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.min_interval", "10m",
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_dashboard.dashboard", "panel.0.query.0.format.0.null_value_display_mode", "connectDotted",
					),
				),
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

        legend {
            show_current = true
            position = "bottom"
            layout = "inline"
        }

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "percent"

            format {
                display_format = "auto"
                input_format = "0-100"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
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
            
            format {
                display_format = "auto"
                input_format = "0-100"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
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

	scope {
		metric = "agent.id"
		comparator = "in"
		value = ["foo", "bar"]
		variable = "agent_id"
	}

	scope {
		metric = "agent.name"
		comparator = "equals"
		value = ["name"]
	}

	scope {
		metric = "kubernetes.namespace.name"
		variable = "k8_ns"
	}

	panel {
		pos_x = 0
		pos_y = 0
		width = 12 # Maximum size: 24
		height = 6
		type = "timechart"
		name = "example panel"
		description = "description"

        legend {
            show_current = true
            position = "bottom"
            layout = "inline"
        }

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "percent"

            format {
                display_format = "auto"
                input_format = "0-100"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
		}
		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent{ns_name=$k8s_ns}[$__interval]))"
			unit = "number"

            format {
                display_format = "auto"
                input_format = "1"
                y_axis = "auto"
				null_value_display_mode = "nullGap"
            }
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

            format {
                display_format = "auto"
                input_format = "ns"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
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

        legend {
            show_current = true
            position = "bottom"
            layout = "inline"
        }

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "percent"

            format {
                display_format = "auto"
                input_format = "0-100"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
		}
		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "number"

            format {
                display_format = "auto"
                input_format = "1"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
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
            
            format {
                display_format = "auto"
                input_format = "ns"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
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

func sharedDashboard(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "a_team" {
  name      = "sample-%s"

  entrypoint {
	type = "Explore"
  }
}

resource "sysdig_monitor_dashboard" "dashboard" {
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

            format {
                display_format = "auto"
                input_format = "0-100"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
		}
	}
	share {
		role = "ROLE_RESOURCE_EDIT"
		member {
			type = "TEAM"
			id = sysdig_monitor_team.a_team.id
		}
   }
}
`, name, name, name)
}

func multiplePanelsDashboardWithDisplayInfo(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_dashboard" "dashboard" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"

	scope {
		metric = "agent.id"
		comparator = "in"
		value = ["foo", "bar"]
		variable = "agent_id"
	}

	scope {
		metric = "agent.name"
		comparator = "equals"
		value = ["name"]
	}

	scope {
		metric = "kubernetes.namespace.name"
		variable = "k8_ns"
	}

	panel {
		pos_x = 0
		pos_y = 0
		width = 12 # Maximum size: 24
		height = 6
		type = "timechart"
		name = "example panel"
		description = "description"

        legend {
            show_current = true
            position = "bottom"
            layout = "inline"
        }

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "percent"
			display_info {
				display_name                      = "hostname"
				time_series_display_name_template = "{{host_hostname}}"
				type                              = "lines"
			}

            format {
                display_format = "auto"
                input_format = "0-100"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
		}
		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent{ns_name=$k8s_ns}[$__interval]))"
			unit = "number"
			display_info {
				time_series_display_name_template = "{{host_hostname}}"
				type                              = "stackedArea"
			}

            format {
                display_format = "auto"
                input_format = "1"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
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
            
            format {
                display_format = "auto"
                input_format = "ns"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
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

func timeChartDashboardWithLegend(name, enabled, showCurrent, position, layout string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_dashboard" "dashboard" {
	name = "TERRAFORM TEST - METRIC %[1]s"
	description = "TERRAFORM TEST - METRIC %[1]s"

	scope {
		metric = "agent.id"
		comparator = "in"
		value = ["foo", "bar"]
		variable = "agent_id"
	}

	scope {
		metric = "agent.name"
		comparator = "equals"
		value = ["name"]
	}

	scope {
		metric = "kubernetes.namespace.name"
		variable = "k8_ns"
	}

	panel {
		pos_x = 0
		pos_y = 0
		width = 12 # Maximum size: 24
		height = 6
		type = "timechart"
		name = "example panel"
		description = "description"
        
        legend {
            enabled = %[2]s
            show_current = %[3]s
            position = "%[4]s"
            layout = "%[5]s"
        }

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "percent"
			display_info {
				display_name                      = "hostname"
				time_series_display_name_template = "{{host_hostname}}"
				type                              = "lines"
			}
            
            format {
                display_format = "auto"
                input_format = "0-100"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
		}
	}
}
`, name, enabled, showCurrent, position, layout)
}

func dashboardWithQueryFormat(name string, format v2.Format) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_dashboard" "dashboard" {
	name = "TERRAFORM TEST - METRIC %[1]s"

	panel {
		pos_x = 0
		pos_y = 0
		width = 12 # Maximum size: 24
		height = 6
		type = "timechart"
		name = "example panel"
		description = "description"

        legend {
            show_current = true
            position = "bottom"
            layout = "inline"
        }

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "%[2]s"
            
            format {
                display_format = "%[3]s"
                input_format = "%[4]s"
                y_axis = "%[5]s"
                null_value_display_mode = "%[6]s"
                decimals = %[7]d
                min_interval = "%[8]s"
            }
		}
	}
}
`,
		name,
		format.Unit,
		*format.DisplayFormat,
		*format.InputFormat,
		*format.YAxis,
		*format.NullValueDisplayMode,
		*format.Decimals,
		*format.MinInterval,
	)
}
