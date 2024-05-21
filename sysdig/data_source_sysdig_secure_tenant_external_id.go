package sysdig

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
	client, err := getSecureCloudAccountClient(meta.(SysdigClients))
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
