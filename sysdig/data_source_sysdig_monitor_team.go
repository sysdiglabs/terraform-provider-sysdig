package sysdig

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func createMonitorDataSourceTeamSchema() map[string]*schema.Schema {
	s := createBaseMonitorTeamSchema()
	applyOnSchema(s, func(s *schema.Schema) {
		s.Computed = true
	})

	s[TeamSchemaNameKey].Required = true
	s[TeamSchemaNameKey].Computed = false
	return s
}

func dataSourceSysdigMonitorTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSysdigMonitorTeamRead,
		Schema:      createMonitorDataSourceTeamSchema(),
	}
}

func dataSourceSysdigMonitorTeamRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getMonitorTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	team, err := client.GetTeamByName(ctx, data.Get(TeamSchemaNameKey).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(team.ID))
	err = teamMonitorToResourceData(data, clients, team)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
