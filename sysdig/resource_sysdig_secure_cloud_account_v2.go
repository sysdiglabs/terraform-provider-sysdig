package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureCloudAccountV2() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		// CreateContext: resourceSysdigSecureCloudAccountCreate,
		// UpdateContext: resourceSysdigSecureCloudAccountUpdate,
		// ReadContext:   resourceSysdigSecureCloudAccountRead,
		// DeleteContext: resourceSysdigSecureCloudAccountDelete,
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
				ValidateFunc: validation.StringInSlice([]string{"gcp"}, false),
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

func resourceSysdigSecureCloudAccountV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSysdigSecureCloudAccountV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSysdigSecureCloudAccountV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSysdigSecureCloudAccountV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func cloudAccountFromResourceDataV2(d *schema.ResourceData) *v2.CloudAccountSecure {
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
