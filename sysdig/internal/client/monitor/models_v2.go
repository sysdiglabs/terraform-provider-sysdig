package monitor

import (
	"bytes"
	"encoding/json"
	"io"
)

// --- Constants --
const (
	// alert types enum
	AlertV2AlertType_AdvancedManual   = "ADVANCED_MANUAL"
	AlertV2AlertType_AnomalyDetection = "ANOMALY_DETECTION"
	AlertV2AlertType_Dowtime          = "DOWNTIME"
	AlertV2AlertType_Event            = "EVENT"
	AlertV2AlertType_GroupOutlier     = "GROUP_OUTLIER"
	AlertV2AlertType_Manual           = "MANUAL"
	AlertV2AlertType_Prometheus       = "PROMETHEUS"

	// severities enum
	AlertV2Severity_High   = "high"
	AlertV2Severity_Medium = "medium"
	AlertV2Severity_Low    = "low"
	AlertV2Severity_Info   = "info"

	// alert link type
	AlertLinkV2Type_Dashboard = "dashboard"
	AlertLinkV2Type_Runbook   = "runbook"

	// others
	AlertV2CaptureFilenameRegexp = `.*?\.scap`
)

// enums severity values
func AlertV2Severity_Values() []string {
	return []string{
		AlertV2Severity_High,
		AlertV2Severity_Medium,
		AlertV2Severity_Low,
		AlertV2Severity_Info,
	}
}

// AlertV2
type AlertV2Common struct {
	ID                            int                            `json:"id,omitempty"`
	Version                       int                            `json:"version,omitempty"`
	Name                          string                         `json:"name"`
	Description                   string                         `json:"description,omitempty"`
	DurationSec                   int                            `json:"durationSec"`
	Type                          string                         `json:"type"`
	Group                         string                         `json:"group,omitempty"`
	Severity                      string                         `json:"severity"`
	TeamID                        int                            `json:"teamId,omitempty"`
	Enabled                       bool                           `json:"enabled"`
	NotificationChannelConfigList *[]NotificationChannelConfigV2 `json:"notificationChannelConfigList,omitempty"`
	CustomNotificationTemplate    *CustomNotificationTemplateV2  `json:"customNotificationTemplate,omitempty"`
	CaptureConfig                 *CaptureConfigV2               `json:"captureConfig,omitempty"`
	Links                         []AlertLinkV2                  `json:"links"`
}

type ScopedSegmentedConfig struct {
	Scope     *AlertScopeV2            `json:"scope,omitempty"`
	SegmentBy []AlertLabelDescriptorV2 `json:"segmentBy"`
}

// Prometheus
type AlertV2ConfigPrometheus struct {
	Query string `json:"query"`
}

type AlertV2Prometheus struct {
	AlertV2Common
	Config *AlertV2ConfigPrometheus `json:"config"`
}

func (a *AlertV2Prometheus) ToJSON() io.Reader {
	data := struct {
		Alert AlertV2Prometheus `json:"alert"`
	}{Alert: *a}
	payload, _ := json.Marshal(data)
	return bytes.NewBuffer(payload)
}

func AlertV2PrometheusFromJSON(body []byte) AlertV2Prometheus {
	var result struct {
		Alert AlertV2Prometheus
	}
	_ = json.Unmarshal(body, &result)
	return result.Alert
}

// Event
type AlertV2ConfigEvent struct {
	ScopedSegmentedConfig

	ConditionOperator        string   `json:"conditionOperator"`
	Threshold                float64  `json:"threshold"`
	WarningConditionOperator string   `json:"warningConditionOperator,omitempty"`
	WarningThreshold         *float64 `json:"warningThreshold,omitempty"`

	Filter string   `json:"filter"`
	Tags   []string `json:"tags"`
}

type AlertV2Event struct {
	AlertV2Common
	Config *AlertV2ConfigEvent `json:"config"`
}

func (a *AlertV2Event) ToJSON() io.Reader {
	data := struct {
		Alert AlertV2Event `json:"alert"`
	}{Alert: *a}
	payload, _ := json.Marshal(data)
	return bytes.NewBuffer(payload)
}

func AlertV2EventFromJSON(body []byte) AlertV2Event {
	var result struct {
		Alert AlertV2Event
	}
	_ = json.Unmarshal(body, &result)
	return result.Alert
}

// Metric
type AlertV2ConfigMetric struct {
	ScopedSegmentedConfig

	ConditionOperator        string   `json:"conditionOperator"`
	Threshold                float64  `json:"threshold"`
	WarningConditionOperator string   `json:"warningConditionOperator,omitempty"`
	WarningThreshold         *float64 `json:"warningThreshold,omitempty"`

	GroupAggregation string                 `json:"groupAggregation"`
	TimeAggregation  string                 `json:"timeAggregation"`
	Metric           AlertLabelDescriptorV2 `json:"metric"`
	NoDataBehaviour  string                 `json:"noDataBehaviour"`
}

type AlertV2Metric struct {
	AlertV2Common
	Config *AlertV2ConfigMetric `json:"config"`
}

func (a *AlertV2Metric) ToJSON() io.Reader {
	data := struct {
		Alert AlertV2Metric `json:"alert"`
	}{Alert: *a}
	payload, _ := json.Marshal(data)
	return bytes.NewBuffer(payload)
}

func AlertV2MetricFromJSON(body []byte) AlertV2Metric {
	var result struct {
		Alert AlertV2Metric
	}
	_ = json.Unmarshal(body, &result)
	return result.Alert
}

// Downtime
type AlertV2ConfigDowntime struct {
	ScopedSegmentedConfig

	ConditionOperator string  `json:"conditionOperator"`
	Threshold         float64 `json:"threshold"`

	GroupAggregation string                 `json:"groupAggregation"`
	TimeAggregation  string                 `json:"timeAggregation"`
	Metric           AlertLabelDescriptorV2 `json:"metric"`
	NoDataBehaviour  string                 `json:"noDataBehaviour"`
}

type AlertV2Downtime struct {
	AlertV2Common
	Config *AlertV2ConfigDowntime `json:"config"`
}

func (a *AlertV2Downtime) ToJSON() io.Reader {
	data := struct {
		Alert AlertV2Downtime `json:"alert"`
	}{Alert: *a}
	payload, _ := json.Marshal(data)
	return bytes.NewBuffer(payload)
}

func AlertV2DowntimeFromJSON(body []byte) AlertV2Downtime {
	var result struct {
		Alert AlertV2Downtime
	}
	_ = json.Unmarshal(body, &result)
	return result.Alert
}

// AlertScopeV2
type AlertScopeV2 struct {
	Expressions []ScopeExpressionV2 `json:"expressions,omitempty"`
}

type AlertLabelDescriptorV2 struct {
	ID       string `json:"id"`
	PublicID string `json:"publicId"`
}

type AlertMetricDescriptorV2 struct {
	ID                string        `json:"id"`
	PublicID          string        `json:"publicId"`
	MetricType        string        `json:"metricType"`
	Type              string        `json:"type"`
	Scale             float64       `json:"scale"`
	GroupAggregations []Aggregation `json:"groupAggregations"`
	TimeAggregations  []Aggregation `json:"timeAggregations"`
}

type Aggregation struct {
	ID               int         `json:"id"`
	Percentile       bool        `json:"percentile"`
	AggregationValue interface{} `json:"aggregationValue"`
}

type ScopeExpressionV2 struct {
	Operand    string                 `json:"operand"`    // old dot notation, e.g. "kubernetes.cluster.name"
	Descriptor AlertLabelDescriptorV2 `json:"descriptor"` // discarded by the backend in put/post
	Operator   string                 `json:"operator"`
	Value      []string               `json:"value"`
}

type NotificationChannelConfigV2 struct {
	// Type can be one of EMAIL, SNS, SLACK, PAGER_DUTY, VICTOROPS, OPSGENIE, WEBHOOK, IBM_FUNCTION, MS_TEAMS, TEAM_EMAIL, IBM_EVENT_NOTIFICATIONS, PROMETHEUS_ALERT_MANAGER
	ChannelID       int                          `json:"channelId,omitempty"`
	Type            string                       `json:"type,omitempty"`
	Name            string                       `json:"nam,omitempty"`
	Enabled         bool                         `json:"enabled,omitempty"`
	OverrideOptions NotificationChannelOptionsV2 `json:"overrideOptions"`
}

type NotificationChannelOptionsV2 struct {
	// commons
	NotifyOnAcknowledge        bool                          `json:"notifyOnAcknowledge,omitempty"`
	NotifyOnResolve            bool                          `json:"notifyOnResolve,omitempty"`
	ReNotifyEverySec           *int                          `json:"reNotifyEverySec"` // must send null to remove this opt
	CustomNotificationTemplate *CustomNotificationTemplateV2 `json:"customNotificationTemplate,omitempty"`
	Thresholds                 []string                      `json:"thresholds"` //Set of thresholds the notification channel will be used for. Possible values [MAIN, WARNING]
}

type CustomNotificationTemplateV2 struct {
	Subject     string `json:"subject"`
	PrependText string `json:"prependText"`
	AppendText  string `json:"appendText"`
}

type CaptureConfigV2 struct {
	DurationSec int    `json:"durationSec"`
	Storage     string `json:"storage"`
	Filter      string `json:"filter"`
	FileName    string `json:"fileName"`
	Enabled     bool   `json:"enabled"`
}

// enums link types values
func AlertLinkV2Type_Values() []string {
	return []string{
		AlertLinkV2Type_Dashboard,
		AlertLinkV2Type_Runbook,
	}
}

type AlertLinkV2 struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	Href string `json:"href,omitempty"`
}

type LabelDescriptorV3 struct {
	ID         string `json:"id"`
	PublicID   string `json:"publicId"`
	CanGroupBy bool   `json:"canGroupBy,omitempty"`
	Documented bool   `json:"documented,omitempty"`
}
