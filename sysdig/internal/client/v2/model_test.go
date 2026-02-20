package v2

import (
	"encoding/json"
	"testing"
)

func TestAlertV2MetricUnmarshalLegacyMetricID(t *testing.T) {
	// Legacy API response: has "metricId" but no "metric.id"
	legacyJSON := `{
		"alert": {
			"id": 1553848,
			"name": "MysqlMateo",
			"type": "MANUAL",
			"config": {
				"metricId": "mysql_global_status_threads_connected",
				"groupAggregation": "avg",
				"timeAggregation": "avg",
				"threshold": 125.0,
				"conditionOperator": ">",
				"noDataBehaviour": "DO_NOTHING",
				"range": 60,
				"duration": 0
			}
		}
	}`

	var wrapper alertV2MetricWrapper
	err := json.Unmarshal([]byte(legacyJSON), &wrapper)
	if err != nil {
		t.Fatalf("failed to unmarshal legacy alert JSON: %v", err)
	}

	alert := wrapper.Alert

	if alert.Config.MetricID != "mysql_global_status_threads_connected" {
		t.Errorf("expected MetricID to be %q, got %q", "mysql_global_status_threads_connected", alert.Config.MetricID)
	}

	// metric.id should be empty since the legacy response doesn't include it
	if alert.Config.Metric.ID != "" {
		t.Errorf("expected Metric.ID to be empty for legacy alert, got %q", alert.Config.Metric.ID)
	}
}

func TestAlertV2MetricUnmarshalCurrentMetricID(t *testing.T) {
	// Current API response: has both "metric.id" and "metricId"
	currentJSON := `{
		"alert": {
			"id": 1760356,
			"name": "cpu usage spike alert",
			"type": "MANUAL",
			"config": {
				"metric": {
					"id": "sysdig_container_cpu_quota_used_percent"
				},
				"metricId": "sysdig_container_cpu_quota_used_percent",
				"groupAggregation": "sum",
				"timeAggregation": "avg",
				"threshold": 90.0,
				"conditionOperator": ">",
				"noDataBehaviour": "DO_NOTHING",
				"range": 60,
				"duration": 0
			}
		}
	}`

	var wrapper alertV2MetricWrapper
	err := json.Unmarshal([]byte(currentJSON), &wrapper)
	if err != nil {
		t.Fatalf("failed to unmarshal current alert JSON: %v", err)
	}

	alert := wrapper.Alert

	if alert.Config.Metric.ID != "sysdig_container_cpu_quota_used_percent" {
		t.Errorf("expected Metric.ID to be %q, got %q", "sysdig_container_cpu_quota_used_percent", alert.Config.Metric.ID)
	}

	if alert.Config.MetricID != "sysdig_container_cpu_quota_used_percent" {
		t.Errorf("expected MetricID to be %q, got %q", "sysdig_container_cpu_quota_used_percent", alert.Config.MetricID)
	}
}
