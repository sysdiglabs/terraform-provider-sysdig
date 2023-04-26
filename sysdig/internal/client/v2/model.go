package v2

type Team struct {
	UserRoles           []UserRoles       `json:"userRoles,omitempty"`
	Description         string            `json:"description"`
	Name                string            `json:"name"`
	ID                  int               `json:"id,omitempty"`
	Version             int               `json:"version,omitempty"`
	Origin              string            `json:"origin,omitempty"`
	LastUpdated         int64             `json:"lastUpdated,omitempty"`
	EntryPoint          *EntryPoint       `json:"entryPoint,omitempty"`
	Theme               string            `json:"theme"`
	CustomerID          int               `json:"customerId,omitempty"`
	DateCreated         int64             `json:"dateCreated,omitempty"`
	Products            []string          `json:"products,omitempty"`
	Show                string            `json:"show,omitempty"`
	Immutable           bool              `json:"immutable,omitempty"`
	CanUseSysdigCapture *bool             `json:"canUseSysdigCapture,omitempty"`
	CanUseCustomEvents  *bool             `json:"canUseCustomEvents,omitempty"`
	CanUseAwsMetrics    *bool             `json:"canUseAwsMetrics,omitempty"`
	CanUseBeaconMetrics *bool             `json:"canUseBeaconMetrics,omitempty"`
	UserCount           int               `json:"userCount,omitempty"`
	Filter              string            `json:"filter,omitempty"`
	NamespaceFilters    *NamespaceFilters `json:"namespaceFilters,omitempty"`
	DefaultTeam         bool              `json:"default,omitempty"`
}

type NamespaceFilters struct {
	IBMPlatformMetrics *string `json:"ibmPlatformMetrics"`
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

type teamWrapper struct {
	Team Team `json:"team"`
}

type User struct {
	ID          int    `json:"id,omitempty"`
	Version     int    `json:"version,omitempty"`
	SystemRole  string `json:"systemRole,omitempty"`
	Email       string `json:"username"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	CurrentTeam *int   `json:"currentTeam"`
}

type userWrapper struct {
	User User `json:"user"`
}

type usersWrapper struct {
	Users []User `json:"users"`
}

type NotificationChannelOptions struct {
	EmailRecipients   []string               `json:"emailRecipients,omitempty"`   // Type: email
	SnsTopicARNs      []string               `json:"snsTopicARNs,omitempty"`      // Type: SNS
	APIKey            string                 `json:"apiKey,omitempty"`            // Type: VictorOps
	RoutingKey        string                 `json:"routingKey,omitempty"`        // Type: VictorOps
	Url               string                 `json:"url,omitempty"`               // Type: OpsGenie, Webhook and Slack
	Channel           string                 `json:"channel,omitempty"`           // Type: Slack
	Account           string                 `json:"account,omitempty"`           // Type: PagerDuty
	ServiceKey        string                 `json:"serviceKey,omitempty"`        // Type: PagerDuty
	ServiceName       string                 `json:"serviceName,omitempty"`       // Type: PagerDuty
	AdditionalHeaders map[string]interface{} `json:"additionalHeaders,omitempty"` // Type: Webhook
	Region            string                 `json:"region,omitempty"`            // Type: OpsGenie

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
	TeamID  *int                       `json:"teamId,omitempty"`
	Options NotificationChannelOptions `json:"options"`
}

type notificationChannelListWrapper struct {
	NotificationChannels []NotificationChannel `json:"notificationChannels"`
}

type notificationChannelWrapper struct {
	NotificationChannel NotificationChannel `json:"notificationChannel"`
}

type TeamMap struct {
	AllTeams bool  `json:"allTeams"`
	TeamIDs  []int `json:"teamIds"`
}

type GroupMapping struct {
	ID         int      `json:"id,omitempty"`
	GroupName  string   `json:"groupName,omitempty"`
	Role       string   `json:"role,omitempty"`
	SystemRole string   `json:"systemRole,omitempty"`
	TeamMap    *TeamMap `json:"teamMap,omitempty"`
}

type alertWrapper struct {
	Alert Alert `json:"alert"`
}

type Alert struct {
	ID                     int                 `json:"id,omitempty"`
	Version                int                 `json:"version,omitempty"`
	Type                   string              `json:"type"`
	Name                   string              `json:"name"`
	Description            string              `json:"description"`
	Enabled                bool                `json:"enabled"`
	GroupName              string              `json:"groupName,omitempty"`
	NotificationChannelIds []int               `json:"notificationChannelIds"`
	Filter                 string              `json:"filter"`
	Severity               int                 `json:"severity"`
	Timespan               int                 `json:"timespan"`
	CustomNotification     *CustomNotification `json:"customNotification"`
	TeamID                 int                 `json:"teamId,omitempty"`
	AutoCreated            bool                `json:"autoCreated"`
	SysdigCapture          *SysdigCapture      `json:"sysdigCapture"`
	RateOfChange           bool                `json:"rateOfChange,omitempty"`
	ReNotifyMinutes        int                 `json:"reNotifyMinutes"`
	ReNotify               bool                `json:"reNotify"`
	Valid                  bool                `json:"valid"`
	SeverityLabel          string              `json:"severityLabel,omitempty"`
	SegmentBy              []string            `json:"segmentBy"`
	SegmentCondition       *SegmentCondition   `json:"segmentCondition"`
	Criteria               *Criteria           `json:"criteria,omitempty"`
	Monitor                []*Monitor          `json:"monitor,omitempty"`
	Condition              string              `json:"condition"`
	SeverityLevel          int                 `json:"severityLevel,omitempty"`
}

type CustomNotification struct {
	TitleTemplate  string `json:"titleTemplate"`
	UseNewTemplate bool   `json:"useNewTemplate"`
	PrependText    string `json:"prependText,omitempty"`
	AppendText     string `json:"appendText,omitempty"`
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

type NotificationChannelConfigV2 struct {
	ChannelID       int                          `json:"channelId,omitempty"`
	Type            string                       `json:"type,omitempty"`
	Name            string                       `json:"nam,omitempty"`
	Enabled         bool                         `json:"enabled,omitempty"`
	OverrideOptions NotificationChannelOptionsV2 `json:"overrideOptions"`
}

type NotificationChannelOptionsV2 struct {
	NotifyOnAcknowledge        bool                          `json:"notifyOnAcknowledge,omitempty"`
	NotifyOnResolve            bool                          `json:"notifyOnResolve"`
	ReNotifyEverySec           *int                          `json:"reNotifyEverySec"`
	CustomNotificationTemplate *CustomNotificationTemplateV2 `json:"customNotificationTemplate,omitempty"`
	Thresholds                 []string                      `json:"thresholds"`
}

type CustomNotificationTemplateV2 struct {
	Subject     string `json:"subject"`
	PrependText string `json:"prependText"`
	AppendText  string `json:"appendText"`
}

type ScopedSegmentedConfig struct {
	Scope     *AlertScopeV2            `json:"scope,omitempty"`
	SegmentBy []AlertLabelDescriptorV2 `json:"segmentBy"`
}

type AlertScopeV2 struct {
	Expressions []ScopeExpressionV2 `json:"expressions,omitempty"`
}

type ScopeExpressionV2 struct {
	Operand    string                  `json:"operand"`
	Descriptor *AlertLabelDescriptorV2 `json:"descriptor,omitempty"`
	Operator   string                  `json:"operator"`
	Value      []string                `json:"value"`
}

type AlertLabelDescriptorV2 struct {
	ID       string `json:"id"`
	PublicID string `json:"publicId,omitempty"`
}

type LabelDescriptorV3 struct {
	ID       string `json:"id"`
	PublicID string `json:"publicId"`
}

var labelsDescriptorsV3Wrapper struct {
	LabelDescriptorV3 `json:"labelDescriptor"`
}

var labelsWrapper struct {
	AllLabels []LabelDescriptorV3 `json:"allLabels"`
}

type AlertV2Prometheus struct {
	AlertV2Common
	Config AlertV2ConfigPrometheus `json:"config"`
}

type AlertV2PrometheusWrapper struct {
	Alert AlertV2Prometheus `json:"alert"`
}

type AlertV2Common struct {
	ID                            int                           `json:"id,omitempty"`
	Version                       int                           `json:"version,omitempty"`
	Name                          string                        `json:"name"`
	Description                   string                        `json:"description,omitempty"`
	DurationSec                   int                           `json:"durationSec"`
	Type                          string                        `json:"type"`
	Group                         string                        `json:"group,omitempty"`
	Severity                      string                        `json:"severity"`
	TeamID                        int                           `json:"teamId,omitempty"`
	Enabled                       bool                          `json:"enabled"`
	NotificationChannelConfigList []NotificationChannelConfigV2 `json:"notificationChannelConfigList"`
	CustomNotificationTemplate    *CustomNotificationTemplateV2 `json:"customNotificationTemplate,omitempty"`
	CaptureConfig                 *CaptureConfigV2              `json:"captureConfig,omitempty"`
	Links                         []AlertLinkV2                 `json:"links"`
}

type CaptureConfigV2 struct {
	DurationSec int    `json:"durationSec"`
	Storage     string `json:"storage"`
	Filter      string `json:"filter"`
	FileName    string `json:"fileName"`
	Enabled     bool   `json:"enabled"`
}

type AlertLinkV2 struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	Href string `json:"href,omitempty"`
}

type AlertV2ConfigPrometheus struct {
	Query string `json:"query"`
}

type AlertV2Metric struct {
	AlertV2Common
	Config AlertV2ConfigMetric `json:"config"`
}

type AlertV2ConfigMetric struct {
	ScopedSegmentedConfig

	ConditionOperator        string   `json:"conditionOperator"`
	Threshold                float64  `json:"threshold"`
	WarningConditionOperator string   `json:"warningConditionOperator,omitempty"`
	WarningThreshold         *float64 `json:"warningThreshold,omitempty"`

	GroupAggregation string                  `json:"groupAggregation"`
	TimeAggregation  string                  `json:"timeAggregation"`
	Metric           AlertMetricDescriptorV2 `json:"metric"`
	NoDataBehaviour  string                  `json:"noDataBehaviour"`
}

type AlertMetricDescriptorV2 struct {
	ID string `json:"id"`
}

type alertV2MetricWrapper struct {
	Alert AlertV2Metric `json:"alert"`
}
