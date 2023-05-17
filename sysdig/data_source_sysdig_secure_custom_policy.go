package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureCustomPolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureCustomPolicyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createPolicyDataSourceSchema(),
	}
}

func dataSourceSysdigSecureCustomPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return commonDataSourceSecurePolicyRead(ctx, d, meta, "custom policy", isCustomPolicy)
}

func isCustomPolicy(policy v2.Policy) bool {
	return !policy.IsDefault && policy.TemplateId == 0
}
