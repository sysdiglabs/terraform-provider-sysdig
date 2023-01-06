package sysdig

import (
	"context"
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigMonitorCloudAccountProvider() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigCloudAccountProviderCreate,
		DeleteContext: resourceSysdigCloudAccountProviderDelete,
		ReadContext:   resourceSysdigCloudAccountProviderRead,
		UpdateContext: resourceSysdigCloudAccountProviderUpdate,
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
			"platform": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"additional_options": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSysdigCloudAccountProviderCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	provider := providerFromResourceData(data)

	providerCreated, err := client.CreateCustomerProviderKey(ctx, &provider)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(providerCreated.Id))

	return nil
}

func resourceSysdigCloudAccountProviderDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteCustomerProviderKeyById(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigCloudAccountProviderRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	provider, err := client.GetCustomerProviderKeyById(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = providerToResourceData(data, provider)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigCloudAccountProviderUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	provider := providerFromResourceData(data)

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCustomerProviderKey(ctx, id, &provider)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func providerFromResourceData(data *schema.ResourceData) monitor.CustomerProviderKey {
	return monitor.CustomerProviderKey{
		Platform:          data.Get("platform").(string),
		IntegrationType:   data.Get("integration_type").(string),
		AdditionalOptions: data.Get("additional_options").(string),
		Credentials: monitor.CustomerProviderCredentials{
			AccountId: data.Get("account_id").(string),
		},
	}
}

func providerToResourceData(data *schema.ResourceData, provider *monitor.CustomerProviderKey) error {
	err := data.Set("platform", provider.Platform)
	if err != nil {
		return err
	}

	err = data.Set("integration_type", provider.IntegrationType)
	if err != nil {
		return err
	}

	err = data.Set("additional_options", provider.AdditionalOptions)
	if err != nil {
		return err
	}

	err = data.Set("account_id", provider.Credentials.AccountId)
	if err != nil {
		return err
	}

	return nil
}
