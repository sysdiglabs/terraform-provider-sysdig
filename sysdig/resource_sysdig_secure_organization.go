package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
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
			"cloud_provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_provider_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"customer_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Optional: true,
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

/* TODO: Make sure these are right inputs for the api */
func secureOrganizationFromResourceData(data *schema.ResourceData) v2.OrganizationSecure {
	return v2.OrganizationSecure{
		Id:         data.Get("organization_id").(string),
		ProviderId: data.Get("cloud_provider_id").(string),
		Provider:   data.Get("cloud_provider_type").(int32),
		CustomerId: data.Get("customer_id").(uint64),
	}
}

func secureOrganizationToResourceData(data *schema.ResourceData, org *v2.OrganizationSecure) error {
	err := data.Set("cloud_provider_id", org.ProviderId)
	if err != nil {
		return err
	}

	err = data.Set("cloud_provider_type", org.Provider)
	if err != nil {
		return err
	}

	err = data.Set("customer_id", org.CustomerId)
	if err != nil {
		return err
	}

	err = data.Set("organization_id", org.Id)
	if err != nil {
		return err
	}

	return nil
}
