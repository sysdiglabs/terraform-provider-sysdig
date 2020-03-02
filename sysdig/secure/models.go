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
}

type Action struct {
	AfterEventNs         int    `json:"afterEventNs,omitempty"`
	BeforeEventNs        int    `json:"beforeEventNs,omitempty"`
	IsLimitedToContainer bool   `json:"isLimitedToContainer"`
	Type                 string `json:"type"`
}

func (policy *Policy) ToJSON() io.Reader {
	payload, _ := json.Marshal(policy)
	return bytes.NewBuffer(payload)
}

func PolicyFromJSON(body []byte) (result Policy) {
	json.Unmarshal(body, &result)

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
	json.Unmarshal(body, &result)

	return result.UserRulesFile
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
	Append    bool       `json:"append,omitempty"`
	Source    string     `json:"source,omitempty"`
	Output    string     `json:"output,omitempty"`
	Condition *Condition `json:"condition,omitempty"`
	Priority  string     `json:"priority,omitempty"`

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

func (r *Rule) ToJSON() io.Reader {
	payload, _ := json.Marshal(r)
	return bytes.NewBuffer(payload)
}

func RuleFromJSON(body []byte) (rule Rule, err error) {
	err = json.Unmarshal(body, &rule)
	return
}

// -------- User --------
type User struct {
	ID         int    `json:"id,omitempty"`
	Version    int    `json:"version,omitempty"`
	SystemRole string `json:"systemRole,omitempty"`
	Email      string `json:"username"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
}

func (u *User) ToJSON() io.Reader {
	payload, _ := json.Marshal(*u)
	return bytes.NewBuffer(payload)
}

func UsersFromJSON(body []byte) User {
	var result usersWrapper
	json.Unmarshal(body, &result)

	return result.Users
}

type usersWrapper struct {
	Users User `json:"user"`
}
