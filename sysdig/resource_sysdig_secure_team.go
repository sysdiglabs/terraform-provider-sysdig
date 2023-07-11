package sysdig

import (
	"context"
	"fmt"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecureTeam() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureTeamCreate,
		UpdateContext: resourceSysdigSecureTeamUpdate,
		ReadContext:   resourceSysdigSecureTeamRead,
		DeleteContext: resourceSysdigSecureTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
			plan := diff.GetRawPlan().AsValueMap()
			zoneIDsPlan := plan[SchemaZonesIDsKey]
			allZonesPlan := plan[SchemaAllZones]

			var nonEmptyZoneIDs bool
			if !zoneIDsPlan.IsNull() && len(zoneIDsPlan.AsValueSlice()) > 0 {
				nonEmptyZoneIDs = true
			}

			if nonEmptyZoneIDs && allZonesPlan.True() {
				return fmt.Errorf("if %s is enabled, %s must be omitted", SchemaAllZones, SchemaZonesIDsKey)
			}

			return nil
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
				Default:  "#73A1F7",
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
				Default:  "container",
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
			"use_sysdig_capture": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"default_team": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			SchemaZonesIDsKey: {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			SchemaAllZones: {
				Optional: true,
				Type:     schema.TypeBool,
				Default:  false,
			},
		},
	}
}

func getSecureTeamClient(c SysdigClients) (v2.TeamInterface, error) {
	var client v2.TeamInterface
	var err error
	switch c.GetClientType() {
	case IBMSecure:
		client, err = c.ibmSecureClient()
		if err != nil {
			return nil, err
		}
	default:
		client, err = c.sysdigSecureClientV2()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func resourceSysdigSecureTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getSecureTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	team := secureTeamFromResourceData(d, clients.GetClientType())
	team.Products = []string{"SDS"}

	team, err = client.CreateTeam(ctx, team)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(team.ID))
	_ = d.Set("version", team.Version)
	resourceSysdigSecureTeamRead(ctx, d, meta)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigSecureTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getSecureTeamClient(clients)
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
	_ = d.Set("use_sysdig_capture", t.CanUseSysdigCapture)
	_ = d.Set("default_team", t.DefaultTeam)
	_ = d.Set("user_roles", userSecureRolesToSet(t.UserRoles))

	err = d.Set(SchemaZonesIDsKey, t.ZoneIDs)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaAllZones, t.AllZones)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigTeamReadIBM(d, &t)

	return nil
}

func userSecureRolesToSet(userRoles []v2.UserRoles) (res []map[string]interface{}) {
	for _, role := range userRoles {
		if role.Admin {
			continue // Admins are added by default, so skip them
		}
		roleMap := map[string]interface{}{
			"email": role.Email,
			"role":  role.Role,
		}
		res = append(res, roleMap)
	}
	return
}

func resourceSysdigSecureTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getSecureTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	t := secureTeamFromResourceData(d, clients.GetClientType())
	t.Products = []string{"SDS"}

	t.Version = d.Get("version").(int)
	t.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateTeam(ctx, t)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigSecureTeamRead(ctx, d, meta)
	return nil
}

func resourceSysdigSecureTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureTeamClient(meta.(SysdigClients))
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

func secureTeamFromResourceData(d *schema.ResourceData, clientType ClientType) v2.Team {
	canUseSysdigCapture := d.Get("use_sysdig_capture").(bool)
	canUseAwsMetrics := new(bool)
	allZones := d.Get(SchemaAllZones).(bool)
	t := v2.Team{
		Theme:               d.Get("theme").(string),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Show:                d.Get("scope_by").(string),
		Filter:              d.Get("filter").(string),
		CanUseSysdigCapture: &canUseSysdigCapture,
		CanUseAwsMetrics:    canUseAwsMetrics,
		DefaultTeam:         d.Get("default_team").(bool),
		AllZones:            allZones,
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

	zonesData := d.Get("zone_ids").([]interface{})
	t.ZoneIDs = make([]int, len(zonesData))
	for i, z := range zonesData {
		t.ZoneIDs[i] = z.(int)
	}

	teamFromResourceDataIBM(d, &t)

	return t
}
