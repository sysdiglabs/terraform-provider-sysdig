package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceSysdigMonitorTeam2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSysdigMonitorTeamCreate2,
		UpdateContext: resourceSysdigMonitorTeamUpdate2,
		ReadContext:   resourceSysdigMonitorTeamRead2,
		DeleteContext: resourceSysdigMonitorTeamDelete2,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func getMonitorTeam2Client(c SysdigClients) (v2.TeamInterface2, error) {
	var client v2.TeamInterface2
	var err error
	switch c.GetClientType() {
	case IBMMonitor:
		client, err = c.ibmMonitorClient()
		if err != nil {
			return nil, err
		}
	default:
		client, err = c.sysdigMonitorClientV2()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func resourceSysdigMonitorTeamCreate2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorTeam2Client(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	team, err := client.CreateTeam2(ctx, teamFromResourceData2(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(team.ID))
	resourceSysdigMonitorTeamRead2(ctx, d, meta)

	return nil
}

func resourceSysdigMonitorTeamRead2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getMonitorTeam2Client(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	t, err := client.GetTeamById2(ctx, id)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	err = d.Set("name", t.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorTeamUpdate2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorTeam2Client(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	t := teamFromResourceData2(d)
	t.ID, err = strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateTeam2(ctx, t)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigMonitorTeamRead2(ctx, d, meta)
	return nil
}

func resourceSysdigMonitorTeamDelete2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorTeam2Client(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteTeam2(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func teamFromResourceData2(d *schema.ResourceData) v2.Team2 {
	return v2.Team2{
		Name: d.Get("name").(string),
		Show: "host",
	}
}
