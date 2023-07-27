package sysdig

import (
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func dataSourceSysdigSecureCustomRolePermissions() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: getDataSourceSysdigCustomRoleMonitorPermissionsRead(v2.SecureProduct),

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: dataSourceSysdigCustomRoleSchema(),
	}
}
