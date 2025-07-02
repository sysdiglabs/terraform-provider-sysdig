package sysdig

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigMonitorTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSysdigMonitorTeamRead,
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
			"prometheus_remote_write_metrics_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"can_use_sysdig_capture": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"can_see_infrastructure_events": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"can_use_aws_data": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"can_use_agent_cli": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default_team": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enable_ibm_platform_metrics": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ibm_platform_metrics": {
				Type:     schema.TypeString,
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
			"entrypoint": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"selection": {
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
		},
	}
}

func dataSourceSysdigMonitorTeamRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getMonitorTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	team, err := client.GetTeamByID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(team.ID))
	_ = d.Set("name", team.Name)
	_ = d.Set("theme", team.Theme)
	_ = d.Set("description", team.Description)
	_ = d.Set("scope_by", team.Show)
	_ = d.Set("filter", team.Filter)
	_ = d.Set("can_use_sysdig_capture", team.CanUseSysdigCapture)
	_ = d.Set("can_see_infrastructure_events", team.CanUseCustomEvents)
	_ = d.Set("can_use_aws_data", team.CanUseAwsMetrics)
	_ = d.Set("can_use_agent_cli", team.CanUseAgentCli)
	_ = d.Set("default_team", team.DefaultTeam)
	_ = d.Set("user_roles", userMonitorRolesToSet(team.UserRoles))
	_ = d.Set("entrypoint", entrypointToSet(team.EntryPoint))
	_ = d.Set("version", team.Version)

	var ibmPlatformMetrics *string
	var prometheusRemoteWrite *string
	if team.NamespaceFilters != nil {
		ibmPlatformMetrics = team.NamespaceFilters.IBMPlatformMetrics
		prometheusRemoteWrite = team.NamespaceFilters.PrometheusRemoteWrite
	}
	_ = d.Set("enable_ibm_platform_metrics", team.CanUseBeaconMetrics)
	_ = d.Set("ibm_platform_metrics", ibmPlatformMetrics)
	_ = d.Set("prometheus_remote_write_metrics_filter", prometheusRemoteWrite)

	return nil
}
