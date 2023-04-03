package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	GetNotificationChannels = "%s/api/notificationChannels"
	GetNotificationChannel  = "%s/api/notificationChannels/%d"
)

type NotificationChannelInterface interface {
	GetNotificationChannelById(ctx context.Context, id int) (NotificationChannel, error)
	GetNotificationChannelByName(ctx context.Context, name string) (NotificationChannel, error)
	CreateNotificationChannel(ctx context.Context, channel NotificationChannel) (NotificationChannel, error)
	UpdateNotificationChannel(ctx context.Context, channel NotificationChannel) (NotificationChannel, error)
	DeleteNotificationChannel(ctx context.Context, id int) error
}

func (client *Client) GetNotificationChannelById(ctx context.Context, id int) (NotificationChannel, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetNotificationChannelUrl(id), nil)
	if err != nil {
		return NotificationChannel{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return NotificationChannel{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[notificationChannelWrapper](response.Body)
	if err != nil {
		return NotificationChannel{}, err
	}

	if wrapper.NotificationChannel.ID == 0 {
		return NotificationChannel{}, fmt.Errorf("NotificationChannel with ID: %d does not exists", id)
	}

	return wrapper.NotificationChannel, nil
}

func (client *Client) GetNotificationChannelByName(ctx context.Context, name string) (NotificationChannel, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetNotificationChannelsUrl(), nil)
	if err != nil {
		return NotificationChannel{}, nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return NotificationChannel{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[notificationChannelListWrapper](response.Body)
	if err != nil {
		return NotificationChannel{}, err
	}

	for _, channel := range wrapper.NotificationChannels {
		if channel.Name == name {
			return channel, nil
		}
	}

	return NotificationChannel{}, fmt.Errorf("notification channel with name: %s does not exist", name)
}

func (client *Client) CreateNotificationChannel(ctx context.Context, channel NotificationChannel) (NotificationChannel, error) {
	payload, err := Marshal(notificationChannelWrapper{
		NotificationChannel: channel,
	})
	if err != nil {
		return NotificationChannel{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.GetNotificationChannelsUrl(), payload)
	if err != nil {
		return NotificationChannel{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return NotificationChannel{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[notificationChannelWrapper](response.Body)
	if err != nil {
		return NotificationChannel{}, err
	}

	return wrapper.NotificationChannel, nil
}

func (client *Client) UpdateNotificationChannel(ctx context.Context, channel NotificationChannel) (NotificationChannel, error) {
	payload, err := Marshal(notificationChannelWrapper{
		NotificationChannel: channel,
	})
	if err != nil {
		return NotificationChannel{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.GetNotificationChannelUrl(channel.ID), payload)
	if err != nil {
		return NotificationChannel{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return NotificationChannel{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[notificationChannelWrapper](response.Body)
	if err != nil {
		return NotificationChannel{}, err
	}

	return wrapper.NotificationChannel, nil
}

func (client *Client) DeleteNotificationChannel(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.GetNotificationChannelUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) GetNotificationChannelsUrl() string {
	return fmt.Sprintf(GetNotificationChannels, client.config.url)
}

func (client *Client) GetNotificationChannelUrl(id int) string {
	return fmt.Sprintf(GetNotificationChannel, client.config.url, id)
}
