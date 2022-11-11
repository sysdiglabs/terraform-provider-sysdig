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
	Links                         *[]AlertLinkV2                 `json:"links,omitempty"`
}

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

// AlertScopeV2
type AlertScopeV2 struct {
	Expressions []ScopeExpressionV2 `json:"expressions"`
}

type AlertLabelDescriptorV2 struct {
	ID       string `json:"id"`
	PublicID string `json:"publicId"`
}

type NotificationGroupingConditionV2 struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
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
	Operand    string                  `json:"operand"`
	Descriptor AlertLabelDescriptorV2  `json:"descriptor"`
	Operator   ScopeExpressionOperator `json:"operator"`
	Value      []string                `json:"value"`
}

type ScopeExpressionOperator struct {
}

type NotificationChannelConfigV2 struct {
	// Type can be one of EMAIL, SNS, SLACK, PAGER_DUTY, VICTOROPS, OPSGENIE, WEBHOOK, IBM_FUNCTION, MS_TEAMS, TEAM_EMAIL, IBM_EVENT_NOTIFICATIONS, PROMETHEUS_ALERT_MANAGER
	ChannelID int                           `json:"channelId,omitempty"`
	Type      string                        `json:"type,omitempty"`
	Name      string                        `json:"nam,omitempty"`
	Enabled   bool                          `json:"enabled,omitempty"`
	Options   *NotificationChannelOptionsV2 `json:"options,omitempty"`
}

type NotificationChannelOptionsV2 struct {
	// commons
	NotifyOnAcknowledge        bool                          `json:"notifyOnAcknowledge,omitempty"`
	NotifyOnResolve            bool                          `json:"notifyOnResolve,omitempty"`
	ReNotifyEverySec           int                           `json:"reNotifyEverySec,omitempty"`
	CustomNotificationTemplate *CustomNotificationTemplateV2 `json:"customNotificationTemplate,omitempty"`
	Thresholds                 []string                      `json:"thresholds,omitempty"`
}

type CustomNotificationTemplateV2 struct {
	Subject     string `json:"subject,omitempty"`
	PrependText string `json:"prependText,omitempty"`
	AppendText  string `json:"appendText,omitempty"`
}

type CaptureConfigV2 struct {
	DurationSec int    `json:"durationSec"`
	Storage     string `json:"storage"`
	Filter      string `json:"filter,omitempty"`
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
	Name string `json:"name"`
	Type string `json:"type"`
	ID   string `json:"id"`
	Href string `json:"href"`
}
