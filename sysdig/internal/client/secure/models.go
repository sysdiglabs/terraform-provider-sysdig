package secure

import (
	"bytes"
	"encoding/json"
	"io"
)

// -------- Policies --------

type Policy struct {
	ID                     int      `json:"id,omitempty"`
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	Severity               int      `json:"severity"`
	Enabled                bool     `json:"enabled"`
	RuleNames              []string `json:"ruleNames"`
	Actions                []Action `json:"actions"`
	Scope                  string   `json:"scope,omitempty"`
	Version                int      `json:"version,omitempty"`
	NotificationChannelIds []int    `json:"notificationChannelIds"`
	Type                   string   `json:"type"`
	Runbook                string   `json:"runbook"`
}

type Action struct {
	AfterEventNs         int    `json:"afterEventNs,omitempty"`
	BeforeEventNs        int    `json:"beforeEventNs,omitempty"`
	Name                 string `json:"name,omitempty"`
	IsLimitedToContainer bool   `json:"isLimitedToContainer"`
	Type                 string `json:"type"`
}

func (policy *Policy) ToJSON() io.Reader {
	payload, _ := json.Marshal(policy)
	return bytes.NewBuffer(payload)
}

func PolicyFromJSON(body []byte) (result Policy) {
	_ = json.Unmarshal(body, &result)

	return result
}

// -------- User Rules File --------

type UserRulesFile struct {
	Content string `json:"content"`
	Version int    `json:"version"`
}

type userRulesFileWrapper struct {
	UserRulesFile UserRulesFile `json:"userRulesFile"`
}

func (userRulesFile *UserRulesFile) ToJSON() io.Reader {
	payload, _ := json.Marshal(userRulesFileWrapper{*userRulesFile})
	return bytes.NewBuffer(payload)
}

func UserRulesFileFromJSON(body []byte) UserRulesFile {
	var result userRulesFileWrapper
	_ = json.Unmarshal(body, &result)

	return result.UserRulesFile
}

// -------- Notification Channels --------

type NotificationChannelOptions struct {
	EmailRecipients []string `json:"emailRecipients,omitempty"` // Type: email
	SnsTopicARNs    []string `json:"snsTopicARNs,omitempty"`    // Type: SNS
	APIKey          string   `json:"apiKey,omitempty"`          // Type: VictorOps, OpsGenie
	RoutingKey      string   `json:"routingKey,omitempty"`      // Type: VictorOps
	Url             string   `json:"url,omitempty"`             // Type: OpsGenie, Webhook and Slack
	Channel         string   `json:"channel,omitempty"`         // Type: Slack
	Account         string   `json:"account,omitempty"`         // Type: PagerDuty
	ServiceKey      string   `json:"serviceKey,omitempty"`      // Type: PagerDuty
	ServiceName     string   `json:"serviceName,omitempty"`     // Type: PagerDuty
	Region          string   `json:"region,omitempty"`          // Type: OpsGenie

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
	_ = json.Unmarshal(body, &result)

	return result.NotificationChannel
}

func NotificationChannelListFromJSON(body []byte) []NotificationChannel {
	var result notificationChannelListWrapper
	_ = json.Unmarshal(body, &result)

	return result.NotificationChannels
}

type notificationChannelListWrapper struct {
	NotificationChannels []NotificationChannel `json:"notificationChannels"`
}

type notificationChannelWrapper struct {
	NotificationChannel NotificationChannel `json:"notificationChannel"`
}

// -------- Rules --------

type Rule struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Details     Details  `json:"details"`
	Version     int      `json:"version,omitempty"`
}

type Details struct {
	// Containers
	Containers *Containers `json:"containers,omitempty"`

	// Filesystems
	ReadWritePaths *ReadWritePaths `json:"readWritePaths,omitempty"`
	ReadPaths      *ReadPaths      `json:"readPaths,omitempty"`

	// Network
	AllOutbound    bool            `json:"allOutbound,omitempty"`
	AllInbound     bool            `json:"allInbound,omitempty"`
	TCPListenPorts *TCPListenPorts `json:"tcpListenPorts,omitempty"`
	UDPListenPorts *UDPListenPorts `json:"udpListenPorts,omitempty"`

	// Processes
	Processes *Processes `json:"processes,omitempty"`

	// Syscalls
	Syscalls *Syscalls `json:"syscalls,omitempty"`

	// Falco
	Append     *bool        `json:"append,omitempty"`
	Source     string       `json:"source,omitempty"`
	Output     string       `json:"output"`
	Condition  *Condition   `json:"condition,omitempty"`
	Priority   string       `json:"priority,omitempty"`
	Exceptions []*Exception `json:"exceptions,omitempty"`

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

func (r *Rule) ToJSON() io.Reader {
	payload, _ := json.Marshal(r)
	return bytes.NewBuffer(payload)
}

func RuleFromJSON(body []byte) (rule Rule, err error) {
	err = json.Unmarshal(body, &rule)
	return
}

// -------- List --------

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

func (l *List) ToJSON() io.Reader {
	payload, _ := json.Marshal(l)
	return bytes.NewBuffer(payload)
}

func ListFromJSON(body []byte) (list List, err error) {
	err = json.Unmarshal(body, &list)
	return
}

// -------- Macro -------

type Macro struct {
	ID        int            `json:"id,omitempty"`
	Version   int            `json:"version,omitempty"`
	Name      string         `json:"name"`
	Condition MacroCondition `json:"condition"`
	Append    bool           `json:"append"`
}

type MacroCondition struct {
	Condition string `json:"condition"`
}

func (l *Macro) ToJSON() io.Reader {
	payload, _ := json.Marshal(l)
	return bytes.NewBuffer(payload)
}

func MacroFromJSON(body []byte) (macro Macro, err error) {
	err = json.Unmarshal(body, &macro)
	return
}

// -------- Team --------
type Team struct {
	ID                  int         `json:"id,omitempty"`
	Version             int         `json:"version,omitempty"`
	Theme               string      `json:"theme"`
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	ScopeBy             string      `json:"show"`
	Filter              string      `json:"filter"`
	CanUseSysdigCapture bool        `json:"canUseSysdigCapture"`
	UserRoles           []UserRoles `json:"userRoles,omitempty"`
	DefaultTeam         bool        `json:"default"`
	Products            []string    `json:"products"`
}

type UserRoles struct {
	UserId int    `json:"userId"`
	Email  string `json:"userName,omitempty"`
	Role   string `json:"role"`
	Admin  bool   `json:"admin,omitempty"`
}

func (t *Team) ToJSON() io.Reader {
	payload, _ := json.Marshal(*t)
	return bytes.NewBuffer(payload)
}

func TeamFromJSON(body []byte) Team {
	var result teamWrapper
	_ = json.Unmarshal(body, &result)

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
	_ = json.Unmarshal(body, &result)

	return result.UsersList
}

type usersListWrapper struct {
	UsersList []UsersList `json:"users"`
}

// -------- VulnerabilityExceptionList --------

type VulnerabilityExceptionList struct {
	ID      string `json:"id,omitempty"`
	Version string `json:"version"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

func (l *VulnerabilityExceptionList) ToJSON() io.Reader {
	payload, _ := json.Marshal(*l)
	return bytes.NewBuffer(payload)
}

func VulnerabilityExceptionListFromJSON(body []byte) *VulnerabilityExceptionList {
	var result VulnerabilityExceptionList
	_ = json.Unmarshal(body, &result)

	return &result
}

// -------- VulnerabilityException --------

type VulnerabilityException struct {
	ID             string `json:"id"`
	Gate           string `json:"gate"`
	TriggerID      string `json:"trigger_id"`
	Notes          string `json:"notes"`
	ExpirationDate *int   `json:"expiration_date,omitempty"`
	Enabled        bool   `json:"enabled"`
}

func (e *VulnerabilityException) ToJSON() io.Reader {
	payload, _ := json.Marshal(*e)
	return bytes.NewBuffer(payload)
}

func VulnerabilityExceptionFromJSON(body []byte) *VulnerabilityException {
	var result VulnerabilityException
	_ = json.Unmarshal(body, &result)

	return &result
}

// -------- CloudAccount --------

type CloudAccount struct {
	AccountID                    string `json:"accountId"`
	Provider                     string `json:"provider"`
	Alias                        string `json:"alias"`
	RoleAvailable                bool   `json:"roleAvailable"`
	RoleName                     string `json:"roleName"`
	ExternalID                   string `json:"externalId,omitempty"`
	WorkLoadIdentityAccountID    string `json:"workloadIdentityAccountId,omitempty"`
	WorkLoadIdentityAccountAlias string `json:"workLoadIdentityAccountAlias,omitempty"`
}

func (e *CloudAccount) ToJSON() io.Reader {
	payload, _ := json.Marshal(*e)
	return bytes.NewBuffer(payload)
}

func CloudAccountFromJSON(body []byte) *CloudAccount {
	var result CloudAccount
	_ = json.Unmarshal(body, &result)

	return &result
}

// -------- BenchmarkTask --------

var SupportedBenchmarkTaskSchemas = []string{
	"aws_foundations_bench-1.3.0",
	"gcp_foundations_bench-1.2.0",
	"azure_foundations_bench-1.3.0",
}

type BenchmarkTask struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Schema   string `json:"schema"`
	Scope    string `json:"scope"`
	Schedule string `json:"schedule"`
	Enabled  bool   `json:"enabled"`
}

func (t *BenchmarkTask) ToJSON() io.Reader {
	payload, _ := json.Marshal(*t)
	return bytes.NewBuffer(payload)
}

func BenchmarkTaskFromJSON(body []byte) *BenchmarkTask {
	var result BenchmarkTask
	_ = json.Unmarshal(body, &result)

	return &result
}

// -------- Scanning Policies --------
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

func (policy *ScanningPolicy) ToJSON() io.Reader {
	payload, _ := json.Marshal(policy)
	return bytes.NewBuffer(payload)
}

func ScanningPolicyFromJSON(body []byte) (result ScanningPolicy) {
	_ = json.Unmarshal(body, &result)
	return result
}

// -------- Scanning Policy Assignments --------
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

func (policy *ScanningPolicyAssignmentList) ToJSON() io.Reader {
	payload, _ := json.Marshal(policy)
	return bytes.NewBuffer(payload)
}

func ScanningPolicyAssignmentFromJSON(body []byte) (result ScanningPolicyAssignmentList) {
	_ = json.Unmarshal(body, &result)
	return result
}
