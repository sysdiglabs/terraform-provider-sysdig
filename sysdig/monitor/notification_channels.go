package monitor

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigMonitorClient) GetNotificationChannelById(id int) (nc NotificationChannel, err error) {
	response, err := client.doSysdigMonitorRequest(http.MethodGet, client.GetNotificationChannelUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	nc = NotificationChannelFromJSON(body)

	if nc.Version == 0 {
		err = fmt.Errorf("NotificationChannel with ID: %d does not exists", id)
		return
	}
	return
}

func (client *sysdigMonitorClient) GetNotificationChannelByName(name string) (nc NotificationChannel, err error) {
	response, err := client.doSysdigMonitorRequest(http.MethodGet, client.GetNotificationChannelsUrl(), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
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

func (client *sysdigMonitorClient) CreateNotificationChannel(ncRequest NotificationChannel) (nc NotificationChannel, err error) {
	response, err := client.doSysdigMonitorRequest(http.MethodPost, client.GetNotificationChannelsUrl(), ncRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errors.New(response.Status)
		return
	}

	nc = NotificationChannelFromJSON(body)
	return
}

func (client *sysdigMonitorClient) UpdateNotificationChannel(ncRequest NotificationChannel) (nc NotificationChannel, err error) {
	response, err := client.doSysdigMonitorRequest(http.MethodPut, client.GetNotificationChannelUrl(ncRequest.ID), ncRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	nc = NotificationChannelFromJSON(body)
	return
}

func (client *sysdigMonitorClient) DeleteNotificationChannel(id int) error {
	response, err := client.doSysdigMonitorRequest(http.MethodDelete, client.GetNotificationChannelUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}
	return nil
}

func (client *sysdigMonitorClient) GetNotificationChannelsUrl() string {
	return fmt.Sprintf("%s/api/notificationChannels", client.URL)
}

func (client *sysdigMonitorClient) GetNotificationChannelUrl(id int) string {
	return fmt.Sprintf("%s/api/notificationChannels/%d", client.URL, id)
}
