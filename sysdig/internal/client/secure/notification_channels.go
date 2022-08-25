package secure

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (client *sysdigSecureClient) GetNotificationChannelById(ctx context.Context, id int) (nc NotificationChannel, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.GetNotificationChannelUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	nc = NotificationChannelFromJSON(body)

	if nc.ID == 0 {
		err = fmt.Errorf("NotificationChannel with ID: %d does not exists", id)
		return
	}
	return
}

func (client *sysdigSecureClient) GetNotificationChannelByName(ctx context.Context, name string) (nc NotificationChannel, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.GetNotificationChannelsUrl(), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
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

func (client *sysdigSecureClient) CreateNotificationChannel(ctx context.Context, ncRequest NotificationChannel) (nc NotificationChannel, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPost, client.GetNotificationChannelsUrl(), ncRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	nc = NotificationChannelFromJSON(body)
	return
}

func (client *sysdigSecureClient) UpdateNotificationChannel(ctx context.Context, ncRequest NotificationChannel) (nc NotificationChannel, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPut, client.GetNotificationChannelUrl(ncRequest.ID), ncRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	nc = NotificationChannelFromJSON(body)
	return
}

func (client *sysdigSecureClient) DeleteNotificationChannel(ctx context.Context, id int) error {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodDelete, client.GetNotificationChannelUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}
	return nil
}

func (client *sysdigSecureClient) GetNotificationChannelsUrl() string {
	return fmt.Sprintf("%s/api/notificationChannels", client.URL)
}

func (client *sysdigSecureClient) GetNotificationChannelUrl(id int) string {
	return fmt.Sprintf("%s/api/notificationChannels/%d", client.URL, id)
}
