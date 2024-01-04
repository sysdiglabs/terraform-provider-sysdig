package sysdig

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/spf13/cast"
)

func resourceSysdigMonitorDashboard() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigDashboardCreate,
		UpdateContext: resourceSysdigDashboardUpdate,
		ReadContext:   resourceSysdigDashboardRead,
		DeleteContext: resourceSysdigDashboardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
			"share": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"id": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"role": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"scope": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric": {
							Type:     schema.TypeString,
							Required: true,
						},
						"comparator": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"in", "notIn", "equals", "notEquals", "contains", "notContains", "startsWith"}, false)),
						},
						"value": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"variable": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
							ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"timechart", "number", "text"}, false)),
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"visible_title": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"autosize_text": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"transparent_background": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"query": {
							Type:     schema.TypeSet,
							Optional: true,
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
									"format": {
										Type:     schema.TypeSet,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"decimals": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"display_format": {
													Type:     schema.TypeString,
													Required: true,
												},
												"input_format": {
													Type:     schema.TypeString,
													Required: true,
												},
												"min_interval": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"null_value_display_mode": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"y_axis": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"display_info": {
										Type:     schema.TypeSet,
										Optional: true,
										MinItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"display_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"time_series_display_name_template": {
													Type:     schema.TypeString,
													Required: true,
												},
												"type": {
													Type:             schema.TypeString,
													Required:         true,
													ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"lines", "stackedArea", "stackedBar"}, false)),
												},
											},
										},
									},
								},
							},
						},
						"legend": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"show_current": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"position": {
										Type:     schema.TypeString,
										Required: true,
									},
									"layout": {
										Type:     schema.TypeString,
										Required: true,
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
			"min_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "60s",
			},
		},
	}
}

func getMonitorDashboardClient(c SysdigClients) (v2.DashboardInterface, error) {
	var client v2.DashboardInterface
	var err error
	switch c.GetClientType() {
	case IBMMonitor:
		client, err = c.ibmMonitorClient()
		if err != nil {
			return nil, err
		}
	default:
		client, err = c.sysdigMonitorClientV2()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func resourceSysdigDashboardCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getMonitorDashboardClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	dashboard, err := dashboardFromResourceData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	dashboardCreated, err := client.CreateDashboard(ctx, dashboard)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(dashboardCreated.ID))
	_ = data.Set("version", dashboardCreated.Version)

	return nil
}

func resourceSysdigDashboardUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getMonitorDashboardClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	dashboard, err := dashboardFromResourceData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	dashboard.ID, _ = strconv.Atoi(data.Id())

	_, err = client.UpdateDashboard(ctx, dashboard)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigDashboardRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getMonitorDashboardClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dashboard, err := client.GetDashboard(ctx, id)
	if err != nil {
		data.SetId("")
		return nil
	}

	err = dashboardToResourceData(dashboard, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigDashboardDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getMonitorDashboardClient(i.(SysdigClients))
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

func dashboardFromResourceData(data *schema.ResourceData) (dashboard *v2.Dashboard, err error) {
	dashboard = v2.NewDashboard(data.Get("name").(string), data.Get("description").(string)).AsPublic(data.Get("public").(bool))
	dashboard.Version = cast.ToInt(data.Get("version"))
	dashboard.PublicToken = data.Get("public_token").(string)

	panels, err := panelsFromResourceData(data)
	if err != nil {
		return nil, err
	}

	scopes, err := scopeFromResourceData(data)
	if err != nil {
		return nil, err
	}
	dashboard.ScopeExpressionList = scopes

	dashboard.AddPanels(panels...)

	shares, err := sharingFromResourceData(data)
	if err != nil {
		return nil, err
	}
	dashboard.SharingSettings = shares
	dashboard.MinInterval = data.Get("min_interval").(string)
	return dashboard, nil
}

func sharingFromResourceData(data *schema.ResourceData) (sharingSettings []*v2.SharingOptions, err error) {
	for _, share := range data.Get("share").(*schema.Set).List() {
		shareInfo := share.(map[string]interface{})
		memberInfo := shareInfo["member"].(*schema.Set).List()[0].(map[string]interface{})
		sharingSettings = append(sharingSettings,
			&v2.SharingOptions{
				Member: v2.SharingMember{
					Type: memberInfo["type"].(string),
					ID:   memberInfo["id"].(int),
				},
				Role: shareInfo["role"].(string),
			})
	}
	return
}

func panelsFromResourceData(data *schema.ResourceData) (panels []*v2.Panels, err error) {
	for _, panelItr := range data.Get("panel").(*schema.Set).List() {
		panelInfo := panelItr.(map[string]interface{})

		var panel *v2.Panels
		switch panelInfo["type"].(string) {
		case "timechart":
			panel, err = timechartPanelFromResourceData(panelInfo)
		case "number":
			panel, err = numberPanelFromResourceData(panelInfo)
		case "text":
			panel, err = textPanelFromResourceData(panelInfo)
		default:
			return nil, fmt.Errorf("unsupported panel type %s", panelInfo["type"])
		}
		if err != nil {
			return nil, err
		}

		panels = append(panels, panel)
	}
	return
}

func defaultLegendConfiguration() *v2.LegendConfiguration {
	return &v2.LegendConfiguration{
		Enabled:     false,
		Position:    "bottom",
		Layout:      "table",
		ShowCurrent: false,
	}
}

func legendFromResourceData(data interface{}) *v2.LegendConfiguration {
	if data == nil {
		return defaultLegendConfiguration()
	}

	legendList := data.(*schema.Set).List()

	if len(legendList) == 0 {
		return defaultLegendConfiguration()
	}

	legend := legendList[0].(map[string]interface{})
	return &v2.LegendConfiguration{
		Enabled:     legend["enabled"].(bool),
		Position:    legend["position"].(string),
		Layout:      legend["layout"].(string),
		ShowCurrent: legend["show_current"].(bool),
	}
}

func timechartPanelFromResourceData(panelInfo map[string]interface{}) (*v2.Panels, error) {
	panel := &v2.Panels{
		ID:                     0,
		Name:                   panelInfo["name"].(string),
		Description:            panelInfo["description"].(string),
		Type:                   v2.PanelTypeTimechart,
		ApplyScopeToAll:        false,
		ApplySegmentationToAll: false,
		AxesConfiguration: &v2.AxesConfiguration{
			Bottom: v2.Bottom{Enabled: true},
			Left: v2.Left{
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
			Right: v2.Right{
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
		LegendConfiguration:   legendFromResourceData(panelInfo["legend"]),
		MarkdownSource:        nil,
		PanelTitleVisible:     false,
		TextAutosized:         false,
		TransparentBackground: false,
	}

	_, err := panel.WithLayout(panelInfo["pos_x"].(int), panelInfo["pos_y"].(int), panelInfo["width"].(int), panelInfo["height"].(int))
	if err != nil {
		return nil, err
	}

	queries, err := queriesFromResourceData(panelInfo, panel)
	if err != nil {
		return nil, err
	}
	if len(queries) == 0 {
		return nil, fmt.Errorf("no query defined for timechart panel")
	}

	_, err = panel.AddQueries(queries...)
	if err != nil {
		return nil, err
	}

	return panel, nil
}

func numberPanelFromResourceData(panelInfo map[string]interface{}) (*v2.Panels, error) {
	panel := &v2.Panels{
		ID:                     0,
		Name:                   panelInfo["name"].(string),
		Description:            panelInfo["description"].(string),
		Type:                   v2.PanelTypeNumber,
		ApplyScopeToAll:        false,
		ApplySegmentationToAll: false,
		AxesConfiguration: &v2.AxesConfiguration{
			Bottom: v2.Bottom{Enabled: true},
			Left: v2.Left{
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
			Right: v2.Right{
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
		LegendConfiguration: &v2.LegendConfiguration{
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
		NumberThresholds: &v2.NumberThresholds{
			Values: []interface{}{}, // These values must be not nil in case of type number
			Base: v2.NumberThresholdBase{
				Severity: "none",
			},
		},
	}

	_, err := panel.WithLayout(panelInfo["pos_x"].(int), panelInfo["pos_y"].(int), panelInfo["width"].(int), panelInfo["height"].(int))
	if err != nil {
		return nil, err
	}

	queries, err := queriesFromResourceData(panelInfo, panel)
	if err != nil {
		return nil, err
	}
	if len(queries) == 0 {
		return nil, fmt.Errorf("no query defined for number panel")
	}

	_, err = panel.AddQueries(queries...)
	if err != nil {
		return nil, err
	}

	return panel, nil
}

func textPanelFromResourceData(panelInfo map[string]interface{}) (*v2.Panels, error) {
	content := panelInfo["content"].(string)
	panel := &v2.Panels{
		ID:                    0,
		Name:                  panelInfo["name"].(string),
		Description:           "",
		Type:                  v2.PanelTypeText,
		MarkdownSource:        &content,
		PanelTitleVisible:     panelInfo["visible_title"].(bool),
		TextAutosized:         panelInfo["autosize_text"].(bool),
		TransparentBackground: panelInfo["transparent_background"].(bool),
	}
	_, err := panel.WithLayout(panelInfo["pos_x"].(int), panelInfo["pos_y"].(int), panelInfo["width"].(int), panelInfo["height"].(int))
	if err != nil {
		return nil, err
	}

	return panel, nil
}

func formatFromResourceData(queryInfo map[string]interface{}) *v2.Format {
	formatData, ok := queryInfo["format"]
	if !ok {
		return nil
	}

	formatSet := formatData.(*schema.Set).List()

	if len(formatSet) == 0 {
		return nil
	}

	fields := formatSet[0].(map[string]interface{})
	format := &v2.Format{}

	if inputFormat, ok := fields["input_format"].(string); ok {
		format.InputFormat = &inputFormat
	}

	if displayFormat, ok := fields["display_format"].(string); ok {
		format.DisplayFormat = &displayFormat
	}

	if decimals, ok := fields["decimals"].(int); ok {
		format.Decimals = &decimals
	}

	if yAxis, ok := fields["y_axis"].(string); ok {
		format.YAxis = &yAxis
	}

	if nullValue, ok := fields["null_value_display_mode"].(string); ok {
		format.NullValueDisplayMode = &nullValue
	}

	if minInterval, ok := fields["min_interval"].(string); ok {
		format.MinInterval = &minInterval
	}

	return format
}

func queriesFromResourceData(panelInfo map[string]interface{}, panel *v2.Panels) (newQueries []*v2.AdvancedQueries, err error) {
	var display v2.DisplayInfo

	for _, queryItr := range panelInfo["query"].(*schema.Set).List() {
		queryInfo := queryItr.(map[string]interface{})

		displayInfo := queryInfo["display_info"].(*schema.Set).List()
		if len(displayInfo) > 0 {
			dip := displayInfo[0].(map[string]interface{})
			display.DisplayName = dip["display_name"].(string)
			display.TimeSeriesDisplayNameTemplate = dip["time_series_display_name_template"].(string)
			display.Type = dip["type"].(string)
		} else {
			display.DisplayName = ""
			display.TimeSeriesDisplayNameTemplate = ""
			display.Type = "lines"
		}

		promqlQuery := v2.NewPromqlQuery(queryInfo["promql"].(string), panel, display)

		format := formatFromResourceData(queryInfo)

		switch queryInfo["unit"].(string) {
		case "percent":
			promqlQuery.WithPercentFormat(format)
		case "data":
			promqlQuery.WithDataFormat(format)
		case "data rate":
			promqlQuery.WithDataRateFormat(format)
		case "number":
			promqlQuery.WithNumberFormat(format)
		case "number rate":
			promqlQuery.WithNumberRateFormat(format)
		case "time":
			promqlQuery.WithTimeFormat(format)
		default:
			return nil, fmt.Errorf("unsupported query format unit: %s", queryInfo["unit"])
		}

		newQueries = append(newQueries, promqlQuery)
	}
	return
}

func dashboardToResourceData(dashboard *v2.Dashboard, data *schema.ResourceData) (err error) {
	_ = data.Set("name", dashboard.Name)
	_ = data.Set("description", dashboard.Description)
	_ = data.Set("public", dashboard.Public)
	_ = data.Set("public_token", dashboard.PublicToken)

	var panels []map[string]interface{}
	for i, panel := range dashboard.Panels {
		panelsData := data.Get("panel").(*schema.Set).List()
		panelData := map[string]interface{}{}
		if len(panelsData) > i {
			panelData = panelsData[i].(map[string]interface{})
		}

		dPanel, err := panelToResourceData(panel, dashboard.Layout, panelData)
		if err != nil {
			return err
		}
		panels = append(panels, dPanel)
	}
	_ = data.Set("panel", panels)

	var scopes []map[string]interface{}
	for _, scope := range dashboard.ScopeExpressionList {
		dScope, err := scopeToResourceData(scope)
		if err != nil {
			return err
		}
		scopes = append(scopes, dScope)
	}
	_ = data.Set("scope", scopes)
	_ = data.Set("version", dashboard.Version)

	var shares []map[string]interface{}
	for _, share := range dashboard.SharingSettings {
		dShare, err := shareToResourceData(share)
		if err != nil {
			return err
		}
		shares = append(shares, dShare)
	}
	_ = data.Set("share", shares)

	return nil
}

func shareToResourceData(share *v2.SharingOptions) (map[string]interface{}, error) {
	res := map[string]interface{}{
		"role": share.Role,
		"member": []map[string]interface{}{{
			"type": share.Member.Type,
			"id":   share.Member.ID,
		}},
	}
	return res, nil
}

func scopeToResourceData(scope *v2.ScopeExpressionList) (map[string]interface{}, error) {
	res := map[string]interface{}{
		"metric": scope.Operand,
	}

	if len(scope.Value) > 0 {
		res["value"] = scope.Value
		res["comparator"] = scope.Operator
	}

	if scope.IsVariable && scope.DisplayName != "" {
		res["variable"] = scope.DisplayName
	}

	return res, nil
}

func scopeFromResourceData(data *schema.ResourceData) ([]*v2.ScopeExpressionList, error) {
	scopes := []*v2.ScopeExpressionList{}
	for _, scopeItr := range data.Get("scope").(*schema.Set).List() {
		scopeInfo := (scopeItr).(map[string]interface{})

		scope := &v2.ScopeExpressionList{}
		scope.Operand = cast.ToString(scopeInfo["metric"])
		scope.Value = []string{}
		comparator := cast.ToString(scopeInfo["comparator"])
		value := cast.ToStringSlice(scopeInfo["value"])
		if comparator != "" {
			scope.Operator = comparator
			if len(value) == 0 {
				return nil, errors.New(`"value" field is required if the comparator is set up`)
			}
			if scope.Operator != "in" && scope.Operator != "notIn" && len(value) > 1 {
				return nil, errors.New(`"value" can only contain 1 value if the "comparator" is not "in" and "notIn"`)
			}
			scope.Value = value
		}
		variable := cast.ToString(scopeInfo["variable"])
		if variable != "" {
			scope.DisplayName = variable
			scope.IsVariable = true
			if scope.Operator == "" {
				scope.Operator = "in"
			}
		} else if comparator == "" || len(value) == 0 {
			return nil, errors.New(`"comparator" and "value" must be set if "variable" is not set`)
		}

		scopes = append(scopes, scope)
	}
	return scopes, nil
}

func panelToResourceData(panel *v2.Panels, layout []*v2.Layout, panelData map[string]interface{}) (map[string]interface{}, error) {
	var panelLayout *v2.Layout

	for _, l := range layout {
		if l.PanelID == panel.ID {
			panelLayout = l
		}
	}
	if panelLayout == nil {
		return nil, fmt.Errorf("inconsistent layout for dashboard trying to find panel ID: %d", panel.ID)
	}

	switch panel.Type {
	case v2.PanelTypeTimechart:
		return timechartPanelToResourceData(panel, panelLayout, panelData)
	case v2.PanelTypeNumber:
		return numberPanelToResourceData(panel, panelLayout, panelData)
	case v2.PanelTypeText:
		return textPanelToResourceData(panel, panelLayout)
	default:
		return nil, fmt.Errorf("unsupported panel type %s", panel.Type)
	}
}

func timechartPanelToResourceData(panel *v2.Panels, panelLayout *v2.Layout, panelData map[string]interface{}) (map[string]interface{}, error) {
	queries, err := queriesToResourceData(panel.AdvancedQueries, panelData)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"pos_x":       panelLayout.X,
		"pos_y":       panelLayout.Y,
		"width":       panelLayout.W,
		"height":      panelLayout.H,
		"name":        panel.Name,
		"description": panel.Description,
		"type":        "timechart",
		"query":       queries,
		"legend":      legendConfigurationToResourceData(panel.LegendConfiguration, panelData),
	}, nil
}

func legendConfigurationToResourceData(legend *v2.LegendConfiguration, panelData map[string]interface{}) []map[string]interface{} {

	legendData := panelData["legend"]
	// If legend is not defined in the user configuration and the dashboard legend is the same as the default one
	// we don't set the legend in the resource data to avoid drifts
	if legendData == nil || len(legendData.(*schema.Set).List()) == 0 {
		if *legend == *defaultLegendConfiguration() {
			return nil
		}
	}

	return []map[string]interface{}{{
		"enabled":      legend.Enabled,
		"show_current": legend.ShowCurrent,
		"position":     legend.Position,
		"layout":       legend.Layout,
	}}
}

func numberPanelToResourceData(panel *v2.Panels, panelLayout *v2.Layout, panelData map[string]interface{}) (map[string]interface{}, error) {
	queries, err := queriesToResourceData(panel.AdvancedQueries, panelData)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"pos_x":       panelLayout.X,
		"pos_y":       panelLayout.Y,
		"width":       panelLayout.W,
		"height":      panelLayout.H,
		"name":        panel.Name,
		"description": panel.Description,
		"type":        "number",
		"query":       queries,
	}, nil
}

func textPanelToResourceData(panel *v2.Panels, panelLayout *v2.Layout) (map[string]interface{}, error) {
	return map[string]interface{}{
		"pos_x":                  panelLayout.X,
		"pos_y":                  panelLayout.Y,
		"width":                  panelLayout.W,
		"height":                 panelLayout.H,
		"name":                   panel.Name,
		"content":                panel.MarkdownSource,
		"type":                   "text",
		"visible_title":          panel.PanelTitleVisible,
		"autosize_text":          panel.TextAutosized,
		"transparent_background": panel.TransparentBackground,
	}, nil
}

func queriesToResourceData(advancedQueries []*v2.AdvancedQueries, panelsData map[string]interface{}) ([]map[string]interface{}, error) {
	var queries []map[string]interface{}
	for queryIndex, query := range advancedQueries {
		unit := ""
		var defaultFormat v2.Format
		switch query.Format.Unit {
		case v2.FormatUnitPercentage:
			unit = "percent"
			defaultFormat = v2.NewPercentFormat()
		case v2.FormatUnitData:
			unit = "data"
			defaultFormat = v2.NewDataFormat()
		case v2.FormatUnitDataRate:
			unit = "data rate"
			defaultFormat = v2.NewDataRateFormat()
		case v2.FormatUnitNumber:
			unit = "number"
			defaultFormat = v2.NewNumberFormat()
		case v2.FormatUnitNumberRate:
			unit = "number rate"
			defaultFormat = v2.NewNumberRateFormat()
		case v2.FormatUnitTime:
			unit = "time"
			defaultFormat = v2.NewTimeFormat()
		default:
			return nil, fmt.Errorf("unsupported query format unit: %s", query.Format.Unit)
		}

		q := map[string]interface{}{
			"unit":   unit,
			"promql": query.Query,
		}

		if query.DisplayInfo.DisplayName != "" || query.DisplayInfo.TimeSeriesDisplayNameTemplate != "" {
			q["display_info"] = []map[string]interface{}{{
				"display_name":                      query.DisplayInfo.DisplayName,
				"time_series_display_name_template": query.DisplayInfo.TimeSeriesDisplayNameTemplate,
				"type":                              query.DisplayInfo.Type,
			}}
		}

		q["format"] = []map[string]interface{}{{
			"decimals":                query.Format.Decimals,
			"display_format":          query.Format.DisplayFormat,
			"input_format":            query.Format.InputFormat,
			"min_interval":            query.Format.MinInterval,
			"null_value_display_mode": query.Format.NullValueDisplayMode,
			"y_axis":                  query.Format.YAxis,
		}}

		queriesData := panelsData["query"]
		queryData := map[string]interface{}{}
		if queriesData != nil && queriesData.(*schema.Set).Len() > queryIndex {
			queryData = queriesData.(*schema.Set).List()[queryIndex].(map[string]interface{})
		}
		formatData := queryData["format"]

		// If format is not defined in the user configuration and the dashboard format is the same as the default one
		// we don't set the format in the resource data to avoid drifts
		if formatData == nil || formatData.(*schema.Set).Len() == 0 {
			if reflect.DeepEqual(query.Format, defaultFormat) {
				q["format"] = nil
			}
		}

		queries = append(queries, q)
	}
	return queries, nil
}
