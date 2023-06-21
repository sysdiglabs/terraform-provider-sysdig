package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func dataSourceSysdigSecurePosturePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSysdigSecurePosturePoliciesRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func getPosturePolicyClient(c SysdigClients) (v2.PosturePolicyInterface, error) {
	var client v2.PosturePolicyInterface
	var err error
	switch c.GetClientType() {
	case IBMSecure:
		client, err = c.ibmSecureClient()
		if err != nil {
			return nil, err
		}
	default:
		client, err = c.sysdigSecureClientV2()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func dataSourceSysdigSecurePosturePoliciesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getPosturePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListPosturePolicies(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	policies := make([]map[string]interface{}, len(resp))
	for i, p := range resp {
		policies[i] = map[string]interface{}{
			"id":   p.ID,
			"name": p.Name,
			"type": p.Type,
			"kind": p.Kind,
		}
	}

	d.SetId("policies")
	_ = d.Set("policies", policies)

	return nil
}
