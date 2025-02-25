package v2

import (
	"encoding/json"
	"errors"

	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
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

type IPFilter struct {
	ID          int    `json:"id,omitempty"`
	IPRange     string `json:"ipRange"`
	Note        string `json:"note,omitempty"`
	Enabled     bool   `json:"isEnabled"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

type IPFiltersSettings struct {
	IPFilteringEnabled bool `json:"isFilteringEnabled"`
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
	Timespan               *int                `json:"timespan,omitempty"`
	Duration               *int                `json:"duration,omitempty"`
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
	Origin                 string        `json:"origin"`
	Runbook                string        `json:"runbook"`
	TemplateId             int           `json:"templateId"`
	TemplateVersion        string        `json:"templateVersion"`
}

type PolicyRulesComposite struct {
	Policy *Policy              `json:"policy"`
	Rules  []*RuntimePolicyRule `json:"rules"`
}

type (
	FlexInt                   int
	RuntimePolicyObjectOrigin string // NOTE: This is an int in model_rules.go#L247
	ElementType               string
)

type RuntimePolicyRuleDetails interface {
	GetRuleType() ElementType
}

type RuntimePolicyRule struct {
	Id          *FlexInt                   `json:"id"`
	Name        string                     `json:"name"`
	Origin      *RuntimePolicyObjectOrigin `json:"origin"`
	VersionId   string                     `json:"versionId"`
	Filename    string                     `json:"filename"`
	Description string                     `json:"description"`
	Details     RuntimePolicyRuleDetails   `json:"details"`
	Tags        []string                   `json:"tags"`
	Version     *int                       `json:"version"`
	CreatedOn   int64                      `json:"createdOn"`
	ModifiedOn  int64                      `json:"modifiedOn"`
}

// TODO: This should be exported into a common package
// Copied from: https://github.com/draios/secure-backend/blob/main/policies/model/model_rules.go#L676C1-L779C2
func (r *RuntimePolicyRule) UnmarshalJSON(b []byte) error {
	type findType struct {
		RuleType string `json:"ruleType"`
	}
	findDetails := struct {
		FindType    findType                   `json:"details"`
		Id          *FlexInt                   `json:"id"`
		Name        string                     `json:"name"`
		Origin      *RuntimePolicyObjectOrigin `json:"origin"`
		VersionId   string                     `json:"versionId"`
		Filename    string                     `json:"filename"`
		Description string                     `json:"description"`
		Tags        []string                   `json:"tags"`
		Version     *int                       `json:"version"`
		CreatedOn   int64                      `json:"createdOn"`
		ModifiedOn  int64                      `json:"modifiedOn"`
	}{}

	err := json.Unmarshal(b, &findDetails)
	if err != nil {
		return err
	}

	var d RuntimePolicyRuleDetails
	switch findDetails.FindType.RuleType {
	// case "FALCO":
	// 	d = &FalcoRuleDetails{}
	// case "CONTAINER":
	// 	d = &ContainerImageRuleDetails{}
	// case "FILESYSTEM":
	// 	d = &FilesystemRuleDetails{}
	// case "NETWORK":
	// 	d = &NetworkRuleDetails{
	// 		AllInbound:  true,
	// 		AllOutbound: true,
	// 	}
	// case "PROCESS":
	// 	d = &ProcessRuleDetails{}
	// case "SYSCALL":
	// 	d = &SyscallRuleDetails{}
	case "DRIFT":
		d = &DriftRuleDetails{}
	case "MACHINE_LEARNING":
		d = &MLRuleDetails{}
	case "AWS_MACHINE_LEARNING":
		d = &AWSMLRuleDetails{}
	case "MALWARE":
		d = &MalwareRuleDetails{}
	// case "OKTA_MACHINE_LEARNING":
	// 	d = &OktaMLRuleDetails{}
	default:
		return errors.New("The field details has an unknown ruleType: " + findDetails.FindType.RuleType)
	}

	getRawDetails := struct {
		RawDetails json.RawMessage `json:"details"`
	}{}
	err = json.Unmarshal(b, &getRawDetails)
	if err != nil {
		return err
	}

	err = json.Unmarshal(getRawDetails.RawDetails, d)
	if err != nil {
		return err
	}

	var findDetailsIdPtr *FlexInt
	if findDetails.Id != nil {
		findDetailsId := FlexInt(*findDetails.Id)
		findDetailsIdPtr = &findDetailsId
	}

	r.Id = findDetailsIdPtr
	r.Name = findDetails.Name
	r.Origin = findDetails.Origin
	r.VersionId = findDetails.VersionId
	r.Filename = findDetails.Filename
	r.Description = findDetails.Description
	r.Tags = findDetails.Tags
	r.Version = findDetails.Version
	r.Details = d
	r.CreatedOn = findDetails.CreatedOn
	r.ModifiedOn = findDetails.ModifiedOn

	return nil
}

type MLRuleThresholdAndSeverity struct {
	Enabled   bool    `json:"enabled" yaml:"enabled"`
	Threshold float64 `json:"threshold" yaml:"threshold"`
	Severity  float64 `json:"severity" yaml:"severity"`
}

type MLRuleDetails struct {
	RuleType                ElementType                 `json:"ruleType" yaml:"ruleType"`
	AnomalyDetectionTrigger *MLRuleThresholdAndSeverity `json:"anomalyDetectionTrigger" yaml:"anomalyDetectionTrigger"`
	CryptominingTrigger     *MLRuleThresholdAndSeverity `json:"cryptominingTrigger" yaml:"cryptominingTrigger"`
	Details                 `json:"-"`
}

func (p MLRuleDetails) GetRuleType() ElementType {
	return p.RuleType
}

type MalwareRuleDetails struct {
	RuleType         ElementType         `json:"ruleType"`
	UseManagedHashes bool                `json:"useManagedHashes"`
	AdditionalHashes map[string][]string `json:"additionalHashes"`
	IgnoreHashes     map[string][]string `json:"ignoreHashes"`
	Details          `json:"-"`
}

func (p MalwareRuleDetails) GetRuleType() ElementType {
	return p.RuleType
}

type RuntimePolicyRuleList struct {
	Items      []string `json:"items"`
	MatchItems bool     `json:"matchItems"`
}

type DriftRuleDetails struct {
	RuleType                  ElementType            `json:"ruleType"`
	Exceptions                *RuntimePolicyRuleList `json:"exceptionList"`
	ProcessBasedExceptions    *RuntimePolicyRuleList `json:"allowlistProcess"`
	ProcessBasedDenylist      *RuntimePolicyRuleList `json:"denylistProcess"`
	ProhibitedBinaries        *RuntimePolicyRuleList `json:"prohibitedBinaries"`
	Mode                      string                 `json:"mode"`
	MountedVolumeDriftEnabled bool                   `json:"mountedVolumeDriftEnabled"`
	Details                   `json:"-"`
}

func (p DriftRuleDetails) GetRuleType() ElementType {
	return p.RuleType
}

type AWSMLRuleDetails struct {
	RuleType              ElementType                 `json:"ruleType" yaml:"ruleType"`
	AnomalousConsoleLogin *MLRuleThresholdAndSeverity `json:"anomalousConsoleLogin" yaml:"anomalousConsoleLogin"`
	Details               `json:"-"`
}

func (p AWSMLRuleDetails) GetRuleType() ElementType {
	return p.RuleType
}

type PolicyRule struct {
	Name    string `json:"ruleName"`
	Enabled bool   `json:"enabled"`
}

// Did not add support storageId because FE does not support it yet
type Action struct {
	AfterEventNs         int     `json:"afterEventNs,omitempty"`
	BeforeEventNs        int     `json:"beforeEventNs,omitempty"`
	Name                 string  `json:"name,omitempty"`
	Filter               string  `json:"filter,omitempty"`
	StorageType          string  `json:"storageType,omitempty"`
	BucketName           string  `json:"bucketName,omitempty"`
	Folder               string  `json:"folder,omitempty"`
	IsLimitedToContainer bool    `json:"isLimitedToContainer"`
	Type                 string  `json:"type"`
	Msg                  *string `json:"msg,omitempty"`
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
	RuleTypeContainer           = "CONTAINER"
	RuleTypeFalco               = "FALCO"
	RuleTypeFilesystem          = "FILESYSTEM"
	RuleTypeNetwork             = "NETWORK"
	RuleTypeProcess             = "PROCESS"
	RuleTypeSyscall             = "SYSCALL"
	RuleTypeStatefulSequence    = "STATEFUL_SEQUENCE"
	StatefulUniqPercentRuleType = "STATEFUL_UNIQ_PERCENT"
	StatefulCountRuleType       = "STATEFUL_COUNT"
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

type CloudauthAccountSecure struct {
	cloudauth.CloudAccount
}

type CloudauthAccountComponentSecure struct {
	cloudauth.AccountComponent
}

type CloudauthAccountFeatureSecure struct {
	cloudauth.AccountFeature
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
	Labels                        map[string]interface{}        `json:"labels,omitempty"`
}

type AlertV2ConfigPrometheus struct {
	Query            string `json:"query"`
	KeepFiringForSec *int   `json:"keepFiringForSec,omitempty"`

	Duration int `json:"duration"`
}

type AlertV2Prometheus struct {
	AlertV2Common
	Config AlertV2ConfigPrometheus `json:"config"`
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

	Range int `json:"range"`
}

type AlertV2Event struct {
	AlertV2Common
	Config AlertV2ConfigEvent `json:"config"`
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

	Range    int `json:"range"`
	Duration int `json:"duration"`
}

type AlertV2Metric struct {
	AlertV2Common
	Config                                   AlertV2ConfigMetric `json:"config"`
	UnreportedAlertNotificationsRetentionSec *int                `json:"unreportedAlertNotificationsRetentionSec"`
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

	Range int `json:"range"`
}

type AlertV2Downtime struct {
	AlertV2Common
	Config                                   AlertV2ConfigDowntime `json:"config"`
	UnreportedAlertNotificationsRetentionSec *int                  `json:"unreportedAlertNotificationsRetentionSec"`
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

	Duration int `json:"duration"`
}

type AlertV2FormBasedPrometheus struct {
	AlertV2Common
	Config                                   AlertV2ConfigFormBasedPrometheus `json:"config"`
	UnreportedAlertNotificationsRetentionSec *int                             `json:"unreportedAlertNotificationsRetentionSec"`
}

type alertV2FormBasedPrometheusWrapper struct {
	Alert AlertV2FormBasedPrometheus `json:"alert"`
}

type AlertV2ConfigGroupOutlier struct {
	ScopedSegmentedConfig

	Algorithm       string  `json:"algorithm"`
	MadThreshold    float64 `json:"madThreshold,omitempty"`
	MadTolerance    float64 `json:"madTolerance,omitempty"`
	DbscanTolerance float64 `json:"dbscanTolerance,omitempty"`

	GroupAggregation string                  `json:"groupAggregation"`
	TimeAggregation  string                  `json:"timeAggregation"`
	Metric           AlertMetricDescriptorV2 `json:"metric"`
	NoDataBehaviour  string                  `json:"noDataBehaviour"`

	ObservationWindow int `json:"observationWindow"`
}

type AlertV2GroupOutlier struct {
	AlertV2Common
	Config                                   AlertV2ConfigGroupOutlier `json:"config"`
	UnreportedAlertNotificationsRetentionSec *int                      `json:"unreportedAlertNotificationsRetentionSec"`
}

type alertV2GroupOutlierWrapper struct {
	Alert AlertV2GroupOutlier `json:"alert"`
}

type AlertV2Change struct {
	AlertV2Common
	Config                                   AlertV2ConfigChange `json:"config"`
	UnreportedAlertNotificationsRetentionSec *int                `json:"unreportedAlertNotificationsRetentionSec"`
}

type alertV2ChangeWrapper struct {
	Alert AlertV2Change `json:"alert"`
}

type CloudAccountCredentialsMonitor struct {
	AccountId   string `json:"accountId"`
	RoleName    string `json:"roleName"`
	SecretKey   string `json:"key"`
	AccessKeyId string `json:"id"`
}

type CloudAccountMonitor struct {
	Id                int                            `json:"id"`
	Platform          string                         `json:"platform"`
	IntegrationType   string                         `json:"integrationType"`
	Credentials       CloudAccountCredentialsMonitor `json:"credentials"`
	AdditionalOptions string                         `json:"additionalOptions"`
}

type CloudAccountMonitorForCost struct {
	Feature       string                         `json:"feature"`
	Platform      string                         `json:"platform"`
	Configuration CloudCostConfiguration         `json:"config"`
	Credentials   CloudAccountCredentialsMonitor `json:"credentials"`
}

type CloudCostConfiguration struct {
	AthenaBucketName     string `json:"athenaBucketName"`
	AthenaDatabaseName   string `json:"athenaDatabaseName"`
	AthenaRegion         string `json:"athenaRegion"`
	AthenaWorkgroup      string `json:"athenaWorkgroup"`
	AthenaTableName      string `json:"athenaTableName"`
	SpotPricesBucketName string `json:"spotPricesBucketName"`
}

type CloudAccountCreatedForCost struct {
	Id              string `json:"id"`
	CustomerId      int    `json:"customerId"`
	ProviderId      string `json:"providerId"`
	Provider        string `json:"provider"`
	SkipFetch       bool   `json:"skipFetch"`
	IntegrationType string `json:"integrationType"`
	CredentialsType string `json:"credentialsType"`
	RoleArn         string `json:"roleArn"`
	ExternalId      string `json:"externalId"`
}

type cloudAccountWrapperMonitor struct {
	CloudAccount CloudAccountMonitor `json:"provider"`
}

type CloudConfigForCost struct {
	AthenaProjectId      string `json:"athenaProjectId"`
	AthenaBucketName     string `json:"athenaBucketName"`
	AthenaRegion         string `json:"athenaRegion"`
	AthenaDatabaseName   string `json:"athenaDatabaseName"`
	AthenaTableName      string `json:"athenaTableName"`
	AthenaWorkgroup      string `json:"athenaWorkgroup"`
	SpotPricesBucketName string `json:"spotPricesBucketName"`
	IntegrationType      string `json:"integrationType"`
}

type CloudAccountCostProvider struct {
	CustomerId      int                `json:"customerId"`
	ProviderId      string             `json:"providerId"`
	Provider        string             `json:"provider"`
	CredentialsId   string             `json:"credentialsId"`
	Feature         string             `json:"feature"`
	Config          CloudConfigForCost `json:"config"`
	Enabled         bool               `json:"enabled"`
	CredentialsType string             `json:"credentialsType"`
	RoleArn         string             `json:"roleArn"`
	ExternalId      string             `json:"externalId"`
}

type CloudAccountCostProviderWrapper struct {
	CloudAccountCostProvider CloudAccountCostProvider `json:"item"`
}

type PosturePolicyZoneMeta struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PosturePolicy struct {
	ID             string                  `json:"id,omitempty"`
	Name           string                  `json:"name,omitempty"`
	Type           int                     `json:"type,omitempty"`
	Kind           int                     `json:"kind,omitempty"`
	Description    string                  `json:"description,omitempty"`
	Version        string                  `json:"version,omitempty"`
	Link           string                  `json:"link,omitempty"`
	Authors        string                  `json:"authors,omitempty"`
	PublishedData  string                  `json:"publishedDate,omitempty"`
	MinKubeVersion float64                 `json:"minKubeVersion,omitempty"`
	MaxKubeVersion float64                 `json:"maxKubeVersion,omitempty"`
	IsCustom       bool                    `json:"isCustom,omitempty"`
	IsActive       bool                    `json:"isActive,omitempty"`
	Platform       string                  `json:"platform,omitempty"`
	Zones          []PosturePolicyZoneMeta `json:"zones,omitempty"`
}

type FullPosturePolicy struct {
	ID                 string              `json:"id,omitempty"`
	Name               string              `json:"name,omitempty"`
	Type               string              `json:"type,omitempty"`
	Description        string              `json:"description,omitempty"`
	Version            string              `json:"version,omitempty"`
	Link               string              `json:"link,omitempty"`
	Authors            string              `json:"authors,omitempty"`
	PublishedData      string              `json:"publishedDate,omitempty"`
	RequirementsGroup  []RequirementsGroup `json:"requirementFolders,omitempty"`
	MinKubeVersion     float64             `json:"minKubeVersion,omitempty"`
	MaxKubeVersion     float64             `json:"maxKubeVersion,omitempty"`
	IsCustom           bool                `json:"isCustom,omitempty"`
	IsActive           bool                `json:"isActive,omitempty"`
	Platform           string              `json:"platform,omitempty"`
	VersionConstraints []VersionConstraint `json:"targets,omitempty"`
}

type VersionConstraint struct {
	Platform   string  `json:"platform"`
	MinVersion float64 `json:"minVersion,omitempty"`
	MaxVersion float64 `json:"maxVersion,omitempty"`
}

type RequirementsGroup struct {
	ID                        string              `json:"id,omitempty"`
	Name                      string              `json:"name,omitempty"`
	Requirements              []Requirement       `json:"requirements,omitempty"`
	Description               string              `json:"description,omitempty"`
	Authors                   string              `json:"author,omitempty"`
	Folders                   []RequirementsGroup `json:"folders,omitempty"`
	RequirementFolderParentID string              `json:"requirementFolderParentId,omitempty"`
}

type Requirement struct {
	ID                  string    `json:"id,omitempty"`
	Name                string    `json:"name,omitempty"`
	RequirementFolderId string    `json:"requirementFolderId,omitempty"`
	Description         string    `json:"description,omitempty"`
	Controls            []Control `json:"controls,omitempty"`
	Authors             string    `json:"authors,omitempty"`
}

type Control struct {
	Name   string `json:"name,omitempty"`
	Status bool   `json:"status,omitempty"`
}

type CreatePosturePolicy struct {
	ID                 string                    `json:"id,omitempty"`
	Name               string                    `json:"name,omitempty"`
	Description        string                    `json:"description,omitempty"`
	Type               string                    `json:"type,omitempty"`
	Link               string                    `json:"link,omitempty"`
	Version            string                    `json:"version,omitempty"`
	RequirementGroups  []CreateRequirementsGroup `json:"groups,omitempty"`
	MinKubeVersion     float64                   `json:"minKubeVersion,omitempty"`
	MaxKubeVersion     float64                   `json:"maxKubeVersion,omitempty"`
	IsActive           bool                      `json:"isActive,omitempty"`
	Platform           string                    `json:"platform,omitempty"`
	VersionConstraints []VersionConstraint       `json:"targets,omitempty"`
}

type CreateRequirementsGroup struct {
	ID           string                    `json:"id,omitempty"`
	Name         string                    `json:"name,omitempty"`
	Requirements []CreateRequirement       `json:"requirements,omitempty"`
	Description  string                    `json:"description,omitempty"`
	Folders      []CreateRequirementsGroup `json:"groups,omitempty"`
}

type CreateRequirement struct {
	ID          string                     `json:"id,omitempty"`
	Name        string                     `json:"name,omitempty"`
	Description string                     `json:"description,omitempty"`
	Controls    []CreateRequirementControl `json:"controls,omitempty"`
}

type CreateRequirementControl struct {
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}
type PosturePolicyResponse struct {
	Data PosturePolicy `json:"data"`
}

type FullPosturePolicyResponse struct {
	Data FullPosturePolicy `json:"data"`
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

type InhibitionRule struct {
	Name           string          `json:"name,omitempty"`
	Description    string          `json:"description,omitempty"`
	Enabled        bool            `json:"isEnabled"`
	SourceMatchers []LabelMatchers `json:"sourceMatchers"`
	TargetMatchers []LabelMatchers `json:"targetMatchers"`
	Equal          []string        `json:"equal,omitempty"`

	Version int `json:"version,omitempty"`
	ID      int `json:"id,omitempty"`
}

type LabelMatchers struct {
	LabelName string `json:"labelName"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
}

type AgentAccessKey struct {
	ID             int               `json:"id,omitempty"`
	Reservation    int               `json:"agentReservation,omitempty"`
	Limit          int               `json:"agentLimit,omitempty"`
	TeamID         int               `json:"teamId,omitempty"`
	AgentAccessKey string            `json:"accessKey,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	Enabled        bool              `json:"isEnabled"`
	DateCreated    string            `json:"dateCreated,omitempty"`
	DateDisabled   string            `json:"dateDisabled,omitempty"`
}

type AgentAccessKeyReadWrapper struct {
	CustomerAccessKey []AgentAccessKey `json:"customerAccessKeys"`
}

type AgentAccessKeyWriteWrapper struct {
	CustomerAccessKey AgentAccessKey `json:"customerAccessKey"`
}

type OrganizationSecure struct {
	cloudauth.CloudOrganization
}
