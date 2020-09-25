package monitor

import (
	"bytes"
	"encoding/json"
	"io"
)

// -------- Alert --------
type CustomNotification struct {
	TitleTemplate  string `json:"titleTemplate"`
	UseNewTemplate bool   `json:"useNewTemplate"`
}

type SysdigCapture struct {
	Name       string      `json:"name"`
	Filters    string      `json:"filters,omitempty"`
	Duration   int         `json:"duration"`
	Type       string      `json:"type,omitempty"`
	BucketName string      `json:"bucketName"`
	Folder     string      `json:"folder,omitempty"`
	Enabled    bool        `json:"enabled"`
	StorageID  interface{} `json:"storageId,omitempty"`
}
type SegmentCondition struct {
	Type string `json:"type"`
}

type Criteria struct {
	Text   string `json:"text"`
	Source string `json:"source"`
}

type Monitor struct {
	Metric       string  `json:"metric"`
	StdDevFactor float64 `json:"stdDevFactor"`
}

type alertWrapper struct {
	Alert Alert `json:"alert"`
}

type Alert struct {
	ID                     int                 `json:"id,omitempty"`
	Version                int                 `json:"version,omitempty"`
	Type                   string              `json:"type"` // computed MANUAL
	Name                   string              `json:"name"`
	Description            string              `json:"description"`
	Enabled                bool                `json:"enabled"`
	NotificationChannelIds []int               `json:"notificationChannelIds"`
	Filter                 string              `json:"filter"`
	Severity               int                 `json:"severity"` // 6 == INFO, 4 == LOW, 2 == MEDIUM, 0 == HIGH // NOT USED
	Timespan               int                 `json:"timespan"` // computed 600000000
	CustomNotification     *CustomNotification `json:"customNotification"`
	TeamID                 int                 `json:"teamId,omitempty"` // computed
	AutoCreated            bool                `json:"autoCreated"`
	SysdigCapture          *SysdigCapture      `json:"sysdigCapture"`
	RateOfChange           bool                `json:"rateOfChange,omitempty"`
	ReNotifyMinutes        int                 `json:"reNotifyMinutes"`
	ReNotify               bool                `json:"reNotify"`
	Valid                  bool                `json:"valid"`
	SeverityLabel          string              `json:"severityLabel,omitempty"` // MEDIUM == MEDIUM, LOW == LOW, NONE == INFO, HIGH == HIGH
	SegmentBy              []string            `json:"segmentBy"`
	SegmentCondition       *SegmentCondition   `json:"segmentCondition"`
	Criteria               *Criteria           `json:"criteria,omitempty"`
	Monitor                []*Monitor          `json:"monitor,omitempty"`
	Condition              string              `json:"condition"`
	SeverityLevel          int                 `json:"severityLevel,omitempty"` // 0 == MEDIUM, 2 == LOW, 4 == INFO, 6 == HIGH
}

func (a *Alert) ToJSON() io.Reader {
	payload, _ := json.Marshal(alertWrapper{Alert: *a})
	return bytes.NewBuffer(payload)
}

func AlertFromJSON(body []byte) Alert {
	var result alertWrapper
	json.Unmarshal(body, &result)

	return result.Alert
}

// -------- Team --------
type Team struct {
	UserRoles           []UserRoles `json:"userRoles"`
	Description         string      `json:"description"`
	Name                string      `json:"name"`
	ID                  int         `json:"id,omitempty"`
	Version             int         `json:"version,omitempty"`
	Origin              string      `json:"origin,omitempty"`
	LastUpdated         int64       `json:"lastUpdated,omitempty"`
	EntryPoint          EntryPoint  `json:"entryPoint"`
	Theme               string      `json:"theme"`
	CustomerID          int         `json:"customerId"`
	DateCreated         int64       `json:"dateCreated"`
	Products            []string    `json:"products,omitempty"`
	Show                string      `json:"show"`
	Immutable           bool        `json:"immutable"`
	CanUseSysdigCapture bool        `json:"canUseSysdigCapture"`
	CanUseCustomEvents  bool        `json:"canUseCustomEvents"`
	CanUseAwsMetrics    bool        `json:"canUseAwsMetrics"`
	CanUseBeaconMetrics bool        `json:"canUseBeaconMetrics"`
	UserCount           int         `json:"userCount"`
	Filter              string      `json:"filter,omitempty"`
	DefaultTeam         bool        `json:"default,omitempty"`
}

type UserRoles struct {
	UserId int    `json:"userId"`
	Email  string `json:"userName,omitempty"`
	Role   string `json:"role"`
	Admin  bool   `json:"admin,omitempty"`
}

type EntryPoint struct {
	Module    string `json:"module"`
	Selection string `json:"selection,omitempty"`
}

func (t *Team) ToJSON() io.Reader {
	payload, _ := json.Marshal(*t)
	return bytes.NewBuffer(payload)
}

func TeamFromJSON(body []byte) Team {
	var result teamWrapper
	json.Unmarshal(body, &result)

	return result.Team
}

type teamWrapper struct {
	Team Team `json:"team"`
}

// -------- UsersList --------
type UsersList struct {
	ID    int    `json:"id"`
	Email string `json:"username"`
}

func UsersListFromJSON(body []byte) []UsersList {
	var result usersListWrapper
	json.Unmarshal(body, &result)

	return result.UsersList
}

type usersListWrapper struct {
	UsersList []UsersList `json:"users"`
}

// -------- Notification Channels --------

type NotificationChannelOptions struct {
	EmailRecipients []string `json:"emailRecipients,omitempty"` // Type: email
	SnsTopicARNs    []string `json:"snsTopicARNs,omitempty"`    // Type: SNS
	APIKey          string   `json:"apiKey,omitempty"`          // Type: VictorOps
	RoutingKey      string   `json:"routingKey,omitempty"`      // Type: VictorOps
	Url             string   `json:"url,omitempty"`             // Type: OpsGenie, Webhook and Slack
	Channel         string   `json:"channel,omitempty"`         // Type: Slack
	Account         string   `json:"account,omitempty"`         // Type: PagerDuty
	ServiceKey      string   `json:"serviceKey,omitempty"`      // Type: PagerDuty
	ServiceName     string   `json:"serviceName,omitempty"`     // Type: PagerDuty

	NotifyOnOk           bool `json:"notifyOnOk"`
	NotifyOnResolve      bool `json:"notifyOnResolve"`
	SendTestNotification bool `json:"sendTestNotification"`
}

type NotificationChannel struct {
	ID      int                        `json:"id,omitempty"`
	Version int                        `json:"version,omitempty"`
	Type    string                     `json:"type"`
	Name    string                     `json:"name"`
	Enabled bool                       `json:"enabled"`
	Options NotificationChannelOptions `json:"options"`
}

func (n *NotificationChannel) ToJSON() io.Reader {
	payload, _ := json.Marshal(notificationChannelWrapper{*n})
	return bytes.NewBuffer(payload)
}

func NotificationChannelFromJSON(body []byte) NotificationChannel {
	var result notificationChannelWrapper
	json.Unmarshal(body, &result)

	return result.NotificationChannel
}

func NotificationChannelListFromJSON(body []byte) []NotificationChannel {
	var result notificationChannelListWrapper
	json.Unmarshal(body, &result)

	return result.NotificationChannels
}

type notificationChannelListWrapper struct {
	NotificationChannels []NotificationChannel `json:"notificationChannels"`
}

type notificationChannelWrapper struct {
	NotificationChannel NotificationChannel `json:"notificationChannel"`
}

// ---- Dashboards -----

type Layout struct {
	X       int `json:"x"`
	Y       int `json:"y"`
	W       int `json:"w"`
	H       int `json:"h"`
	PanelID int `json:"panelId"`
}
type QueryParams struct {
	Severities    []interface{} `json:"severities"`
	AlertStatuses []interface{} `json:"alertStatuses"`
	Categories    []interface{} `json:"categories"`
	Filter        string        `json:"filter"`
	TeamScope     bool          `json:"teamScope"`
}
type EventDisplaySettings struct {
	Enabled     bool        `json:"enabled"`
	QueryParams QueryParams `json:"queryParams"`
}
type Bottom struct {
	Enabled bool `json:"enabled"`
}
type Left struct {
	Enabled        bool        `json:"enabled"`
	DisplayName    interface{} `json:"displayName"`
	Unit           string      `json:"unit"`
	DisplayFormat  string      `json:"displayFormat"`
	Decimals       interface{} `json:"decimals"`
	MinValue       int         `json:"minValue"`
	MaxValue       interface{} `json:"maxValue"`
	MinInputFormat string      `json:"minInputFormat"`
	MaxInputFormat string      `json:"maxInputFormat"`
	Scale          string      `json:"scale"`
}
type Right struct {
	Enabled        bool        `json:"enabled"`
	DisplayName    interface{} `json:"displayName"`
	Unit           string      `json:"unit"`
	DisplayFormat  string      `json:"displayFormat"`
	Decimals       interface{} `json:"decimals"`
	MinValue       int         `json:"minValue"`
	MaxValue       interface{} `json:"maxValue"`
	MinInputFormat string      `json:"minInputFormat"`
	MaxInputFormat string      `json:"maxInputFormat"`
	Scale          string      `json:"scale"`
}
type AxesConfiguration struct {
	Bottom Bottom `json:"bottom"`
	Left   Left   `json:"left"`
	Right  Right  `json:"right"`
}
type LegendConfiguration struct {
	Enabled     bool        `json:"enabled"`
	Position    string      `json:"position"`
	Layout      string      `json:"layout"`
	ShowCurrent bool        `json:"showCurrent"`
	Width       interface{} `json:"width"`
	Height      interface{} `json:"height"`
}
type DisplayInfo struct {
	DisplayName                   string `json:"displayName"`
	TimeSeriesDisplayNameTemplate string `json:"timeSeriesDisplayNameTemplate"`
	Type                          string `json:"type"`
	Color                         string `json:"color,omitempty"`
	LineWidth                     int    `json:"lineWidth,omitempty"`
}
type Format struct {
	Unit          string      `json:"unit"`
	InputFormat   string      `json:"inputFormat"`
	DisplayFormat string      `json:"displayFormat"`
	Decimals      interface{} `json:"decimals"`
	YAxis         string      `json:"yAxis"`
}

func NewPercentFormat() Format {
	return Format{
		Unit:          "%",
		InputFormat:   "0-100",
		DisplayFormat: "auto",
		Decimals:      nil,
		YAxis:         "auto",
	}
}

func NewDataFormat() Format {
	return Format{
		Unit:          "byte",
		InputFormat:   "B",
		DisplayFormat: "auto",
		Decimals:      nil,
		YAxis:         "auto",
	}
}

func NewDataRateFormat() Format {
	return Format{
		Unit:          "byteRate",
		InputFormat:   "B/s",
		DisplayFormat: "auto",
		Decimals:      nil,
		YAxis:         "auto",
	}
}

func NewNumberFormat() Format {
	return Format{
		Unit:          "number",
		InputFormat:   "1",
		DisplayFormat: "auto",
		Decimals:      nil,
		YAxis:         "auto",
	}
}

func NewNumberRateFormat() Format {
	return Format{
		Unit:          "numberRate",
		InputFormat:   "/s",
		DisplayFormat: "auto",
		Decimals:      nil,
		YAxis:         "auto",
	}
}

func NewTimeFormat() Format {
	return Format{
		Unit:          "relativeTime",
		InputFormat:   "ns",
		DisplayFormat: "auto",
		Decimals:      nil,
		YAxis:         "auto",
	}
}

type AdvancedQueries struct {
	Enabled     bool        `json:"enabled"`
	DisplayInfo DisplayInfo `json:"displayInfo"`
	Format      Format      `json:"format"`
	Query       string      `json:"query"`
	ID          int         `json:"id"`
}
type Panels struct {
	ID                     int                 `json:"id"`
	Name                   string              `json:"name"`
	Description            string              `json:"description"`
	AxesConfiguration      AxesConfiguration   `json:"axesConfiguration"`
	LegendConfiguration    LegendConfiguration `json:"legendConfiguration"`
	ApplyScopeToAll        bool                `json:"applyScopeToAll"`
	ApplySegmentationToAll bool                `json:"applySegmentationToAll"`
	AdvancedQueries        []AdvancedQueries   `json:"advancedQueries"`
	NumberThresholds       NumberThresholds    `json:"numberThresholds"`
	MarkdownSource         interface{}         `json:"markdownSource"`
	PanelTitleVisible      bool                `json:"panelTitleVisible"`
	TextAutosized          bool                `json:"textAutosized"`
	TransparentBackground  bool                `json:"transparentBackground"`
	Type                   string              `json:"type"`
}

type NumberThresholds struct {
	Base   Base          `json:"base"`
	Values []interface{} `json:"values"`
}

type Base struct {
	DisplayText string `json:"displayText"`
	Severity    string `json:"severity"`
}

type TeamSharingOptions struct {
	Type          string        `json:"type"`
	UserTeamsRole string        `json:"userTeamsRole"`
	SelectedTeams []interface{} `json:"selectedTeams"`
}
type Dashboard struct {
	Version                 int                  `json:"version,omitempty"`
	CustomerID              interface{}          `json:"customerId"`
	TeamID                  int                  `json:"teamId"`
	Schema                  int                  `json:"schema"`
	AutoCreated             bool                 `json:"autoCreated"`
	PublicToken             string               `json:"publicToken"`
	ScopeExpressionList     interface{}          `json:"scopeExpressionList"`
	Layout                  []Layout             `json:"layout"`
	TeamScope               interface{}          `json:"teamScope"`
	EventDisplaySettings    EventDisplaySettings `json:"eventDisplaySettings"`
	ID                      int                  `json:"id,omitempty"`
	Name                    string               `json:"name"`
	Description             string               `json:"description"`
	Username                string               `json:"username"`
	Shared                  bool                 `json:"shared"`
	SharingSettings         []interface{}        `json:"sharingSettings"`
	Public                  bool                 `json:"public"`
	Favorite                bool                 `json:"favorite"`
	CreatedOn               int64                `json:"createdOn"`
	ModifiedOn              int64                `json:"modifiedOn"`
	Panels                  []Panels             `json:"panels"`
	TeamScopeExpressionList []interface{}        `json:"teamScopeExpressionList"`
	CreatedOnDate           string               `json:"createdOnDate"`
	ModifiedOnDate          string               `json:"modifiedOnDate"`
	TeamSharingOptions      TeamSharingOptions   `json:"teamSharingOptions"`
}

type dashboardWrapper struct {
	Dashboard Dashboard `json:"dashboard"`
}

func (db *Dashboard) ToJSON() io.Reader {
	payload, _ := json.Marshal(dashboardWrapper{*db})
	return bytes.NewBuffer(payload)
}

func DashboardFromJSON(body []byte) Dashboard {
	var result dashboardWrapper
	json.Unmarshal(body, &result)

	return result.Dashboard
}
