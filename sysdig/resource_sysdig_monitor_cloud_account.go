package sysdig

import (
	"context"
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigMonitorCloudAccount() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorCloudAccountCreate,
		DeleteContext: resourceSysdigMonitorCloudAccountDelete,
		ReadContext:   resourceSysdigMonitorCloudAccountRead,
		UpdateContext: resourceSysdigMonitorCloudAccountUpdate,
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
			"cloud_provider": {
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

func resourceSysdigMonitorCloudAccountCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudAccount := monitorCloudAccountFromResourceData(data)

	cloudAccountCreated, err := client.CreateCloudAccount(ctx, &cloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(cloudAccountCreated.Id))

	return nil
}

func resourceSysdigMonitorCloudAccountDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteCloudAccountById(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorCloudAccountRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudAccount, err := client.GetCloudAccountById(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = monitorCloudAccountToResourceData(data, cloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorCloudAccountUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudAccount := monitorCloudAccountFromResourceData(data)

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCloudAccount(ctx, id, &cloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func monitorCloudAccountFromResourceData(data *schema.ResourceData) monitor.CloudAccount {
	return monitor.CloudAccount{
		Platform:          data.Get("cloud_provider").(string),
		IntegrationType:   data.Get("integration_type").(string),
		AdditionalOptions: data.Get("additional_options").(string),
		Credentials: monitor.CloudAccountCredentials{
			AccountId: data.Get("account_id").(string),
		},
	}
}

func monitorCloudAccountToResourceData(data *schema.ResourceData, cloudAccount *monitor.CloudAccount) error {
	err := data.Set("cloud_provider", cloudAccount.Platform)
	if err != nil {
		return err
	}

	err = data.Set("integration_type", cloudAccount.IntegrationType)
	if err != nil {
		return err
	}

	err = data.Set("additional_options", cloudAccount.AdditionalOptions)
	if err != nil {
		return err
	}

	err = data.Set("account_id", cloudAccount.Credentials.AccountId)
	if err != nil {
		return err
	}

	return nil
}
