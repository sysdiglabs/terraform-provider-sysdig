package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigMonitorTeam() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorTeamCreate,
		UpdateContext: resourceSysdigMonitorTeamUpdate,
		ReadContext:   resourceSysdigMonitorTeamRead,
		DeleteContext: resourceSysdigMonitorTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(5 * time.Minute), // Removing the team is for some reason slower.
		},

		Schema: map[string]*schema.Schema{
			"theme": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "#05C391",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope_by": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "host",
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_ibm_platform_metrics": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ibm_platform_metrics": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"can_use_sysdig_capture": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_see_infrastructure_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_use_aws_data": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"user_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ROLE_TEAM_STANDARD",
						},
					},
				},
			},
			"entrypoint": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Explore", "Dashboards", "Events", "Alerts", "Settings"}, false),
						},

						"selection": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"default_team": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func getMonitorTeamClient(c SysdigClients) (v2.TeamInterface, error) {
	var client v2.TeamInterface
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

func resourceSysdigMonitorTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getMonitorTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	team := teamFromResourceData(d, clients.GetClientType())
	team.Products = []string{"SDC"}

	team, err = client.CreateTeam(ctx, team)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(team.ID))
	resourceSysdigMonitorTeamRead(ctx, d, meta)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigMonitorTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getMonitorTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	t, err := client.GetTeamById(ctx, id)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_ = d.Set("version", t.Version)
	_ = d.Set("theme", t.Theme)
	_ = d.Set("name", t.Name)
	_ = d.Set("description", t.Description)
	_ = d.Set("scope_by", t.Show)
	_ = d.Set("filter", t.Filter)
	_ = d.Set("can_use_sysdig_capture", t.CanUseSysdigCapture)
	_ = d.Set("can_see_infrastructure_events", t.CanUseCustomEvents)
	_ = d.Set("can_use_aws_data", t.CanUseAwsMetrics)
	_ = d.Set("default_team", t.DefaultTeam)
	_ = d.Set("user_roles", userMonitorRolesToSet(t.UserRoles))
	_ = d.Set("entrypoint", entrypointToSet(t.EntryPoint))

	resourceSysdigTeamReadIBM(d, &t)

	return nil
}

func userMonitorRolesToSet(userRoles []v2.UserRoles) (res []map[string]interface{}) {
	for _, role := range userRoles {
		if role.Admin { // Admins are added by default, so skip them
			continue
		}

		roleMap := map[string]interface{}{
			"email": role.Email,
			"role":  role.Role,
		}
		res = append(res, roleMap)
	}
	return
}

func entrypointToSet(entrypoint *v2.EntryPoint) (res []map[string]interface{}) {
	if entrypoint == nil {
		return
	}

	entrypointMap := map[string]interface{}{
		"type":      entrypoint.Module,
		"selection": entrypoint.Selection,
	}
	return append(res, entrypointMap)
}

func resourceSysdigMonitorTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getMonitorTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	t := teamFromResourceData(d, clients.GetClientType())
	t.Products = []string{"SDC"}

	t.Version = d.Get("version").(int)
	t.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateTeam(ctx, t)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigMonitorTeamRead(ctx, d, meta)
	return nil
}

func resourceSysdigMonitorTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorTeamClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeleteTeam(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func teamFromResourceData(d *schema.ResourceData, clientType ClientType) v2.Team {
	canUseSysdigCapture := d.Get("can_use_sysdig_capture").(bool)
	canUseCustomEvents := d.Get("can_see_infrastructure_events").(bool)
	canUseAwsMetrics := d.Get("can_use_aws_data").(bool)
	canUseBeaconMetrics := false
	t := v2.Team{
		Theme:               d.Get("theme").(string),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Show:                d.Get("scope_by").(string),
		Filter:              d.Get("filter").(string),
		CanUseSysdigCapture: &canUseSysdigCapture,
		CanUseCustomEvents:  &canUseCustomEvents,
		CanUseAwsMetrics:    &canUseAwsMetrics,
		CanUseBeaconMetrics: &canUseBeaconMetrics,
		DefaultTeam:         d.Get("default_team").(bool),
	}

	userRoles := make([]v2.UserRoles, 0)
	for _, userRole := range d.Get("user_roles").(*schema.Set).List() {
		ur := userRole.(map[string]interface{})
		userRoles = append(userRoles, v2.UserRoles{
			Email: ur["email"].(string),
			Role:  ur["role"].(string),
		})
	}
	t.UserRoles = userRoles

	t.EntryPoint = &v2.EntryPoint{}
	t.EntryPoint.Module = d.Get("entrypoint.0.type").(string)
	if val, ok := d.GetOk("entrypoint.0.selection"); ok {
		t.EntryPoint.Selection = val.(string)
	}

	teamFromResourceDataIBM(d, &t)

	return t
}
