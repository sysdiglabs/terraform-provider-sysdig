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
			SchemaNameKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaDescriptionKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaMonitorPermKey: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SchemaSecurePermKey: {
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

	name := d.Get(SchemaNameKey).(string)

	customRole, err := client.GetCustomRoleByName(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(customRole.ID))
	_ = d.Set(SchemaNameKey, customRole.Name)
	_ = d.Set(SchemaDescriptionKey, customRole.Description)
	_ = d.Set(SchemaMonitorPermKey, customRole.MonitorPermissions)
	_ = d.Set(SchemaSecurePermKey, customRole.SecurePermissions)

	return nil
}
