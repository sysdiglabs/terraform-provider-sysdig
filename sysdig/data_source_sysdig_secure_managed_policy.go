package sysdig

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureManagedPolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigManagedPolicyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createPolicyDataSourceSchema(),
	}
}

func dataSourceSysdigManagedPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	policyName := d.Get("name").(string)
	policyType := d.Get("type").(string)

	policy, err := getManagedPolicy(ctx, client, policyName, policyType)

	loadedPolicy, _, err := client.GetPolicyByID(ctx, policy.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	policyDataSourceToResourceData(loadedPolicy, d)

	return nil
}
