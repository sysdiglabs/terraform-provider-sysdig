package v2

import (
	proto "github.com/draios/protorepo/cloudauth/go"
)

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
	ZoneIDs             []int             `json:"zoneIds,omitempty"`
	AllZones            bool              `json:"allZones"`
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

type CustomRole struct {
	ID                 int      `json:"id,omitempty"`
	Name               string   `json:"name"`
	Description        string   `json:"description,omitempty"`
	MonitorPermissions []string `json:"monitorPermissions,omitempty"`
	SecurePermissions  []string `json:"securePermissions,omitempty"`
}
type customRoleListWrapper struct {
	Roles []CustomRole `json:"roles"`
}

type Dependency struct {
	PermissionAuthority string   `json:"permissionAuthority"`
	Dependencies        []string `json:"dependencies"`
}

type Dependencies []Dependency

type TeamServiceAccount struct {
	ID             int    `json:"id,omitempty"`
	Name           string `json:"name"`
	SystemRole     string `json:"systemRole"`
	TeamId         int    `json:"teamId"`
	TeamRole       string `json:"teamRole"`
	DateCreated    int64  `json:"dateCreated,omitempty"`
	ExpirationDate int64  `json:"expirationDate"`
	ApiKey         string `json:"apiKey,omitempty"`
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

type NotificationChannelTemplateConfigurationSection struct {
	SectionName string `json:"sectionName,omitempty"`
	ShouldShow  bool   `json:"shouldShow"`
}

type NotificationChannelTemplateConfiguration struct {
	TemplateKey                   string                                            `json:"templateKey,omitempty"`
	TemplateConfigurationSections []NotificationChannelTemplateConfigurationSection `json:"templateConfigurationSections,omitempty"`
}

type NotificationChannelOptions struct {
	EmailRecipients          []string                                   `json:"emailRecipients,omitempty"`          // Type: email
	SnsTopicARNs             []string                                   `json:"snsTopicARNs,omitempty"`             // Type: SNS
	APIKey                   string                                     `json:"apiKey,omitempty"`                   // Type: VictorOps, ibm event function
	RoutingKey               string                                     `json:"routingKey,omitempty"`               // Type: VictorOps
	Url                      string                                     `json:"url,omitempty"`                      // Type: OpsGenie, Webhook, Slack, google chat, prometheus alert manager, custom webhook, ms teams
	Channel                  string                                     `json:"channel,omitempty"`                  // Type: Slack
	Account                  string                                     `json:"account,omitempty"`                  // Type: PagerDuty
	ServiceKey               string                                     `json:"serviceKey,omitempty"`               // Type: PagerDuty
	ServiceName              string                                     `json:"serviceName,omitempty"`              // Type: PagerDuty
	AdditionalHeaders        map[string]interface{}                     `json:"additionalHeaders,omitempty"`        // Type: Webhook, prometheus alert manager, custom webhook, ibm function
	Region                   string                                     `json:"region,omitempty"`                   // Type: OpsGenie
	AllowInsecureConnections *bool                                      `json:"allowInsecureConnections,omitempty"` // Type: prometheus alert manager, custom webhook, Webhook
	TeamId                   int                                        `json:"teamId,omitempty"`                   // Type: team email
	HttpMethod               string                                     `json:"httpMethod,omitempty"`               // Type: custom webhook
	MonitorTemplate          string                                     `json:"monitorTemplate,omitempty"`          // Type: custom webhook
	InstanceId               string                                     `json:"instanceId,omitempty"`               // Type: ibm event notification
	IbmFunctionType          string                                     `json:"ibmFunctionType,omitempty"`          // Type: ibm event function
	CustomData               map[string]interface{}                     `json:"customData,omitempty"`               // Type: ibm function, Webhook
	TemplateConfiguration    []NotificationChannelTemplateConfiguration `json:"templateConfiguration,omitempty"`    // Type: slack, ms teams

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
	Weight     int      `json:"weight,omitempty"`
}

type GroupMappingConfig struct {
	NoMappingStrategy             string `json:"noMappingStrategy"`
	DifferentTeamSameRoleStrategy string `json:"differentRolesSameTeamStrategy"`
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

type Policy struct {
	ID                     int           `json:"id,omitempty"`
	IsDefault              bool          `json:"isDefault"`
	Name                   string        `json:"name"`
	Description            string        `json:"description"`
	Severity               int           `json:"severity"`
	Enabled                bool          `json:"enabled"`
	RuleNames              []string      `json:"ruleNames"`
	Rules                  []*PolicyRule `json:"rules,omitempty"`
	Actions                []Action      `json:"actions"`
	Scope                  string        `json:"scope,omitempty"`
	Version                int           `json:"version,omitempty"`
	NotificationChannelIds []int         `json:"notificationChannelIds"`
	Type                   string        `json:"type"`
	Runbook                string        `json:"runbook"`
	TemplateId             int           `json:"templateId"`
	TemplateVersion        string        `json:"templateVersion"`
}

type PolicyRule struct {
	Name    string `json:"ruleName"`
	Enabled bool   `json:"enabled"`
}

type Action struct {
	AfterEventNs         int    `json:"afterEventNs,omitempty"`
	BeforeEventNs        int    `json:"beforeEventNs,omitempty"`
	Name                 string `json:"name,omitempty"`
	IsLimitedToContainer bool   `json:"isLimitedToContainer"`
	Type                 string `json:"type"`
}

type List struct {
	Name    string `json:"name"`
	Items   Items  `json:"items"`
	Append  bool   `json:"append"`
	ID      int    `json:"id,omitempty"`
	Version int    `json:"version,omitempty"`
}

type Items struct {
	Items []string `json:"items"`
}

type Macro struct {
	ID                   int            `json:"id,omitempty"`
	Version              int            `json:"version,omitempty"`
	Name                 string         `json:"name"`
	Condition            MacroCondition `json:"condition"`
	Append               bool           `json:"append"`
	MinimumEngineVersion *int           `json:"minimumEngineVersion,omitempty"`
}

type MacroCondition struct {
	Condition string `json:"condition"`
}

type VulnerabilityExceptionList struct {
	ID      string `json:"id,omitempty"`
	Version string `json:"version"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type VulnerabilityException struct {
	ID             string `json:"id"`
	Gate           string `json:"gate"`
	TriggerID      string `json:"trigger_id"`
	Notes          string `json:"notes"`
	ExpirationDate *int   `json:"expiration_date,omitempty"`
	Enabled        bool   `json:"enabled"`
}

type Rule struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags"`
	Details     Details  `json:"details"`
	Version     int      `json:"version,omitempty"`
}

const (
	RuleTypeContainer  = "CONTAINER"
	RuleTypeFalco      = "FALCO"
	RuleTypeFilesystem = "FILESYSTEM"
	RuleTypeNetwork    = "NETWORK"
	RuleTypeProcess    = "PROCESS"
	RuleTypeSyscall    = "SYSCALL"
)

type Details struct {
	// Containers
	Containers *Containers `json:"containers,omitempty"`

	// Filesystems
	ReadWritePaths *ReadWritePaths `json:"readWritePaths,omitempty"`
	ReadPaths      *ReadPaths      `json:"readPaths,omitempty"`

	// Network
	AllOutbound    bool            `json:"allOutbound"`
	AllInbound     bool            `json:"allInbound"`
	TCPListenPorts *TCPListenPorts `json:"tcpListenPorts,omitempty"`
	UDPListenPorts *UDPListenPorts `json:"udpListenPorts,omitempty"`

	// Processes
	Processes *Processes `json:"processes,omitempty"`

	// Syscalls
	Syscalls *Syscalls `json:"syscalls,omitempty"`

	// Falco
	Append               *bool        `json:"append,omitempty"`
	Source               string       `json:"source,omitempty"`
	Output               string       `json:"output,omitempty"`
	Condition            *Condition   `json:"condition,omitempty"`
	Priority             string       `json:"priority,omitempty"`
	Exceptions           []*Exception `json:"exceptions,omitempty"`
	MinimumEngineVersion *int         `json:"minimumEngineVersion,omitempty"`

	RuleType string `json:"ruleType"`
}

type Containers struct {
	Items      []string `json:"items"`
	MatchItems bool     `json:"matchItems"`
}

type ReadWritePaths struct {
	Items      []string `json:"items"`
	MatchItems bool     `json:"matchItems"`
}
type ReadPaths struct {
	Items      []string `json:"items"`
	MatchItems bool     `json:"matchItems"`
}

type TCPListenPorts struct {
	Items      []string `json:"items"`
	MatchItems bool     `json:"matchItems"`
}

type UDPListenPorts struct {
	Items      []string `json:"items"`
	MatchItems bool     `json:"matchItems"`
}

type Processes struct {
	Items      []string `json:"items"`
	MatchItems bool     `json:"matchItems"`
}

type Syscalls struct {
	Items      []string `json:"items"`
	MatchItems bool     `json:"matchItems"`
}

type Condition struct {
	Condition  string        `json:"condition"`
	Components []interface{} `json:"components"`
}

type Exception struct {
	Name   string      `json:"name"`
	Fields interface{} `json:"fields,omitempty"`
	Comps  interface{} `json:"comps,omitempty"`
	Values interface{} `json:"values,omitempty"`
}

type CloudAccountSecure struct {
	AccountID                    string `json:"accountId"`
	Provider                     string `json:"provider"`
	Alias                        string `json:"alias"`
	RoleAvailable                bool   `json:"roleAvailable"`
	RoleName                     string `json:"roleName"`
	ExternalID                   string `json:"externalId,omitempty"`
	WorkLoadIdentityAccountID    string `json:"workloadIdentityAccountId,omitempty"`
	WorkLoadIdentityAccountAlias string `json:"workLoadIdentityAccountAlias,omitempty"`
}
type ScanningPolicy struct {
	ID             string         `json:"id,omitempty"`
	Version        string         `json:"version,omitempty"`
	Name           string         `json:"name"`
	Comment        string         `json:"comment"`
	IsDefault      bool           `json:"isDefault,omitempty"`
	PolicyBundleId string         `json:"policyBundleId,omitempty"`
	Rules          []ScanningGate `json:"rules"`
}

type ScanningGate struct {
	ID      string              `json:"id,omitempty"`
	Gate    string              `json:"gate"`
	Trigger string              `json:"trigger"`
	Action  string              `json:"action"`
	Params  []ScanningGateParam `json:"params"`
}

type ScanningGateParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ScanningPolicyAssignmentList struct {
	Items          []ScanningPolicyAssignment `json:"items"`
	PolicyBundleId string                     `json:"policyBundleId"`
}

type ScanningPolicyAssignment struct {
	ID           string                        `json:"id,omitempty"`
	Name         string                        `json:"name"`
	Registry     string                        `json:"registry"`
	Repository   string                        `json:"repository"`
	Image        ScanningPolicyAssignmentImage `json:"image"`
	PolicyIDs    []string                      `json:"policy_ids"`
	WhitelistIDs []string                      `json:"whitelist_ids"`
}

type ScanningPolicyAssignmentImage struct {
	Type  string `json:"type"`
	Value string `json:"value"`
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

type AlertV2Common struct {
	ID                            int                           `json:"id,omitempty"`
	Version                       int                           `json:"version,omitempty"`
	Name                          string                        `json:"name"`
	Description                   string                        `json:"description,omitempty"`
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

type AlertV2ConfigPrometheus struct {
	Query            string `json:"query"`
	KeepFiringForSec *int   `json:"keepFiringForSec,omitempty"`
}

type AlertV2Prometheus struct {
	AlertV2Common
	DurationSec int                     `json:"durationSec"`
	Config      AlertV2ConfigPrometheus `json:"config"`
}

type alertV2PrometheusWrapper struct {
	Alert AlertV2Prometheus `json:"alert"`
}

type AlertLabelDescriptorV2 struct {
	ID       string `json:"id"`
	PublicID string `json:"publicId,omitempty"`
}

type ScopeExpressionV2 struct {
	Operand    string                  `json:"operand"`
	Descriptor *AlertLabelDescriptorV2 `json:"descriptor,omitempty"`
	Operator   string                  `json:"operator"`
	Value      []string                `json:"value"`
}

type AlertScopeV2 struct {
	Expressions []ScopeExpressionV2 `json:"expressions,omitempty"`
}

type ScopedSegmentedConfig struct {
	Scope     *AlertScopeV2            `json:"scope,omitempty"`
	SegmentBy []AlertLabelDescriptorV2 `json:"segmentBy"`
}

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
	DurationSec int                `json:"durationSec"`
	Config      AlertV2ConfigEvent `json:"config"`
}

type alertV2EventWrapper struct {
	Alert AlertV2Event `json:"alert"`
}

type LabelDescriptorV3 struct {
	ID       string `json:"id"`
	PublicID string `json:"publicId"`
}

type labelsDescriptorV3 struct {
	AllLabels []LabelDescriptorV3 `json:"allLabels"`
}

type labelDescriptorV3 struct {
	LabelDescriptor LabelDescriptorV3 `json:"labelDescriptor"`
}

type AlertMetricDescriptorV2 struct {
	ID string `json:"id"`
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

type AlertV2Metric struct {
	AlertV2Common
	DurationSec int                 `json:"durationSec"`
	Config      AlertV2ConfigMetric `json:"config"`
}

type alertV2MetricWrapper struct {
	Alert AlertV2Metric `json:"alert"`
}

type AlertV2ConfigDowntime struct {
	ScopedSegmentedConfig

	ConditionOperator string  `json:"conditionOperator"`
	Threshold         float64 `json:"threshold"`

	GroupAggregation string                  `json:"groupAggregation"`
	TimeAggregation  string                  `json:"timeAggregation"`
	Metric           AlertMetricDescriptorV2 `json:"metric"`
	NoDataBehaviour  string                  `json:"noDataBehaviour"`
}

type AlertV2Downtime struct {
	AlertV2Common
	DurationSec int                   `json:"durationSec"`
	Config      AlertV2ConfigDowntime `json:"config"`
}

type alertV2DowntimeWrapper struct {
	Alert AlertV2Downtime `json:"alert"`
}

type AlertV2ConfigChange struct {
	ScopedSegmentedConfig

	ConditionOperator        string   `json:"conditionOperator"`
	Threshold                float64  `json:"threshold"`
	WarningConditionOperator string   `json:"warningConditionOperator,omitempty"`
	WarningThreshold         *float64 `json:"warningThreshold,omitempty"`

	GroupAggregation string                  `json:"groupAggregation"`
	TimeAggregation  string                  `json:"timeAggregation"`
	Metric           AlertMetricDescriptorV2 `json:"metric"`

	ShorterRangeSec int `json:"shorterRangeSec"`
	LongerRangeSec  int `json:"longerRangeSec"`
}

type AlertV2ConfigFormBasedPrometheus struct {
	ScopedSegmentedConfig

	Query                    string   `json:"query"`
	ConditionOperator        string   `json:"conditionOperator"`
	Threshold                float64  `json:"threshold"`
	WarningConditionOperator string   `json:"warningConditionOperator,omitempty"`
	WarningThreshold         *float64 `json:"warningThreshold,omitempty"`
	NoDataBehaviour          string   `json:"noDataBehaviour"`
}

type AlertV2FormBasedPrometheus struct {
	AlertV2Common
	DurationSec int                              `json:"durationSec"` // not really used but the api wants it set to 0 in POST/PUT
	Config      AlertV2ConfigFormBasedPrometheus `json:"config"`
}

type alertV2FormBasedPrometheusWrapper struct {
	Alert AlertV2FormBasedPrometheus `json:"alert"`
}

type AlertV2Change struct {
	AlertV2Common
	DurationSec int                 `json:"durationSec"` // not really used but the api wants it set to 0 in POST/PUT
	Config      AlertV2ConfigChange `json:"config"`
}

type alertV2ChangeWrapper struct {
	Alert AlertV2Change `json:"alert"`
}

type CloudAccountCredentialsMonitor struct {
	AccountId string `json:"accountId"`
}

type CloudAccountMonitor struct {
	Id                int                            `json:"id"`
	Platform          string                         `json:"platform"`
	IntegrationType   string                         `json:"integrationType"`
	Credentials       CloudAccountCredentialsMonitor `json:"credentials"`
	AdditionalOptions string                         `json:"additionalOptions"`
}

type cloudAccountWrapperMonitor struct {
	CloudAccount CloudAccountMonitor `json:"provider"`
}

type PosturePolicyZoneMeta struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PosturePolicy struct {
	ID             string                  `json:"id,omitempty"`
	Name           string                  `json:"name"`
	Type           int                     `json:"type"`
	Kind           int                     `json:"kind"`
	Description    string                  `json:"description"`
	Version        string                  `json:"version"`
	Link           string                  `json:"link"`
	Authors        string                  `json:"authors"`
	PublishedData  string                  `json:"publishedDate"`
	MinKubeVersion float64                 `json:"minKubeVersion"`
	MaxKubeVersion float64                 `json:"maxKubeVersion"`
	IsCustom       bool                    `json:"isCustom"`
	IsActive       bool                    `json:"isActive"`
	Platform       string                  `json:"platform"`
	Zones          []PosturePolicyZoneMeta `json:"zones"`
}

type PostureZonePolicyListResponse struct {
	Data []PosturePolicy `json:"data"`
}

type PostureZoneScope struct {
	ID         string `json:"id,omitempty"`
	TargetType string `json:"targetType"`
	Rules      string `json:"rules"`
}

type PostureZonePolicySlim struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Type int    `json:"type"`
	Kind int    `json:"kind"`
}

type PostureZone struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	Description    string                  `json:"description"`
	Author         string                  `json:"author"`
	LastModifiedBy string                  `json:"lastModifiedBy"`
	LastUpdated    string                  `json:"lastUpdated"`
	IsSystem       bool                    `json:"isSystem"`
	Scopes         []PostureZoneScope      `json:"scopes"`
	Policies       []PostureZonePolicySlim `json:"policies"`
}

type PostureZoneRequest struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	PolicyIDs   []string           `json:"policyIds"`
	Scopes      []PostureZoneScope `json:"scopes"`
}

type PostureZoneResponse struct {
	Data PostureZone `json:"data"`
}

type IdentityContext struct {
	IdentityType       string `json:"identityType"`
	CustomerID         int    `json:"customerId"`
	TeamID             int    `json:"teamId"`
	TeamName           string `json:"teamName"`
	UserID             int    `json:"userId"`
	Username           string `json:"username"`
	ServiceAccountID   int    `json:"serviceAccountId"`
	ServiceAccountName string `json:"serviceAccountName"`
}

type SilenceRule struct {
	Name                   string `json:"name"`
	Enabled                bool   `json:"enabled"`
	StartTs                int64  `json:"startTs"`
	DurationInSec          int    `json:"durationInSec"`
	Scope                  string `json:"scope,omitempty"`
	AlertIds               []int  `json:"alertIds,omitempty"`
	NotificationChannelIds []int  `json:"notificationChannelIds,omitempty"`

	Version int `json:"version,omitempty"`
	ID      int `json:"id,omitempty"`
}

type OrganizationSecure proto.CloudOrganization
