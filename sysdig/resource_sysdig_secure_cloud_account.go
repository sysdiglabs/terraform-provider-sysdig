package sysdig

import (
	"context"
	"log"
	"time"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureCloudAccount() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureCloudAccountCreate,
		UpdateContext: resourceSysdigSecureCloudAccountUpdate,
		ReadContext:   resourceSysdigSecureCloudAccountRead,
		DeleteContext: resourceSysdigSecureCloudAccountDelete,
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
			"cloud_provider": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"aws", "gcp", "azure"}, false),
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"role_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "SysdigCloudBench",
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSysdigSecureCloudAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudAccount, err := client.CreateCloudAccount(ctx, cloudAccountFromResourceData(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudAccount.AccountID)
	err = d.Set("account_id", cloudAccount.AccountID)
	if err != nil {
		log.Println("error asigning 'account_id'")
	}
	err = d.Set("cloud_provider", cloudAccount.Provider)
	if err != nil {
		log.Println("error asigning 'cloud_provider'")
	}
	err = d.Set("alias", cloudAccount.Alias)
	if err != nil {
		log.Println("error asigning 'alias'")
	}
	err = d.Set("role_enabled", cloudAccount.RoleAvailable)
	if err != nil {
		log.Println("error asigning 'role_enabled'")
	}
	err = d.Set("role_name", cloudAccount.RoleName)
	if err != nil {
		log.Println("error asigning 'role_name'")
	}
	err = d.Set("external_id", cloudAccount.ExternalID)
	if err != nil {
		log.Println("error asigning 'external_id'")
	}

	return nil
}

func resourceSysdigSecureCloudAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	cloudAccount, err := client.GetCloudAccountById(ctx, d.Id())
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	err = d.Set("account_id", cloudAccount.AccountID)
	if err != nil {
		log.Println("error asigning 'account_id'")
	}

	err = d.Set("cloud_provider", cloudAccount.Provider)
	if err != nil {
		log.Println("error asigning 'cloud_provider'")
	}

	err = d.Set("alias", cloudAccount.Alias)
	if err != nil {
		log.Println("error asigning 'alias'")
	}

	err = d.Set("role_enabled", cloudAccount.RoleAvailable)
	if err != nil {
		log.Println("error asigning 'role_enabled'")
	}

	err = d.Set("role_name", cloudAccount.RoleName)
	if err != nil {
		log.Println("error asigning 'role_name'")
	}

	err = d.Set("external_id", cloudAccount.ExternalID)
	if err != nil {
		log.Println("error asigning 'external_id'")
	}

	return nil
}

func resourceSysdigSecureCloudAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCloudAccount(ctx, d.Id(), cloudAccountFromResourceData(d))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureCloudAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteCloudAccount(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func cloudAccountFromResourceData(d *schema.ResourceData) *secure.CloudAccount {
	return &secure.CloudAccount{
		AccountID:     d.Get("account_id").(string),
		Provider:      d.Get("cloud_provider").(string),
		Alias:         d.Get("alias").(string),
		RoleAvailable: d.Get("role_enabled").(bool),
		RoleName:      d.Get("role_name").(string),
	}
}
