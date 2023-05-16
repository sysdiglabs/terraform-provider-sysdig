package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureManagedRuleset() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigManagedRulesetRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createPolicyDataSourceSchema(),
	}
}

func dataSourceSysdigManagedRulesetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	policyName := d.Get("name").(string)
	policyType := d.Get("type").(string)

	policies, _, err := client.GetPolicies(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var policy v2.Policy
	for _, existingPolicy := range policies {
		if existingPolicy.Name == policyName && existingPolicy.Type == policyType {
			if existingPolicy.IsDefault || existingPolicy.TemplateId == 0 {
				return diag.Errorf("policy is not a managed ruleset")
			}
			policy = existingPolicy
			break
		}
	}

	if policy.ID == 0 {
		return diag.Errorf("unable to find managed ruleset")
	}

	loadedPolicy, _, err := client.GetPolicyByID(ctx, policy.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	policyDataSourceToResourceData(loadedPolicy, d)

	return nil
}
