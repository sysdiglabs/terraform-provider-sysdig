package sysdig

import (
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigTeamReadIBM(d *schema.ResourceData, t *v2.Team) {
	var ibmPlatformMetrics *string
	if t.NamespaceFilters != nil {
		ibmPlatformMetrics = t.NamespaceFilters.IBMPlatformMetrics
	}
	_ = d.Set("enable_ibm_platform_metrics", t.CanUseBeaconMetrics)
	_ = d.Set("ibm_platform_metrics", ibmPlatformMetrics)
}

func updateNamespaceFilters(filters *v2.NamespaceFilters, update v2.NamespaceFilters) *v2.NamespaceFilters {
	if filters == nil {
		filters = &v2.NamespaceFilters{}
	}

	if update.IBMPlatformMetrics != nil {
		filters.IBMPlatformMetrics = update.IBMPlatformMetrics
	}

	return filters
}

func teamFromResourceDataIBM(d *schema.ResourceData, t *v2.Team) {
	canUseBeaconMetrics := d.Get("enable_ibm_platform_metrics").(bool)
	t.CanUseBeaconMetrics = &canUseBeaconMetrics

	if v, ok := d.GetOk("ibm_platform_metrics"); ok {
		metrics := v.(string)
		t.NamespaceFilters = updateNamespaceFilters(t.NamespaceFilters, v2.NamespaceFilters{
			IBMPlatformMetrics: &metrics,
		})
	}
}
