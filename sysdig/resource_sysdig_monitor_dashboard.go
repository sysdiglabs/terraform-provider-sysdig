package sysdig

import (
	"context"
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"
	"time"
)

func resourceSysdigMonitorDashboard() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		CreateContext: resourceSysdigDashboardCreate,
		UpdateContext: resourceSysdigDashboardUpdate,
		ReadContext:   resourceSysdigDashboardRead,
		DeleteContext: resourceSysdigDashboardDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"public_token": {
				Type:         schema.TypeString,
				ComputedWhen: []string{"public"},
				Computed:     true,
			},
			"panel": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pos_x": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validateDiagFunc(validation.IntBetween(0, 23)),
						},
						"pos_y": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validateDiagFunc(validation.IntAtLeast(0)),
						},
						"width": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validateDiagFunc(validation.IntBetween(1, 24)),
						},
						"height": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validateDiagFunc(validation.IntAtLeast(1)),
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"timechart", "number"}, false)),
						},
						"query": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"promql": {
										Type:     schema.TypeString,
										Required: true,
									},
									"unit": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"percent", "data", "data rate", "number", "number rate", "time"}, false)),
									},
								},
							},
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSysdigDashboardCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	dashboard, err := dashboardFromResourceData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	dashboardCreated, err := client.CreateDashboard(ctx, *dashboard)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(dashboardCreated.ID))
	data.Set("version", dashboardCreated.Version)
	return nil
}

func resourceSysdigDashboardUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	dashboard, err := dashboardFromResourceData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	dashboard.ID, _ = strconv.Atoi(data.Id())

	_, err = client.UpdateDashboard(ctx, *dashboard)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigDashboardRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dashboard, err := client.GetDashboardByID(ctx, id)

	if err != nil {
		data.SetId("")
		return nil
	}

	err = dashboardToResourceData(&dashboard, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigDashboardDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteDashboard(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func dashboardFromResourceData(data *schema.ResourceData) (dashboard *monitor.Dashboard, err error) {
	dashboard = &monitor.Dashboard{
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Schema:      3,
		Public:      data.Get("public").(bool),
		PublicToken: data.Get("public_token").(string),
	}

	for panelId, panelItr := range data.Get("panel").(*schema.Set).List() {
		panelInfo := panelItr.(map[string]interface{})

		var panelType string
		switch panelInfo["type"].(string) {
		case "timechart":
			panelType = "advancedTimechart"
		case "number":
			panelType = "advancedNumber"
		default:
			panic("unreachable code")
		}

		panel := monitor.Panels{
			ID:                     panelId + 1,
			Name:                   panelInfo["name"].(string),
			Description:            panelInfo["description"].(string),
			Type:                   panelType,
			AdvancedQueries:        dashboardFromResourceData_Queries(panelInfo, panelType),
			ApplyScopeToAll:        false,
			ApplySegmentationToAll: false,
			AxesConfiguration: monitor.AxesConfiguration{
				Bottom: monitor.Bottom{Enabled: true},
				Left: monitor.Left{
					Enabled:        true,
					DisplayName:    nil,
					Unit:           "auto",
					DisplayFormat:  "auto",
					Decimals:       "",
					MinValue:       0,
					MaxValue:       "",
					MinInputFormat: "ns",
					MaxInputFormat: "ns",
					Scale:          "linear",
				},
				Right: monitor.Right{
					Enabled:        true,
					DisplayName:    nil,
					Unit:           "auto",
					DisplayFormat:  "auto",
					Decimals:       "",
					MinValue:       0,
					MaxValue:       "",
					MinInputFormat: "1",
					MaxInputFormat: "1",
					Scale:          "linear",
				},
			},
			LegendConfiguration: monitor.LegendConfiguration{
				Enabled:     true,
				Position:    "right",
				Layout:      "table",
				ShowCurrent: true,
				Width:       nil,
				Height:      nil,
			},
			MarkdownSource:        nil,
			PanelTitleVisible:     false,
			TextAutosized:         false,
			TransparentBackground: false,
		}
		if panel.Type == "advancedNumber" {
			if len(panel.AdvancedQueries) > 1 {
				return nil, fmt.Errorf("a panel of type 'number' can only contain one query")
			}
			panel.NumberThresholds.Values = []interface{}{} // These values must be not nil in case of type number
			panel.NumberThresholds.Base.Severity = "none"
		}

		dashboard.Panels = append(dashboard.Panels, panel)

		dashboard.Layout = append(dashboard.Layout, monitor.Layout{
			X:       panelInfo["pos_x"].(int),
			Y:       panelInfo["pos_y"].(int),
			W:       panelInfo["width"].(int),
			H:       panelInfo["height"].(int),
			PanelID: panelId + 1,
		})

	}

	return
}

func dashboardFromResourceData_Queries(panelInfo map[string]interface{}, panelType string) (queries []monitor.AdvancedQueries) {
	for queryID, queryItr := range panelInfo["query"].(*schema.Set).List() {
		queryInfo := queryItr.(map[string]interface{})
		var queryFormat monitor.Format
		switch queryInfo["unit"].(string) {
		case "percent":
			queryFormat = monitor.NewPercentFormat()
		case "data":
			queryFormat = monitor.NewDataFormat()
		case "data rate":
			queryFormat = monitor.NewDataRateFormat()
		case "number":
			queryFormat = monitor.NewNumberFormat()
		case "number rate":
			queryFormat = monitor.NewNumberRateFormat()
		case "time":
			queryFormat = monitor.NewTimeFormat()
		default:
			panic("unreachable code")
		}

		query := monitor.AdvancedQueries{
			Query:   queryInfo["promql"].(string),
			Enabled: true,
			ID:      queryID + 1,
			DisplayInfo: monitor.DisplayInfo{ // API expects this value to be set like this
				DisplayName:                   "",
				TimeSeriesDisplayNameTemplate: "",
				Type:                          "lines",
			},
			Format: queryFormat,
		}
		if panelType == "advancedNumber" {
			query.DisplayInfo.Color = "mixed"
			query.DisplayInfo.LineWidth = 2
		}
		queries = append(queries, query)
	}

	return
}

func dashboardToResourceData(dashboard *monitor.Dashboard, data *schema.ResourceData) (err error) {
	data.Set("name", dashboard.Name)
	data.Set("description", dashboard.Description)
	data.Set("public", dashboard.Public)
	data.Set("public_token", dashboard.PublicToken)

	var panels []map[string]interface{}
	for _, panel := range dashboard.Panels {
		var queries []map[string]interface{}
		for _, query := range panel.AdvancedQueries {
			unit := ""
			switch query.Format.Unit {
			case "%":
				unit = "percent"
			case "byte":
				unit = "data"
			case "byteRate":
				unit = "data rate"
			case "number":
				unit = "number"
			case "numberRate":
				unit = "number rate"
			case "relativeTime":
				unit = "time"
			default:
				panic("unreachable code")
			}

			queries = append(queries, map[string]interface{}{
				"unit":   unit,
				"promql": query.Query,
			})
		}

		var panelLayout monitor.Layout

		for _, layout := range dashboard.Layout {
			if layout.PanelID == panel.ID {
				panelLayout = layout
			}
		}

		var panelType string
		switch panel.Type {
		case "advancedTimechart":
			panelType = "timechart"
		case "advancedNumber":
			panelType = "number"
		default:
			panic("unreachable code")
		}

		panels = append(panels, map[string]interface{}{
			"pos_x":       panelLayout.X,
			"pos_y":       panelLayout.Y,
			"width":       panelLayout.W,
			"height":      panelLayout.H,
			"name":        panel.Name,
			"description": panel.Description,
			"type":        panelType,
			"query":       queries,
		})
	}
	data.Set("panel", panels)
	data.Set("version", dashboard.Version)

	return nil
}
