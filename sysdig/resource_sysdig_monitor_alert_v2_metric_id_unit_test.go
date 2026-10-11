//go:build unit

package sysdig

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
)

// When the API cannot resolve the metric descriptor, it omits "metric" and only returns the
// legacy "metricId" field. The resources must fall back to it to avoid a perpetual diff.
func TestAlertV2ReadFallsBackToLegacyMetricID(t *testing.T) {
	const legacyOnly = `{"id":1,"name":"test","config":{"metricId":"sysdig_container_up"}}`
	const both = `{"id":1,"name":"test","config":{"metric":{"id":"sysdig_container_cpu_used_percent"},"metricId":"sysdig_container_up"}}`

	cases := map[string]struct {
		resource    *schema.Resource
		updateState func(d *schema.ResourceData, body string) error
	}{
		"metric": {
			resource: resourceSysdigMonitorAlertV2Metric(),
			updateState: func(d *schema.ResourceData, body string) error {
				var a v2.AlertV2Metric
				if err := json.Unmarshal([]byte(body), &a); err != nil {
					return err
				}
				return updateAlertV2MetricState(d, &a)
			},
		},
		"downtime": {
			resource: resourceSysdigMonitorAlertV2Downtime(),
			updateState: func(d *schema.ResourceData, body string) error {
				var a v2.AlertV2Downtime
				if err := json.Unmarshal([]byte(body), &a); err != nil {
					return err
				}
				return updateAlertV2DowntimeState(d, &a)
			},
		},
		"change": {
			resource: resourceSysdigMonitorAlertV2Change(),
			updateState: func(d *schema.ResourceData, body string) error {
				var a v2.AlertV2Change
				if err := json.Unmarshal([]byte(body), &a); err != nil {
					return err
				}
				return updateAlertV2ChangeState(d, &a)
			},
		},
		"group_outlier": {
			resource: resourceSysdigMonitorAlertV2GroupOutlier(),
			updateState: func(d *schema.ResourceData, body string) error {
				var a v2.AlertV2GroupOutlier
				if err := json.Unmarshal([]byte(body), &a); err != nil {
					return err
				}
				return updateAlertV2GroupOutlierState(d, &a)
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			for body, want := range map[string]string{
				legacyOnly: "sysdig_container_up",
				both:       "sysdig_container_cpu_used_percent",
			} {
				d := tc.resource.TestResourceData()
				if err := tc.updateState(d, body); err != nil {
					t.Fatalf("updating state: %v", err)
				}
				if got := d.Get("metric").(string); got != want {
					t.Errorf("metric = %q, want %q (body %s)", got, want, body)
				}
			}
		})
	}
}
