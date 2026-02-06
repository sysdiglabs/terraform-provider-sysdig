package sysdig

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSSOGroupMapping() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSSOGroupMappingCreate,
		ReadContext:   resourceSysdigSSOGroupMappingRead,
		UpdateContext: resourceSysdigSSOGroupMappingUpdate,
		DeleteContext: resourceSysdigSSOGroupMappingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
			teamMaps := diff.Get("team_map").([]any)
			if len(teamMaps) > 0 {
				teamMap := teamMaps[0].(map[string]any)
				isForAllTeams := teamMap["is_for_all_teams"].(bool)
				teamIDs := teamMap["team_ids"].([]any)
				if !isForAllTeams && len(teamIDs) == 0 {
					return fmt.Errorf("team_ids must be set when is_for_all_teams is false")
				}
			}

			return nil
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"standard_team_role": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"custom_team_role_id"},
				AtLeastOneOf:  []string{"standard_team_role", "custom_team_role_id"},
			},
			"custom_team_role_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"standard_team_role"},
				AtLeastOneOf:  []string{"standard_team_role", "custom_team_role_id"},
			},
			"is_admin": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"team_map": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_for_all_teams": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"team_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      32767,
				ValidateFunc: validation.IntBetween(1, 32767),
			},
		},
	}
}

func resourceSysdigSSOGroupMappingCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	gm := ssoGroupMappingFromResourceData(d)

	created, err := client.CreateSSOGroupMapping(ctx, gm)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(created.ID))

	return resourceSysdigSSOGroupMappingRead(ctx, d, m)
}

func resourceSysdigSSOGroupMappingRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	gm, err := client.GetSSOGroupMapping(ctx, id)
	if err != nil {
		if err == v2.ErrSSOGroupMappingNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = ssoGroupMappingToResourceData(gm, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSSOGroupMappingUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	gm := ssoGroupMappingFromResourceData(d)

	_, err = client.UpdateSSOGroupMapping(ctx, id, gm)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSysdigSSOGroupMappingRead(ctx, d, m)
}

func resourceSysdigSSOGroupMappingDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteSSOGroupMapping(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ssoGroupMappingFromResourceData(d *schema.ResourceData) *v2.SSOGroupMapping {
	gm := &v2.SSOGroupMapping{
		GroupName: d.Get("group_name").(string),
		IsAdmin:   d.Get("is_admin").(bool),
		Weight:    d.Get("weight").(int),
	}

	if v, ok := d.GetOk("standard_team_role"); ok {
		gm.StandardTeamRole = v.(string)
	}

	if v, ok := d.GetOk("custom_team_role_id"); ok {
		gm.CustomTeamRoleID = v.(int)
	}

	teamMaps := d.Get("team_map").([]any)
	if len(teamMaps) > 0 {
		teamMap := teamMaps[0].(map[string]any)
		teamIDsInterface := teamMap["team_ids"].([]any)
		teamIDs := make([]int, len(teamIDsInterface))
		for i, id := range teamIDsInterface {
			teamIDs[i] = id.(int)
		}
		gm.TeamMap = &v2.SSOGroupMappingTeamMap{
			IsForAllTeams: teamMap["is_for_all_teams"].(bool),
			TeamIDs:       teamIDs,
		}
	}

	return gm
}

func ssoGroupMappingToResourceData(gm *v2.SSOGroupMapping, d *schema.ResourceData) error {
	if err := d.Set("group_name", gm.GroupName); err != nil {
		return err
	}
	if err := d.Set("standard_team_role", gm.StandardTeamRole); err != nil {
		return err
	}
	if err := d.Set("custom_team_role_id", gm.CustomTeamRoleID); err != nil {
		return err
	}
	if err := d.Set("is_admin", gm.IsAdmin); err != nil {
		return err
	}
	if err := d.Set("weight", gm.Weight); err != nil {
		return err
	}

	if gm.TeamMap != nil {
		teamMap := map[string]any{
			"is_for_all_teams": gm.TeamMap.IsForAllTeams,
			"team_ids":         gm.TeamMap.TeamIDs,
		}
		if err := d.Set("team_map", []map[string]any{teamMap}); err != nil {
			return err
		}
	}

	return nil
}
