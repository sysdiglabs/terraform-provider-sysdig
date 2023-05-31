package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigGroupMappingConfig() *schema.Resource {
	timeout := 5 * time.Minute
	return &schema.Resource{
		ReadContext:   resourceSysdigGroupMappingConfigRead,
		CreateContext: resourceSysdigGroupMappingConfigCreate,
		UpdateContext: resourceSysdigGroupMappingConfigUpdate,
		DeleteContext: resourceSysdigGroupMappingConfigDelete,
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
			"no_mapping_strategy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"different_team_same_role_strategy": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSysdigGroupMappingConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	groupMappingConfig, err := client.GetGroupMappingConfig(ctx)
	if err != nil {
		if err == v2.GroupMappingConfigNotFound {
			return nil
		}
		return diag.FromErr(err)
	}

	err = groupMappingConfigToResourceData(groupMappingConfig, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigGroupMappingConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	groupMappingConfig := groupMappingConfigFromResourceData(d)
	groupMappingConfig, err = client.CreateGroupMappingConfig(ctx, groupMappingConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("conflicts_resolution_strategies")

	resourceSysdigGroupMappingConfigRead(ctx, d, m)

	return nil
}

func resourceSysdigGroupMappingConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	groupMappingConfig := groupMappingConfigFromResourceData(d)
	groupMappingConfig, err = client.UpdateGroupMappingConfig(ctx, groupMappingConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO
	resourceSysdigGroupMappingConfigRead(ctx, d, m)

	return nil
}

func resourceSysdigGroupMappingConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func groupMappingConfigToResourceData(groupMappingConfig *v2.GroupMappingConfig, d *schema.ResourceData) error {
	err := d.Set("no_mapping_strategy", groupMappingConfig.NoMappingStrategy)
	if err != nil {
		return err
	}
	err = d.Set("different_team_same_role_strategy", groupMappingConfig.DifferentTeamSameRoleStrategy)

	return nil
}

func groupMappingConfigFromResourceData(d *schema.ResourceData) *v2.GroupMappingConfig {
	return &v2.GroupMappingConfig{
		NoMappingStrategy:             d.Get("no_mapping_strategy").(string),
		DifferentTeamSameRoleStrategy: d.Get("different_team_same_role_strategy").(string),
	}
}
