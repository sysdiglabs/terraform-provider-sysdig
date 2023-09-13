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
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_provider_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{cloudauth.Provider_PROVIDER_AWS.String(), cloudauth.Provider_PROVIDER_GCP.String(), cloudauth.Provider_PROVIDER_AZURE.String()}, false),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"components": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"instance": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cloudConnectorMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"trustedRoleMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"eventBridgeMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"servicePrincipalMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"webhookDatasourceMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cryptoKeyMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cloudLogsMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
	data.Set("cloud_provider_type", cloudauthAccount.Provider.String())

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
		CloudAccount: cloudauth.CloudAccount{
			Enabled:    data.Get("enabled").(bool),
			ProviderId: data.Get("cloud_provider_id").(string),
			Provider:   cloudauth.Provider(cloudauth.Provider_value[data.Get("cloud_provider_type").(string)]),
			Components: cloudauth.FeatureComponents(data.Get("components").([]map[string]interface{})),
		},
	}
}

func cloudauthAccountToResourceData(data *schema.ResourceData, cloudAccount *v2.CloudauthAccountSecure) error {
	err := data.Set("id", cloudAccount.Id)

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

	err = data.Set("cloud_provider_type", cloudAccount.Provider.String())

	if err != nil {
		return err
	}

	err = data.Set("components", cloudAccount.Components)

	if err != nil {
		return err
	}

	return nil
}
