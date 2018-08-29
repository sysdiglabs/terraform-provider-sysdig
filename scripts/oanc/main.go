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

type NotificationChannelList []NotificationChannel

func (nc NotificationChannelList) String() string {
	resourceTemplate := `resource "sysdig_secure_notification_channel" "example-slack-%d" {
  name = "%s"
  enabled = %t
  type = "%s"
  url = "%s"
  channel = "%s"
  notify_when_ok = %t
  notify_when_resolved = %t
}
`
	builder := strings.Builder{}
	for id, slack := range nc {
		ncStr := fmt.Sprintf(resourceTemplate, id, slack.Name, slack.Enabled, slack.Type, slack.Options.URL, slack.Options.Channel, slack.Options.NotifyOnOk, slack.Options.NotifyOnResolve)
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
		URL             string `json:"url"`
		Channel         string `json:"channel"`
		NotifyOnResolve bool   `json:"notifyOnResolve"`
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
		if nc.Type == "SLACK" {
			slackNC = append(slackNC, nc)
		}
	}

	fmt.Println(slackNC)
}
