package monitor

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigMonitorClient) GetNotificationChannelById(ctx context.Context, id int) (nc NotificationChannel, err error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodGet, client.GetNotificationChannelUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	nc = NotificationChannelFromJSON(body)

	if nc.Version == 0 {
		err = fmt.Errorf("NotificationChannel with ID: %d does not exists", id)
		return
	}
	return
}

func (client *sysdigMonitorClient) GetNotificationChannelByName(ctx context.Context, name string) (nc NotificationChannel, err error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodGet, client.GetNotificationChannelsUrl(), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	ncList := NotificationChannelListFromJSON(body)

	for _, channel := range ncList {
		if channel.Name == name {
			nc = channel
			return
		}
	}

	err = fmt.Errorf("Notification channel with Name: %s does not exist", name)
	return
}

func (client *sysdigMonitorClient) CreateNotificationChannel(ctx context.Context, ncRequest NotificationChannel) (nc NotificationChannel, err error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodPost, client.GetNotificationChannelsUrl(), ncRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errorFromResponse(response)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	nc = NotificationChannelFromJSON(body)
	return
}

func (client *sysdigMonitorClient) UpdateNotificationChannel(ctx context.Context, ncRequest NotificationChannel) (nc NotificationChannel, err error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodPut, client.GetNotificationChannelUrl(ncRequest.ID), ncRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	nc = NotificationChannelFromJSON(body)
	return
}

func (client *sysdigMonitorClient) DeleteNotificationChannel(ctx context.Context, id int) error {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodDelete, client.GetNotificationChannelUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return errorFromResponse(response)
	}
	return nil
}

func (client *sysdigMonitorClient) GetNotificationChannelsUrl() string {
	return fmt.Sprintf("%s/api/notificationChannels", client.URL)
}

func (client *sysdigMonitorClient) GetNotificationChannelUrl(id int) string {
	return fmt.Sprintf("%s/api/notificationChannels/%d", client.URL, id)
}
