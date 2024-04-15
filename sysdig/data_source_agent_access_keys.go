package sysdig

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigAgentAccessKey() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigAgentAccessKeyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
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
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"date_disabled": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_key": {
				Type:     schema.TypeString,
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

	agentKeyId := d.Get("id").(int)
	agentAccessKey, err := client.GetAgentAccessKeyByID(ctx, strconv.Itoa(agentKeyId))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(agentAccessKey.ID))
	_ = d.Set("reservation", agentAccessKey.Reservation)
	_ = d.Set("limit", agentAccessKey.Limit)
	_ = d.Set("team_id", agentAccessKey.TeamID)
	_ = d.Set("metadata", agentAccessKey.Metadata)
	_ = d.Set("enabled", agentAccessKey.Enabled)
	_ = d.Set("date_disabled", agentAccessKey.DateDisabled)
	_ = d.Set("date_created", agentAccessKey.DateCreated)
	_ = d.Set("access_key", agentAccessKey.AgentAccessKey)

	return nil
}
