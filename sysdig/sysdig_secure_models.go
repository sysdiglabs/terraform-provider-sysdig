package sysdig

import (
	"bytes"
	"encoding/json"
	"io"
)

// -------- Policies --------

type Policy struct {
	ID                           int                           `json:"id,omitempty"`
	Name                         string                        `json:"name"`
	Description                  string                        `json:"description"`
	Severity                     int                           `json:"severity"`
	ContainerScope               bool                          `json:"containerScope"`
	HostScope                    bool                          `json:"hostScope"`
	Enabled                      bool                          `json:"enabled"`
	Actions                      []Action                      `json:"actions,omitempty"`
	Scope                        string                        `json:"scope,omitempty"`
	FalcoConfiguration           FalcoConfiguration            `json:"falcoConfiguration,omitempty"`
	ProcessesConfiguration       *ProcessesConfiguration       `json:"processesConfiguration,omitempty"`
	ContainerImagesConfiguration *ContainerImagesConfiguration `json:"containerImagesConfiguration,omitempty"`
	FileSystemConfiguration      *FileSystemConfiguration      `json:"fileSystemConfiguration,omitempty"`
	SyscallsConfiguration        *SyscallsConfiguration        `json:"syscallsConfiguration,omitempty"`
	NetworkConfiguration         *NetworkConfiguration         `json:"networkConfiguration,omitempty"`
	Version                      int                           `json:"version,omitempty"`
	NotificationChannelIds       []int                         `json:"notificationChannelIds,omitempty"`
}

type Action struct {
	AfterEventNs         int    `json:"afterEventNs,omitempty"`
	BeforeEventNs        int    `json:"beforeEventNs,omitempty"`
	IsLimitedToContainer bool   `json:"isLimitedToContainer,omitempty"`
	Type                 string `json:"type"`
}

type FalcoConfiguration struct {
	RuleNameRegEx string `json:"ruleNameRegEx"`
}

type ProcessesConfigurationElement struct {
	Values  []string `json:"values"`
	OnMatch string   `json:"onMatch"`
}

type ProcessesConfiguration struct {
	OnDefault string                          `json:"onDefault"`
	Fields    []string                        `json:"fields"`
	List      []ProcessesConfigurationElement `json:"list,omitempty"`
}

type ContainerImagesConfigurationElement struct {
	Values  []string `json:"values"`
	OnMatch string   `json:"onMatch"`
}

type ContainerImagesConfiguration struct {
	OnDefault string                                `json:"onDefault"`
	Fields    []string                              `json:"fields"`
	List      []ContainerImagesConfigurationElement `json:"list"`
}

type FileSystemConfigurationListElement struct {
	Values     []string `json:"values"`
	OnMatch    string   `json:"onMatch"`
	AccessType string   `json:"accessType"`
}

type FileSystemConfiguration struct {
	OnDefault string                               `json:"onDefault"`
	Fields    []string                             `json:"fields"`
	List      []FileSystemConfigurationListElement `json:"list"`
}

type SyscallsConfigurationElement struct {
	Values  []string `json:"values"`
	OnMatch string   `json:"onMatch"`
}

type SyscallsConfiguration struct {
	OnDefault string                         `json:"onDefault"`
	Fields    []string                       `json:"fields"`
	List      []SyscallsConfigurationElement `json:"list"`
}

type NetworkConfigurationListeningPortsListElement struct {
	Values   []string `json:"values"`
	OnMatch  string   `json:"onMatch"`
	NetProto string   `json:"netProto"`
}

type NetworkConfiguration struct {
	ListeningPorts struct {
		OnDefault string                                          `json:"onDefault"`
		Fields    []string                                        `json:"fields"`
		List      []NetworkConfigurationListeningPortsListElement `json:"list"`
	} `json:"listeningPorts"`
	Inbound struct {
		OnDefault string   `json:"onDefault"`
		Fields    []string `json:"fields"`
	} `json:"inbound"`
	Outbound struct {
		OnDefault string   `json:"onDefault"`
		Fields    []string `json:"fields"`
	} `json:"outbound"`
}

type policyWrapper struct {
	Policy Policy `json:"policy"`
}

func (policy *Policy) ToJSON() io.Reader {
	payload, _ := json.Marshal(policyWrapper{*policy})
	return bytes.NewBuffer(payload)
}

func PolicyFromJSON(body []byte) Policy {
	var result policyWrapper
	json.Unmarshal(body, &result)

	return result.Policy
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

// -------- Policies Priority ---------

type PoliciesPriority struct {
	Version   int   `json:"version,omitempty"`
	PolicyIds []int `json:"policyIds"`
}

type policiesPriorityWrapper struct {
	Priorities PoliciesPriority `json:"priorities"`
}

func (pp *PoliciesPriority) ToJSON() io.Reader {
	payload, _ := json.Marshal(policiesPriorityWrapper{*pp})
	return bytes.NewBuffer(payload)
}

func PoliciesPriorityFromJSON(body []byte) PoliciesPriority {
	var result policiesPriorityWrapper
	json.Unmarshal(body, &result)
	return result.Priorities
}

// -------- Notification Channels --------

type NotificationChannelOptions struct {
	EmailRecipients []string `json:"emailRecipients,omitempty"` // Type: email
	SnsTopicARNs    []string `json:"snsTopicARNs,omitempty"`    // Type: SNS
	APIKey          string   `json:"apiKey,omitempty"`          // Type: VictorOps
	RoutingKey      string   `json:"routingKey,omitempty"`      // Type: VictorOps
	Url             string   `json:"url,omitempty"`             // Type: OpsGenie, Webhook and Slack
	Channel         string   `json:"channel,omitempty"`         // Type: Slack

	NotifyOnOk      bool `json:"notifyOnOk"`
	NotifyOnResolve bool `json:"notifyOnResolve"`
}

type NotificationChannel struct {
	ID      int                        `json:"id,omitempty"`
	Version int                        `json:"version,omitempty"`
	Type    string                     `json:"type"`
	Name    string                     `json:"name"`
	Enabled bool                       `json:"enabled"`
	Options NotificationChannelOptions `json:"options"`
}

func (n NotificationChannel) ToJSON() io.Reader {
	payload, _ := json.Marshal(notificationChannelWrapper{n})
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
