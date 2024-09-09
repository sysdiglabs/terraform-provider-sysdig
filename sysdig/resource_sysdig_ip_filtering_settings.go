package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigIPFilteringSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceSysdigIPFilteringSettingsRead,
		CreateContext: resourceSysdigIPFilteringSettingsCreate,
		UpdateContext: resourceSysdigIPFilteringSettingsUpdate,
		DeleteContext: resourceSysdigIPFilteringSettingsDelete,
		Schema: map[string]*schema.Schema{
			"ip_filtering_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceSysdigIPFilteringSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	ipFilteringSettings, err := client.GetIPFilteringSettings(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	err = ipFilteringSettingsToResourceData(ipFilteringSettings, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigIPFilteringSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	d.SetId("ip_filtering_settings_id") // It's singleton resource so we use a fixed ID

	return updateIPFilteringSettings(ctx, d, m)
}

func resourceSysdigIPFilteringSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return updateIPFilteringSettings(ctx, d, m)
}

func resourceSysdigIPFilteringSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func updateIPFilteringSettings(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	ipFiltersSettings, err := ipFilteringSettingsFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateIPFilteringSettings(ctx, ipFiltersSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigIPFilteringSettingsRead(ctx, d, m)

	return nil
}

func ipFilteringSettingsToResourceData(ipFiltersSettings *v2.IPFiltersSettings, d *schema.ResourceData) error {
	err := d.Set("ip_filtering_enabled", ipFiltersSettings.IPFilteringEnabled)
	if err != nil {
		return err
	}

	return nil
}

func ipFilteringSettingsFromResourceData(d *schema.ResourceData) (*v2.IPFiltersSettings, error) {
	return &v2.IPFiltersSettings{
		IPFilteringEnabled: d.Get("ip_filtering_enabled").(bool),
	}, nil
}
