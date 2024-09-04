package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"access_key_id": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"additional_options": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func getMonitorCloudAccountClient(c SysdigClients) (v2.CloudAccountMonitorInterface, error) {
	return c.sysdigMonitorClientV2()
}

func resourceSysdigMonitorCloudAccountCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getMonitorCloudAccountClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	cloudAccount := monitorCloudAccountFromResourceData(data)

	cloudAccountCreated, err := client.CreateCloudAccountMonitor(ctx, &cloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(cloudAccountCreated.Id))

	return nil
}

func resourceSysdigMonitorCloudAccountDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getMonitorCloudAccountClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteCloudAccountMonitor(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorCloudAccountRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getMonitorCloudAccountClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudAccount, err := client.GetCloudAccountMonitor(ctx, id)
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
	client, err := getMonitorCloudAccountClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	cloudAccount := monitorCloudAccountFromResourceData(data)

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCloudAccountMonitor(ctx, id, &cloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func monitorCloudAccountFromResourceData(data *schema.ResourceData) v2.CloudAccountMonitor {
	return v2.CloudAccountMonitor{
		Platform:          data.Get("cloud_provider").(string),
		IntegrationType:   data.Get("integration_type").(string),
		AdditionalOptions: data.Get("additional_options").(string),
		Credentials: v2.CloudAccountCredentialsMonitor{
			AccountId:   data.Get("account_id").(string),
			RoleName:    data.Get("role_name").(string),
			SecretKey:   data.Get("secret_key").(string),
			AccessKeyId: data.Get("access_key_id").(string),
		},
	}
}

func monitorCloudAccountToResourceData(data *schema.ResourceData, cloudAccount *v2.CloudAccountMonitor) error {
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

	err = data.Set("role_name", cloudAccount.Credentials.RoleName)
	if err != nil {
		return err
	}

	err = data.Set("secret_key", cloudAccount.Credentials.SecretKey)
	if err != nil {
		return err
	}

	err = data.Set("access_key_id", cloudAccount.Credentials.AccessKeyId)
	if err != nil {
		return err
	}

	return nil
}
