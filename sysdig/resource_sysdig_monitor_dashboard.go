package sysdig

import (
	"context"
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor/model"
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

	dashboardCreated, err := client.CreateDashboard(ctx, dashboard)
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

	dashboard.PublicToken = data.Get("public_token").(string)

	panels, err := panelsFromResourceData(data)
	if err != nil {
		return nil, err
	}

	dashboard.AddPanels(panels...)
	return dashboard, nil
}

func panelsFromResourceData(data *schema.ResourceData) (panels []*model.Panels, err error) {
	for _, panelItr := range data.Get("panel").(*schema.Set).List() {
		panelInfo := panelItr.(map[string]interface{})

		var panelType model.PanelType
		switch panelInfo["type"].(string) {
		case "timechart":
			panelType = model.PanelTypeTimechart
		case "number":
			panelType = model.PanelTypeNumber
		default:
			return nil, fmt.Errorf("unsupported panel type %s", panelInfo["type"])
		}

		panel := model.NewPanel(panelInfo["name"].(string), panelInfo["description"].(string), panelType)

		_, err = panel.WithLayout(panelInfo["pos_x"].(int), panelInfo["pos_y"].(int), panelInfo["width"].(int), panelInfo["height"].(int))
		if err != nil {
			return nil, err
		}

		queries, err := queriesFromResourceData(panelInfo, panel)
		if err != nil {
			return nil, err
		}

		_, err = panel.AddQueries(queries...)
		if err != nil {
			return nil, err
		}

		panels = append(panels, panel)
	}
	return
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
				return fmt.Errorf("unsupported query format unit: %s", query.Format.Unit)
			}

			queries = append(queries, map[string]interface{}{
				"unit":   unit,
				"promql": query.Query,
			})
		}

		var panelLayout *model.Layout

		for _, layout := range dashboard.Layout {
			if layout.PanelID == panel.ID {
				panelLayout = layout
			}
		}

		var panelType string
		switch panel.Type {
		case model.PanelTypeTimechart:
			panelType = "timechart"
		case model.PanelTypeNumber:
			panelType = "number"
		default:
			return fmt.Errorf("unsupported panel type %s", panel.Type)
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
