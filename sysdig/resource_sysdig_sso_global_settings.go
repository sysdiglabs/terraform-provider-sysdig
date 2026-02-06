package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSSOGlobalSettings() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSSOGlobalSettingsCreate,
		ReadContext:   resourceSysdigSSOGlobalSettingsRead,
		UpdateContext: resourceSysdigSSOGlobalSettingsUpdate,
		DeleteContext: resourceSysdigSSOGlobalSettingsDelete,
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
			"product": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"monitor", "secure"}, false),
			},
			"is_password_login_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceSysdigSSOGlobalSettingsCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	product := d.Get("product").(string)
	settings := ssoGlobalSettingsFromResourceData(d)

	_, err = client.UpdateSSOGlobalSettings(ctx, product, settings)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(product)

	return resourceSysdigSSOGlobalSettingsRead(ctx, d, m)
}

func resourceSysdigSSOGlobalSettingsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	product := d.Id()

	settings, err := client.GetSSOGlobalSettings(ctx, product)
	if err != nil {
		if err == v2.ErrSSOGlobalSettingsNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = ssoGlobalSettingsToResourceData(settings, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSSOGlobalSettingsUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	product := d.Id()
	settings := ssoGlobalSettingsFromResourceData(d)

	_, err = client.UpdateSSOGlobalSettings(ctx, product, settings)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSysdigSSOGlobalSettingsRead(ctx, d, m)
}

func resourceSysdigSSOGlobalSettingsDelete(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
	return nil
}

func ssoGlobalSettingsFromResourceData(d *schema.ResourceData) *v2.SSOGlobalSettings {
	return &v2.SSOGlobalSettings{
		IsPasswordLoginEnabled: d.Get("is_password_login_enabled").(bool),
	}
}

func ssoGlobalSettingsToResourceData(settings *v2.SSOGlobalSettings, d *schema.ResourceData) error {
	// Product may not be returned in API response; use ID (which is the product name)
	product := settings.Product
	if product == "" {
		product = d.Id()
	}
	if err := d.Set("product", product); err != nil {
		return err
	}
	if err := d.Set("is_password_login_enabled", settings.IsPasswordLoginEnabled); err != nil {
		return err
	}

	return nil
}
