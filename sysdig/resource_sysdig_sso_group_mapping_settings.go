package sysdig

import (
	"context"
	"fmt"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSSOGroupMappingSettings() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSSOGroupMappingSettingsCreate,
		ReadContext:   resourceSysdigSSOGroupMappingSettingsRead,
		UpdateContext: resourceSysdigSSOGroupMappingSettingsUpdate,
		DeleteContext: resourceSysdigSSOGroupMappingSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
			strategy := diff.Get("no_mapping_strategy").(string)
			redirectURL, hasRedirectURL := diff.GetOk("no_mappings_error_redirect_url")

			if strategy == "NO_MAPPINGS_ERROR_REDIRECT" && (!hasRedirectURL || redirectURL.(string) == "") {
				return fmt.Errorf("no_mappings_error_redirect_url must be set when no_mapping_strategy is NO_MAPPINGS_ERROR_REDIRECT")
			}

			if strategy != "NO_MAPPINGS_ERROR_REDIRECT" && hasRedirectURL && redirectURL.(string) != "" {
				return fmt.Errorf("no_mappings_error_redirect_url can only be set when no_mapping_strategy is NO_MAPPINGS_ERROR_REDIRECT")
			}

			return nil
		},
		Schema: map[string]*schema.Schema{
			"no_mapping_strategy": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"UNAUTHORIZED",
					"DEFAULT_TEAM_DEFAULT_ROLE",
					"NO_MAPPINGS_ERROR_REDIRECT",
				}, false),
			},
			"different_roles_same_team_strategy": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"UNAUTHORIZED",
					"HIGHEST_ROLE",
					"LOWEST_ROLE",
				}, false),
			},
			"no_mappings_error_redirect_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 2048),
			},
		},
	}
}

func resourceSysdigSSOGroupMappingSettingsCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	settings := ssoGroupMappingSettingsFromResourceData(d)

	_, err = client.UpdateSSOGroupMappingSettings(ctx, settings)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("sso_group_mapping_settings")

	return resourceSysdigSSOGroupMappingSettingsRead(ctx, d, m)
}

func resourceSysdigSSOGroupMappingSettingsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	settings, err := client.GetSSOGroupMappingSettings(ctx)
	if err != nil {
		if err == v2.ErrSSOGroupMappingSettingsNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = ssoGroupMappingSettingsToResourceData(settings, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSSOGroupMappingSettingsUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	settings := ssoGroupMappingSettingsFromResourceData(d)

	_, err = client.UpdateSSOGroupMappingSettings(ctx, settings)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSysdigSSOGroupMappingSettingsRead(ctx, d, m)
}

func resourceSysdigSSOGroupMappingSettingsDelete(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
	return nil
}

func ssoGroupMappingSettingsFromResourceData(d *schema.ResourceData) *v2.SSOGroupMappingSettings {
	settings := &v2.SSOGroupMappingSettings{
		NoMappingStrategy:              d.Get("no_mapping_strategy").(string),
		DifferentRolesSameTeamStrategy: d.Get("different_roles_same_team_strategy").(string),
	}

	if v, ok := d.GetOk("no_mappings_error_redirect_url"); ok {
		settings.NoMappingsErrorRedirectURL = v.(string)
	}

	return settings
}

func ssoGroupMappingSettingsToResourceData(settings *v2.SSOGroupMappingSettings, d *schema.ResourceData) error {
	if err := d.Set("no_mapping_strategy", settings.NoMappingStrategy); err != nil {
		return err
	}
	if err := d.Set("different_roles_same_team_strategy", settings.DifferentRolesSameTeamStrategy); err != nil {
		return err
	}
	if err := d.Set("no_mappings_error_redirect_url", settings.NoMappingsErrorRedirectURL); err != nil {
		return err
	}

	return nil
}
