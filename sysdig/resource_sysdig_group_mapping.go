package sysdig

import (
	"context"
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigGroupMapping() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext:   resourceSysdigGroupMappingRead,
		CreateContext: resourceSysdigGroupMappingCreate,
		UpdateContext: resourceSysdigGroupMappingUpdate,
		DeleteContext: resourceSysdigGroupMappingDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team_map": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_teams": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"team_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
		},
	}
}

func resourceSysdigGroupMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	groupMapping, err := client.GetGroupMapping(ctx, id)
	if err != nil {
		d.SetId("")
		return nil
	}

	groupMappingToResourceData(groupMapping, d)

	return nil
}

func resourceSysdigGroupMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error

	client, err := m.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return diag.FromErr(err)
	}

	groupMapping := groupMappingFromResourceData(d)
	groupMapping, err = client.CreateGroupMapping(ctx, groupMapping)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(groupMapping.ID))

	return nil
}

func resourceSysdigGroupMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error

	client, err := m.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return diag.FromErr(err)
	}

	groupMapping := groupMappingFromResourceData(d)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	groupMapping.ID = id
	groupMapping, err = client.UpdateGroupMapping(ctx, groupMapping, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigGroupMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteGroupMapping(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func groupMappingFromResourceData(d *schema.ResourceData) *common.GroupMapping {
	return &common.GroupMapping{
		GroupName: d.Get("group_name").(string),
		Role:      d.Get("role").(string),
		TeamMap:   teamMapFromResourceData(d),
	}
}

func teamMapFromResourceData(d *schema.ResourceData) *common.TeamMap {
	teamMap := d.Get("team_map").(*schema.Set).List()[0].(map[string]interface{})
	teamIDsInterface := teamMap["team_ids"].([]interface{})
	teamIDs := make([]int, len(teamIDsInterface))
	for i, teamID := range teamIDsInterface {
		teamIDs[i] = teamID.(int)
	}

	return &common.TeamMap{
		AllTeams: teamMap["all_teams"].(bool),
		TeamIDs:  teamIDs,
	}
}

func teamMapToResourceData(teamMap *common.TeamMap) map[string]interface{} {
	return map[string]interface{}{
		"all_teams": teamMap.AllTeams,
		"team_ids":  teamMap.TeamIDs,
	}
}

func groupMappingToResourceData(groupMapping *common.GroupMapping, d *schema.ResourceData) error {
	_ = d.Set("group_name", groupMapping.GroupName)
	_ = d.Set("role", groupMapping.Role)
	_ = d.Set("team_map", []map[string]interface{}{teamMapToResourceData(groupMapping.TeamMap)})

	return nil
}
