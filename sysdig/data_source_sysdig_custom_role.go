package sysdig

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceSysdigCustomRole() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigCustomRoleRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"monitor_permissions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"secure_permissions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSysdigCustomRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	customRole, err := client.GetCustomRoleByName(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(customRole.ID))
	_ = d.Set("name", customRole.Name)
	_ = d.Set("description", customRole.Description)
	_ = d.Set("monitor_permissions", customRole.MonitorPermissions)
	_ = d.Set("secure_permissions", customRole.SecurePermissions)

	return nil
}
