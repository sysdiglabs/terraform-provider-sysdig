package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureManagedPolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureManagedPolicyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createPolicyDataSourceSchema(),
	}
}

func dataSourceSysdigSecureManagedPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return commonDataSourceSecurePolicyRead(ctx, d, meta, "managed policy", isManagedPolicy)
}

func isManagedPolicy(policy v2.Policy) bool {
	return policy.IsDefault
}
