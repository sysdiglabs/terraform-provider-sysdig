package sysdig

import (
	"context"
	"time"

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

func resourceSysdigSecureCloudauthAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSysdigSecureCloudauthAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSysdigSecureCloudauthAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSysdigSecureCloudauthAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func cloudauthAccountFromResourceData(data *schema.ResourceData) *v2.CloudAccountSecure {
	return &v2.CloudAccountSecure{
		AccountID:                    data.Get("account_id").(string),
		Provider:                     data.Get("cloud_provider").(string),
		Alias:                        data.Get("alias").(string),
		RoleAvailable:                data.Get("role_enabled").(bool),
		RoleName:                     data.Get("role_name").(string),
		WorkLoadIdentityAccountID:    data.Get("workload_identity_account_id").(string),
		WorkLoadIdentityAccountAlias: data.Get("workload_identity_account_alias").(string),
	}
}
