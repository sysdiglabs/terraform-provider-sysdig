package sysdig

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/arn"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getSecureOnboardingClient(c SysdigClients) (v2.OnboardingSecureInterface, error) {
	return c.sysdigSecureClientV2()
}

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
			"azure_tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"azure_service_principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Retrieves the information of a resource form the file and loads it in Terraform
func dataSourceSysdigSecureTrustedCloudIdentityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureOnboardingClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	identity, err := client.GetTrustedCloudIdentitySecure(ctx, d.Get("cloud_provider").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(identity)
	_ = d.Set("identity", identity)

	provider := d.Get("cloud_provider")
	switch provider {
	case "aws", "gcp":
		// If identity is an ARN, attempt to extract certain fields
		parsedArn, err := arn.Parse(identity)
		if err == nil {
			_ = d.Set("aws_account_id", parsedArn.AccountID)
			if parsedArn.Service == "iam" && strings.HasPrefix(parsedArn.Resource, "role/") {
				_ = d.Set("aws_role_name", strings.TrimPrefix(parsedArn.Resource, "role/"))
			}
		}
	case "azure":
		// If identity is an Azure tenantID/clientID, separate into each part
		tenantID, spID, err := parseAzureCreds(identity)
		if err == nil {
			_ = d.Set("azure_tenant_id", tenantID)
			_ = d.Set("azure_service_principal_id", spID)

		}
	}
	return nil
}

func dataSourceSysdigSecureTenantExternalID() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureTenantExternalIDRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Retrieves the information of a resource form the file and loads it in Terraform
func dataSourceSysdigSecureTenantExternalIDRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureOnboardingClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	externalId, err := client.GetTenantExternalIDSecure(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(externalId)
	err = d.Set("external_id", externalId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
