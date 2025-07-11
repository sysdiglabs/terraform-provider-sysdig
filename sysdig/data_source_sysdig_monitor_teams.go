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
					},
				},
			},
		},
	}
}

func dataSourceSysdigMonitorTeamsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getMonitorTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	teams, err := client.ListTeams(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var result []map[string]any
	for _, team := range teams {
		result = append(result, map[string]any{
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
