package sysdig

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecureConnectionRead,
		Schema: map[string]*schema.Schema{
			"secure_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sysdig Secure URL basepath to where backend requests will be sent",
			},
			"secure_api_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Sysdig Secure authentication api token",
			},
		},
	}
}

func dataSourceSecureConnectionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	endpoint, err := meta.(SysdigClients).GetSecureEndpoint()
	if err != nil {
		return diag.FromErr(err)
	}

	apiToken, err := meta.(SysdigClients).GetSecureAPIToken()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%x", sha256.Sum256(fmt.Appendf(nil, "%s,%s", endpoint, apiToken))))

	err = d.Set("secure_url", endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("secure_api_token", apiToken)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
