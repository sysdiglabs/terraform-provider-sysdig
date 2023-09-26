package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureOrganization() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureOrganizationCreate,
		DeleteContext: resourceSysdigSecureOrganizationDelete,
		ReadContext:   resourceSysdigSecureOrganizationRead,
		UpdateContext: resourceSysdigSecureOrganizationUpdate,
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
			SchemaIDKey: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			SchemaManagementAccountId: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaCloudProviderType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{cloudauth.Provider_PROVIDER_AWS.String(), cloudauth.Provider_PROVIDER_GCP.String(), cloudauth.Provider_PROVIDER_AZURE.String()}, false),
			},
		},
	}
}

func getSecureOrganizationClient(c SysdigClients) (v2.OrganizationSecureInterface, error) {
	return c.sysdigSecureClientV2()
}

func resourceSysdigSecureOrganizationCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getSecureOrganizationClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	org := secureOrganizationFromResourceData(data)

	orgCreated, err := client.CreateOrganizationSecure(ctx, &org)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(orgCreated.Id)

	return nil
}

func resourceSysdigSecureOrganizationDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getSecureOrganizationClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteOrganizationSecure(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureOrganizationRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getSecureOrganizationClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	org, err := client.GetOrganizationSecure(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = secureOrganizationToResourceData(data, org)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureOrganizationUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getSecureOrganizationClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	org := secureOrganizationFromResourceData(data)

	_, err = client.UpdateOrganizationSecure(ctx, data.Id(), &org)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func secureOrganizationFromResourceData(data *schema.ResourceData) v2.OrganizationSecure {
	return v2.OrganizationSecure{
		CloudOrganization: cloudauth.CloudOrganization{
			ManagementAccountId: data.Get(SchemaManagementAccountId).(string),
			Provider:            cloudauth.Provider(cloudauth.Provider_value[data.Get(SchemaCloudProviderType).(string)]),
		},
	}
}

func secureOrganizationToResourceData(data *schema.ResourceData, org *v2.OrganizationSecure) error {
	err := data.Set(SchemaCloudProviderId, org.ProviderId)
	if err != nil {
		return err
	}

	err = data.Set(SchemaCloudProviderType, org.Provider)
	if err != nil {
		return err
	}

	err = data.Set(SchemaIDKey, org.Id)
	if err != nil {
		return err
	}

	return nil
}
