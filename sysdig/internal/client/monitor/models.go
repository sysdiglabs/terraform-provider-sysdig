package monitor

import (
	"bytes"
	"encoding/json"
	"io"
)

type CloudAccountCredentials struct {
	AccountId string `json:"accountId"`
}

type CloudAccount struct {
	Id                int                     `json:"id"`
	Platform          string                  `json:"platform"`
	IntegrationType   string                  `json:"integrationType"`
	Credentials       CloudAccountCredentials `json:"credentials"`
	AdditionalOptions string                  `json:"additionalOptions"`
}

type cloudAccountWrapper struct {
	CloudAccount CloudAccount `json:"provider"`
}

func CloudAccountFromJSON(body []byte) *CloudAccount {
	var result cloudAccountWrapper
	_ = json.Unmarshal(body, &result)

	return &result.CloudAccount
}

func CloudAccountToJSON(providerInfo *CloudAccount) io.Reader {
	payload, _ := json.Marshal(*providerInfo)
	return bytes.NewBuffer(payload)
}

// -------- Notification Channels --------

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
	Options NotificationChannelOptions `json:"options"`
}

func NotificationChannelFromJSON(body []byte) NotificationChannel {
	var result notificationChannelWrapper
	_ = json.Unmarshal(body, &result)

	return result.NotificationChannel
}

type notificationChannelWrapper struct {
	NotificationChannel NotificationChannel `json:"notificationChannel"`
}
