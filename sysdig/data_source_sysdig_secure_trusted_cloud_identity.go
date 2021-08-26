package sysdig

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceSysdigSecureTrustedCloudIdentity() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureTrustedCloudIdentityRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"cloud_provider": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"aws", "gcp", "azure"}, false),
			},
			"identity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Retrieves the information of a resource form the file and loads it in Terraform
func dataSourceSysdigSecureTrustedCloudIdentityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	identity, err := client.GetTrustedCloudIdentity(ctx, d.Get("cloud_provider").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(identity)
	d.Set("identity", identity)

	// If identity is an ARN, attempt to extract certain fields
	parsedArn, err := arn.Parse(identity)
	if err == nil {
		d.Set("aws_account_id", parsedArn.AccountID)

		if parsedArn.Service == "iam" && strings.HasPrefix(parsedArn.Resource, "role/") {
			d.Set("aws_role_name", strings.TrimPrefix(parsedArn.Resource, "role/"))
		}
	}

	return nil
}
