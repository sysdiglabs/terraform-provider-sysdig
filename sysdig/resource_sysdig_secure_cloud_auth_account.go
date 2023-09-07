package sysdig

import (
	"context"
	"strings"
	"time"

	draiosproto "github.com/draios/protorepo/cloudauth/go"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
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
			"cloud_provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_provider_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"aws", "gcp", "azure"}, false),
			},
			"cloud_provider_alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"customer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": {
				Type:     schema.TypeString,
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
		data.SetId("")
		return diag.FromErr(err)
	}

	cloudauthAccount, err := client.GetCloudauthAccountSecure(ctx, data.Id())

	if err != nil {
		data.SetId("")
		
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
		Id:            data.Get("account_id").(string),
		CustomerId:    data.Get("customer_id").(uint64),
		ProviderId:    data.Get("cloud_provider_id").(string),
		Provider:      data.Get("cloud_provider_type").(draiosproto.Provider),
		ProviderAlias: data.Get("cloud_provider_alias").(string),
	}
}

func cloudauthAccountToResourceData(data *schema.ResourceData, cloudAccount *v2.CloudauthAccountSecure) error {
	err := data.Set("account_id", cloudAccount.Id)

	if err != nil {
		return err
	}

	err = data.Set("customer_id", cloudAccount.CustomerId)

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

	err = data.Set("cloud_provider_alias", cloudAccount.ProviderAlias)

	if err != nil {
		return err
	}

	return nil
}
