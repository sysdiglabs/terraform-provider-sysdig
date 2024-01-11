package sysdig

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func dataSourceSysdigAgentAccessKey() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigAgentAccessKeyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"agent_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reservation": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"team_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"team_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"agents_connected": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// Retrieves the information of a resource form the file and loads it in Terraform
func dataSourceSysdigAgentAccessKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).commonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	agentKeyId := d.Get("agent_key").(string)

	agentAccessKey, err := client.GetAgentAccessKeyById(ctx, agentKeyId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(agentAccessKey.AgentAccessKeyId)
	_ = d.Set("reservation", agentAccessKey.Reservation)
	_ = d.Set("limit", agentAccessKey.Limit)
	_ = d.Set("team_id", agentAccessKey.TeamID)
	_ = d.Set("metadata", agentAccessKey.Metadata)
	_ = d.Set("team_name", agentAccessKey.TeamName)
	_ = d.Set("enabled", agentAccessKey.Enabled)
	_ = d.Set("agents_connected", agentAccessKey.AgentsConnected)

	return nil
}
