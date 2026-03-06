package sysdig

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigBuiltinRole() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigBuiltinRoleRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			SchemaNameKey: {
				Type:     schema.TypeString,
				Required: true,
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

func dataSourceSysdigBuiltinRoleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(SchemaNameKey).(string)

	builtinRole, err := client.GetBuiltinRole(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	err = d.Set(SchemaNameKey, builtinRole.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaMonitorPermKey, builtinRole.MonitorPermissions)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaSecurePermKey, builtinRole.SecurePermissions)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
