package sysdig

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigCustomRoleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		SchemaRequestedPermKey: {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		SchemaEnrichedPermKey: {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func getDataSourceSysdigCustomRoleMonitorPermissionsRead(product v2.Product) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client, err := m.(SysdigClients).sysdigCommonClientV2()
		if err != nil {
			return diag.FromErr(err)
		}
		rp := d.Get(SchemaRequestedPermKey).([]interface{})

		rps := readPermissions(rp)
		dependencies, err := client.GetPermissionsDependencies(ctx, product, rps)
		if err != nil {
			return diag.FromErr(err)
		}
		ps := make([]string, len(dependencies))
		for i, dependency := range dependencies {
			ps[i] = dependency.PermissionAuthority
			ps = append(ps, dependency.Dependencies...)
		}

		cdefChecksum := sha256.Sum256([]byte(strings.Join(rps, ",")))
		d.SetId(fmt.Sprintf("%x", cdefChecksum))
		_ = d.Set(SchemaEnrichedPermKey, ps)

		return nil
	}
}

func readPermissions(rp []interface{}) []string {
	permissions := make([]string, len(rp))
	for i, permission := range rp {
		permissions[i] = permission.(string)
	}
	return permissions
}
