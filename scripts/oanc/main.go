package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	slack     = "SLACK"
	pagerduty = "PAGER_DUTY"
)

type NotificationChannelList []NotificationChannel

func (nc NotificationChannelList) String() string {
	slackResourceTemplate := `resource "sysdig_secure_notification_channel" "example-slack-%d" {
  name = "%s"
  enabled = %t
  type = "%s"
  url = "%s"
  channel = "%s"
  notify_when_ok = %t
  notify_when_resolved = %t
}
`

	pdResourceTemplate := `resource "sysdig_secure_notification_channel" "example-pager-duty-%d" {
  name = "%s"
  enabled = %t
  type = "%s"
  account = "%s"
  service_key = "%s"
  service_name = "%s"
  notify_when_ok = %t
  notify_when_resolved = %t
}
`

	var ncStr string
	builder := strings.Builder{}
	for id, channel := range nc {
		switch channel.Type {
		case slack:
			ncStr = fmt.Sprintf(slackResourceTemplate, id, channel.Name, channel.Enabled, channel.Type, channel.Options.URL, channel.Options.Channel, channel.Options.NotifyOnOk, channel.Options.NotifyOnResolve)
		case pagerduty:
			ncStr = fmt.Sprintf(pdResourceTemplate, id, channel.Name, channel.Enabled, channel.Type, channel.Options.Account, channel.Options.ServiceKey, channel.Options.ServiceName, channel.Options.NotifyOnOk, channel.Options.NotifyOnResolve)
		}
		builder.WriteString(ncStr)
	}
	return builder.String()
}

type NotificationChannel struct {
	ID                   int    `json:"id"`
	Version              int    `json:"version"`
	CreatedOn            int64  `json:"createdOn"`
	ModifiedOn           int64  `json:"modifiedOn"`
	Type                 string `json:"type"`
	Enabled              bool   `json:"enabled"`
	SendTestNotification bool   `json:"sendTestNotification"`
	Name                 string `json:"name"`
	Options              struct {
		NotifyOnOk      bool   `json:"notifyOnOk"`
		NotifyOnResolve bool   `json:"notifyOnResolve"`
		Account         string `json:"account,omitempty"`     // Type: PagerDuty
		ServiceKey      string `json:"serviceKey,omitempty"`  // Type: PagerDuty
		ServiceName     string `json:"serviceName,omitempty"` // Type: PagerDuty
		Channel         string `json:"channel,omitempty"`     // Type: Slack
		URL             string `json:"url,omitempty"`         // Type: Slack
	} `json:"options"`
}

type Response struct {
	NotificationChannels []NotificationChannel `json:"notificationChannels"`
}

func main() {
	token, tokenDefined := os.LookupEnv("SDC_TOKEN")
	if !tokenDefined {
		log.Fatal("SDC_TOKEN env var not set")
	}
	url, urlDefined := os.LookupEnv("SDC_URL")
	if !urlDefined {
		url = "https://secure.sysdig.com"
	}

	endpoint := fmt.Sprintf("%s/api/notificationChannels", url)
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatal(errors.New(response.Status))
	}

	var data Response
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	var slackNC NotificationChannelList
	for _, nc := range data.NotificationChannels {

		fmt.Printf("%+v\n", nc)
		if nc.Type == "SLACK" {
			slackNC = append(slackNC, nc)
		}

		if nc.Type == "PAGER_DUTY" {
			slackNC = append(slackNC, nc)
		}
	}

	fmt.Println(slackNC)
}
