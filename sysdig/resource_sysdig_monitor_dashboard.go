package sysdig

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/spf13/cast"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor/model"
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

	dashboardCreated, err := client.CreateDashboard(ctx, dashboard)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(dashboardCreated.ID))
	_ = data.Set("version", dashboardCreated.Version)

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

	_, err = client.UpdateDashboard(ctx, dashboard)
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

	err = dashboardToResourceData(dashboard, data)
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

func dashboardFromResourceData(data *schema.ResourceData) (dashboard *model.Dashboard, err error) {
	dashboard = model.NewDashboard(data.Get("name").(string), data.Get("description").(string)).AsPublic(data.Get("public").(bool))
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

	return dashboard, nil
}

func sharingFromResourceData(data *schema.ResourceData) (sharingSettings []*model.SharingOptions, err error) {
	for _, share := range data.Get("share").(*schema.Set).List() {
		shareInfo := share.(map[string]interface{})
		memberInfo := shareInfo["member"].(*schema.Set).List()[0].(map[string]interface{})
		sharingSettings = append(sharingSettings,
			&model.SharingOptions{
				Member: model.SharingMember{
					Type: memberInfo["type"].(string),
					ID:   memberInfo["id"].(int),
				},
				Role: shareInfo["role"].(string),
			})
	}
	return
}

func panelsFromResourceData(data *schema.ResourceData) (panels []*model.Panels, err error) {
	for _, panelItr := range data.Get("panel").(*schema.Set).List() {
		panelInfo := panelItr.(map[string]interface{})

		var panel *model.Panels
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

func timechartPanelFromResourceData(panelInfo map[string]interface{}) (*model.Panels, error) {
	panel := &model.Panels{
		ID:                     0,
		Name:                   panelInfo["name"].(string),
		Description:            panelInfo["description"].(string),
		Type:                   model.PanelTypeTimechart,
		ApplyScopeToAll:        false,
		ApplySegmentationToAll: false,
		AxesConfiguration: &model.AxesConfiguration{
			Bottom: model.Bottom{Enabled: true},
			Left: model.Left{
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
			Right: model.Right{
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
		LegendConfiguration: &model.LegendConfiguration{
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

func numberPanelFromResourceData(panelInfo map[string]interface{}) (*model.Panels, error) {
	panel := &model.Panels{
		ID:                     0,
		Name:                   panelInfo["name"].(string),
		Description:            panelInfo["description"].(string),
		Type:                   model.PanelTypeNumber,
		ApplyScopeToAll:        false,
		ApplySegmentationToAll: false,
		AxesConfiguration: &model.AxesConfiguration{
			Bottom: model.Bottom{Enabled: true},
			Left: model.Left{
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
			Right: model.Right{
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
		LegendConfiguration: &model.LegendConfiguration{
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
		NumberThresholds: &model.NumberThresholds{
			Values: []interface{}{}, // These values must be not nil in case of type number
			Base: model.Base{
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

func textPanelFromResourceData(panelInfo map[string]interface{}) (*model.Panels, error) {
	content := panelInfo["content"].(string)
	panel := &model.Panels{
		ID:                    0,
		Name:                  panelInfo["name"].(string),
		Description:           "",
		Type:                  model.PanelTypeText,
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

func queriesFromResourceData(panelInfo map[string]interface{}, panel *model.Panels) (newQueries []*model.AdvancedQueries, err error) {
	for _, queryItr := range panelInfo["query"].(*schema.Set).List() {
		queryInfo := queryItr.(map[string]interface{})

		promqlQuery := model.NewPromqlQuery(queryInfo["promql"].(string), panel)

		switch queryInfo["unit"].(string) {
		case "percent":
			promqlQuery.WithPercentFormat()
		case "data":
			promqlQuery.WithDataFormat()
		case "data rate":
			promqlQuery.WithDataRateFormat()
		case "number":
			promqlQuery.WithNumberFormat()
		case "number rate":
			promqlQuery.WithNumberRateFormat()
		case "time":
			promqlQuery.WithTimeFormat()
		default:
			return nil, fmt.Errorf("unsupported query format unit: %s", queryInfo["unit"])
		}

		newQueries = append(newQueries, promqlQuery)
	}
	return
}

func dashboardToResourceData(dashboard *model.Dashboard, data *schema.ResourceData) (err error) {
	_ = data.Set("name", dashboard.Name)
	_ = data.Set("description", dashboard.Description)
	_ = data.Set("public", dashboard.Public)
	_ = data.Set("public_token", dashboard.PublicToken)

	var panels []map[string]interface{}
	for _, panel := range dashboard.Panels {
		dPanel, err := panelToResourceData(panel, dashboard.Layout)
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

func shareToResourceData(share *model.SharingOptions) (map[string]interface{}, error) {
	res := map[string]interface{}{
		"role": share.Role,
		"member": []map[string]interface{}{{
			"type": share.Member.Type,
			"id":   share.Member.ID,
		}},
	}
	return res, nil
}

func scopeToResourceData(scope *model.ScopeExpressionList) (map[string]interface{}, error) {
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

func scopeFromResourceData(data *schema.ResourceData) ([]*model.ScopeExpressionList, error) {
	scopes := []*model.ScopeExpressionList{}
	for _, scopeItr := range data.Get("scope").(*schema.Set).List() {
		scopeInfo := (scopeItr).(map[string]interface{})

		scope := &model.ScopeExpressionList{}
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

func panelToResourceData(panel *model.Panels, layout []*model.Layout) (map[string]interface{}, error) {
	var panelLayout *model.Layout

	for _, l := range layout {
		if l.PanelID == panel.ID {
			panelLayout = l
		}
	}
	if panelLayout == nil {
		return nil, fmt.Errorf("inconsistent layout for dashboard trying to find panel ID: %d", panel.ID)
	}

	switch panel.Type {
	case model.PanelTypeTimechart:
		return timechartPanelToResourceData(panel, panelLayout)
	case model.PanelTypeNumber:
		return numberPanelToResourceData(panel, panelLayout)
	case model.PanelTypeText:
		return textPanelToResourceData(panel, panelLayout)
	default:
		return nil, fmt.Errorf("unsupported panel type %s", panel.Type)
	}
}

func timechartPanelToResourceData(panel *model.Panels, panelLayout *model.Layout) (map[string]interface{}, error) {
	queries, err := queriesToResourceData(panel.AdvancedQueries)
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
	}, nil
}

func numberPanelToResourceData(panel *model.Panels, panelLayout *model.Layout) (map[string]interface{}, error) {
	queries, err := queriesToResourceData(panel.AdvancedQueries)
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

func textPanelToResourceData(panel *model.Panels, panelLayout *model.Layout) (map[string]interface{}, error) {
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

func queriesToResourceData(advancedQueries []*model.AdvancedQueries) ([]map[string]interface{}, error) {
	var queries []map[string]interface{}
	for _, query := range advancedQueries {
		unit := ""
		switch query.Format.Unit {
		case model.FormatUnitPercentage:
			unit = "percent"
		case model.FormatUnitData:
			unit = "data"
		case model.FormatUnitDataRate:
			unit = "data rate"
		case model.FormatUnitNumber:
			unit = "number"
		case model.FormatUnitNumberRate:
			unit = "number rate"
		case model.FormatUnitTime:
			unit = "time"
		default:
			return nil, fmt.Errorf("unsupported query format unit: %s", query.Format.Unit)
		}

		queries = append(queries, map[string]interface{}{
			"unit":   unit,
			"promql": query.Query,
		})
	}
	return queries, nil
}
