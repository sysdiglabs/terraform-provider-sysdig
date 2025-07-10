package sysdig

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigCustomRole() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext:   resourceSysdigCustomRoleRead,
		CreateContext: resourceSysdigCustomRoleCreate,
		UpdateContext: resourceSysdigCustomRoleUpdate,
		DeleteContext: resourceSysdigCustomRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},
		Schema: map[string]*schema.Schema{
			SchemaNameKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaDescriptionKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaPermissionsKey: {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						SchemaMonitorPermKey: {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						SchemaSecurePermKey: {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceSysdigCustomRoleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	customRole, err := client.GetCustomRoleByID(ctx, id)
	if err != nil {
		if err == v2.ErrCustomRoleNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = customRoleToResourceData(customRole, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigCustomRoleCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var err error

	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	customRole, err := customRoleFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	customRole, err = client.CreateCustomRole(ctx, customRole)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(customRole.ID))

	resourceSysdigCustomRoleRead(ctx, d, m)

	return nil
}

func resourceSysdigCustomRoleUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var err error

	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	customRole, err := customRoleFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	customRole.ID = id
	_, err = client.UpdateCustomRole(ctx, customRole, id)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigCustomRoleRead(ctx, d, m)

	return nil
}

func resourceSysdigCustomRoleDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteCustomRole(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func customRoleFromResourceData(d *schema.ResourceData) (*v2.CustomRole, error) {
	schemaPermissions, ok := d.Get(SchemaPermissionsKey).(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("cast permissions to set resuts in an error")
	}
	p := permission{
		schemaPermissions,
	}
	return &v2.CustomRole{
		Name:               d.Get(SchemaNameKey).(string),
		Description:        d.Get(SchemaDescriptionKey).(string),
		MonitorPermissions: p.readMonitorPermissions(),
		SecurePermissions:  p.readSecurePermissions(),
	}, nil
}

type permission struct {
	s *schema.Set
}

func (p *permission) readPermissions(product string) []string {
	permissionsMap := p.s.List()[0].(map[string]any)
	permissionsInterface := permissionsMap[product].(*schema.Set).List()
	permissions := make([]string, len(permissionsInterface))
	for i, permission := range permissionsInterface {
		permissions[i] = permission.(string)
	}
	return permissions
}

func (p *permission) readSecurePermissions() []string {
	return p.readPermissions(SchemaSecurePermKey)
}

func (p *permission) readMonitorPermissions() []string {
	return p.readPermissions(SchemaMonitorPermKey)
}

func customRoleToResourceData(customRole *v2.CustomRole, d *schema.ResourceData) error {
	err := d.Set(SchemaNameKey, customRole.Name)
	if err != nil {
		return err
	}
	err = d.Set(SchemaDescriptionKey, customRole.Description)
	if err != nil {
		return err
	}
	err = d.Set(SchemaPermissionsKey, []map[string]any{
		permissionsToResourceData(customRole.MonitorPermissions, customRole.SecurePermissions),
	})
	if err != nil {
		return err
	}
	return nil
}

func permissionsToResourceData(monitorPermissions []string, securePermissions []string) map[string]any {
	return map[string]any{
		SchemaMonitorPermKey: monitorPermissions,
		SchemaSecurePermKey:  securePermissions,
	}
}
