package sysdig

import (
	"context"
	"log"
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
	err = d.Set("identity", identity)
	if err != nil {
		log.Println("error assigning 'identity'")
	}

	provider := d.Get("cloud_provider")
	switch provider {
	case "aws", "gcp":
		// If identity is an ARN, attempt to extract certain fields
		parsedArn, err := arn.Parse(identity)
		if err == nil {
			err = d.Set("aws_account_id", parsedArn.AccountID)
			if err != nil {
				log.Println("error assigning 'aws_account_id'")
			}

			if parsedArn.Service == "iam" && strings.HasPrefix(parsedArn.Resource, "role/") {
				err = d.Set("aws_role_name", strings.TrimPrefix(parsedArn.Resource, "role/"))
				if err != nil {
					log.Println("error assigning 'aws_role_name'")
				}
			}
		}
	case "azure":
		// If identity is an Azure tenantID/clientID, separate into each part
		tenantID, clientID, err := parseAzureCreds(identity)
		if err == nil {
			err = d.Set("azure_tenant_id", tenantID)
			if err != nil {
				log.Println("error assigning 'azure_tenant_id'")
			}
			err = d.Set("azure_client_id", clientID)
			if err != nil {
				log.Println("error assigning 'azure_client_id'")
			}
		}
	}
	return nil
}
