package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigGroupMapping() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext:   resourceSysdigGroupMappingRead,
		CreateContext: resourceSysdigGroupMappingCreate,
		UpdateContext: resourceSysdigGroupMappingUpdate,
		DeleteContext: resourceSysdigGroupMappingDelete,
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
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
			},
			"system_role": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceSysdigGroupMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	groupMapping, err := client.GetGroupMapping(ctx, id)
	if err != nil {
		if err == v2.GroupMappingNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = groupMappingToResourceData(groupMapping, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigGroupMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error

	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	groupMapping := groupMappingFromResourceData(d)
	groupMapping, err = client.CreateGroupMapping(ctx, groupMapping)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(groupMapping.ID))

	resourceSysdigGroupMappingRead(ctx, d, m)

	return nil
}

func resourceSysdigGroupMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error

	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	groupMapping := groupMappingFromResourceData(d)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	groupMapping.ID = id
	_, err = client.UpdateGroupMapping(ctx, groupMapping, id)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigGroupMappingRead(ctx, d, m)

	return nil
}

func resourceSysdigGroupMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
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

func groupMappingFromResourceData(d *schema.ResourceData) *v2.GroupMapping {
	return &v2.GroupMapping{
		GroupName:  d.Get("group_name").(string),
		Role:       d.Get("role").(string),
		SystemRole: d.Get("system_role").(string),
		TeamMap:    teamMapFromResourceData(d),
		Weight:     d.Get("weight").(int),
	}
}

func teamMapFromResourceData(d *schema.ResourceData) *v2.TeamMap {
	teamMap := d.Get("team_map").(*schema.Set).List()[0].(map[string]interface{})
	teamIDsInterface := teamMap["team_ids"].([]interface{})
	teamIDs := make([]int, len(teamIDsInterface))
	for i, teamID := range teamIDsInterface {
		teamIDs[i] = teamID.(int)
	}

	return &v2.TeamMap{
		AllTeams: teamMap["all_teams"].(bool),
		TeamIDs:  teamIDs,
	}
}

func teamMapToResourceData(teamMap *v2.TeamMap) map[string]interface{} {
	return map[string]interface{}{
		"all_teams": teamMap.AllTeams,
		"team_ids":  teamMap.TeamIDs,
	}
}

func groupMappingToResourceData(groupMapping *v2.GroupMapping, d *schema.ResourceData) error {
	err := d.Set("group_name", groupMapping.GroupName)
	if err != nil {
		return err
	}
	err = d.Set("role", groupMapping.Role)
	if err != nil {
		return err
	}
	err = d.Set("system_role", groupMapping.SystemRole)
	if err != nil {
		return err
	}
	err = d.Set("team_map", []map[string]interface{}{teamMapToResourceData(groupMapping.TeamMap)})
	if err != nil {
		return err
	}

	return nil
}
