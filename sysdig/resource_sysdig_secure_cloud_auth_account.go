package sysdig

import (
	"context"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureCloudauthAccount() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureCloudauthAccountCreate,
		UpdateContext: resourceSysdigSecureCloudauthAccountUpdate,
		ReadContext:   resourceSysdigSecureCloudauthAccountRead,
		DeleteContext: resourceSysdigSecureCloudauthAccountDelete,
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
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_provider_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{cloudauth.Provider_PROVIDER_AWS.String(), cloudauth.Provider_PROVIDER_GCP.String(), cloudauth.Provider_PROVIDER_AZURE.String()}, false),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func getSecureCloudauthAccountClient(client SysdigClients) (v2.CloudauthAccountSecureInterface, error) {
	return client.sysdigSecureClientV2()
}

func resourceSysdigSecureCloudauthAccountCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountClient((meta.(SysdigClients)))

	if err != nil {
		return diag.FromErr(err)
	}

	cloudauthAccount, err := client.CreateCloudauthAccountSecure(ctx, cloudauthAccountFromResourceData(data))

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(cloudauthAccount.Id)

	return nil
}

func resourceSysdigSecureCloudauthAccountRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountClient(meta.(SysdigClients))

	if err != nil {
		return diag.FromErr(err)
	}

	cloudauthAccount, err := client.GetCloudauthAccountSecure(ctx, data.Id())

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(err)
	}

	err = cloudauthAccountToResourceData(data, cloudauthAccount)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureCloudauthAccountUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountClient(meta.(SysdigClients))

	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCloudauthAccountSecure(ctx, data.Id(), cloudauthAccountFromResourceData(data))

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureCloudauthAccountDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountClient(meta.(SysdigClients))

	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteCloudauthAccountSecure(ctx, data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func cloudauthAccountFromResourceData(data *schema.ResourceData) *v2.CloudauthAccountSecure {
	return &v2.CloudauthAccountSecure{
		Id:         data.Get("account_id").(string),
		Enabled:    data.Get("enabled").(bool),
		ProviderId: data.Get("cloud_provider_id").(string),
		Provider:   data.Get("cloud_provider_type").(cloudauth.Provider),
	}
}

func cloudauthAccountToResourceData(data *schema.ResourceData, cloudAccount *v2.CloudauthAccountSecure) error {
	err := data.Set("account_id", cloudAccount.Id)

	if err != nil {
		return err
	}

	err = data.Set("enabled", cloudAccount.Enabled)

	if err != nil {
		return err
	}

	err = data.Set("cloud_provider_id", cloudAccount.ProviderId)

	if err != nil {
		return err
	}

	err = data.Set("cloud_provider_type", cloudAccount.Provider)

	if err != nil {
		return err
	}

	return nil
}
