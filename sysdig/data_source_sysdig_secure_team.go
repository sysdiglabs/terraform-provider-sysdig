package sysdig

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureTeamRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"theme": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scope_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use_sysdig_capture": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default_team": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"user_roles": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zone_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"all_zones": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSysdigSecureTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getSecureTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	team, err := client.GetTeamById(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(team.ID))
	_ = d.Set("name", team.Name)
	_ = d.Set("theme", team.Theme)
	_ = d.Set("description", team.Description)
	_ = d.Set("scope_by", team.Show)
	_ = d.Set("filter", team.Filter)
	_ = d.Set("use_sysdig_capture", team.CanUseSysdigCapture)
	_ = d.Set("default_team", team.DefaultTeam)
	_ = d.Set("user_roles", userSecureRolesToSet(team.UserRoles))
	_ = d.Set("version", team.Version)
	_ = d.Set("zone_ids", team.ZoneIDs)
	_ = d.Set("all_zones", team.AllZones)

	return nil
}
