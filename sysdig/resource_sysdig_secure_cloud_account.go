package sysdig

import (
	"context"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

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
			"workload_identity_account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workload_identity_account_alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func getSecureCloudAccountClient(c SysdigClients) (v2.CloudAccountSecureInterface, error) {
	return c.sysdigSecureClientV2()
}

func resourceSysdigSecureCloudAccountCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureCloudAccountClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	cloudAccount, err := client.CreateCloudAccountSecure(ctx, cloudAccountFromResourceData(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudAccount.AccountID)
	_ = d.Set("account_id", cloudAccount.AccountID)
	_ = d.Set("cloud_provider", cloudAccount.Provider)
	_ = d.Set("alias", cloudAccount.Alias)
	_ = d.Set("role_enabled", cloudAccount.RoleAvailable)
	_ = d.Set("role_name", cloudAccount.RoleName)
	_ = d.Set("external_id", cloudAccount.ExternalID)
	_ = d.Set("workload_identity_account_id", cloudAccount.WorkLoadIdentityAccountID)
	_ = d.Set("workload_identity_account_alias", cloudAccount.WorkLoadIdentityAccountAlias)

	return nil
}

func resourceSysdigSecureCloudAccountRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureCloudAccountClient(meta.(SysdigClients))
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	cloudAccount, err := client.GetCloudAccountSecureByID(ctx, d.Id())
	if err != nil {
		d.SetId("")
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("account_id", cloudAccount.AccountID)
	_ = d.Set("cloud_provider", cloudAccount.Provider)
	_ = d.Set("alias", cloudAccount.Alias)
	_ = d.Set("role_enabled", cloudAccount.RoleAvailable)
	_ = d.Set("role_name", cloudAccount.RoleName)
	_ = d.Set("external_id", cloudAccount.ExternalID)
	_ = d.Set("workload_identity_account_id", cloudAccount.WorkLoadIdentityAccountID)
	_ = d.Set("workload_identity_account_alias", cloudAccount.WorkLoadIdentityAccountAlias)

	return nil
}

func resourceSysdigSecureCloudAccountUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureCloudAccountClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCloudAccountSecure(ctx, d.Id(), cloudAccountFromResourceData(d))
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureCloudAccountDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureCloudAccountClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteCloudAccountSecure(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(err)
	}
	return nil
}

func cloudAccountFromResourceData(d *schema.ResourceData) *v2.CloudAccountSecure {
	return &v2.CloudAccountSecure{
		AccountID:                    d.Get("account_id").(string),
		Provider:                     d.Get("cloud_provider").(string),
		Alias:                        d.Get("alias").(string),
		RoleAvailable:                d.Get("role_enabled").(bool),
		RoleName:                     d.Get("role_name").(string),
		WorkLoadIdentityAccountID:    d.Get("workload_identity_account_id").(string),
		WorkLoadIdentityAccountAlias: d.Get("workload_identity_account_alias").(string),
	}
}
