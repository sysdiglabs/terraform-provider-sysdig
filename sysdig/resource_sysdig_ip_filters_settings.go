package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigIPFiltersSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceSysdigIPFiltersSettingsRead,
		CreateContext: resourceSysdigIPFiltersSettingsCreate,
		UpdateContext: resourceSysdigIPFiltersSettingsUpdate,
		DeleteContext: resourceSysdigIPFiltersSettingsDelete,
		Schema: map[string]*schema.Schema{
			"ip_filtering_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceSysdigIPFiltersSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	ipFiltersSettings, err := client.GetIPFiltersSettings(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	err = ipFiltersSettingsToResourceData(ipFiltersSettings, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigIPFiltersSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	d.SetId("ip_filters_settings_id") // It's singleton resource so we use a fixed ID

	return updateIPFiltersSettings(ctx, d, m)
}

func resourceSysdigIPFiltersSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return updateIPFiltersSettings(ctx, d, m)
}

func resourceSysdigIPFiltersSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func updateIPFiltersSettings(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	ipFiltersSettings, err := ipFiltersSettingsFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateIPFiltersSettings(ctx, ipFiltersSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigIPFiltersSettingsRead(ctx, d, m)

	return nil
}

func ipFiltersSettingsToResourceData(ipFiltersSettings *v2.IPFiltersSettings, d *schema.ResourceData) error {
	err := d.Set("ip_filtering_enabled", ipFiltersSettings.IPFilteringEnabled)
	if err != nil {
		return err
	}

	return nil
}

func ipFiltersSettingsFromResourceData(d *schema.ResourceData) (*v2.IPFiltersSettings, error) {
	return &v2.IPFiltersSettings{
		IPFilteringEnabled: d.Get("ip_filtering_enabled").(bool),
	}, nil
}
