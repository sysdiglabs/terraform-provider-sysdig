package sysdig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigMonitorTeams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSysdigMonitorTeamsRead,
		Schema: map[string]*schema.Schema{
			"teams": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"theme": {
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
				},
			},
		},
	}
}

func dataSourceSysdigMonitorTeamsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getMonitorTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	teams, err := client.ListTeams(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var result []map[string]interface{}
	for _, team := range teams {
		result = append(result, map[string]interface{}{
			"id":   team.ID,
			"name": team.Name,
		})
	}
	d.SetId("sysdig_monitor_teams")
	if err := d.Set("teams", result); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
