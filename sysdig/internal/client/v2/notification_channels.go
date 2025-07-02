package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	getNotificationChannels = "%s/api/notificationChannels"
	getNotificationChannel  = "%s/api/notificationChannels/%d"
)

var ErrNotificationChannelNotFound = errors.New("notification channel not found")

type NotificationChannelInterface interface {
	Base
	GetNotificationChannelByID(ctx context.Context, id int) (NotificationChannel, error)
	GetNotificationChannelByName(ctx context.Context, name string) (NotificationChannel, error)
	CreateNotificationChannel(ctx context.Context, channel NotificationChannel) (NotificationChannel, error)
	UpdateNotificationChannel(ctx context.Context, channel NotificationChannel) (NotificationChannel, error)
	DeleteNotificationChannel(ctx context.Context, id int) error
}

func (c *Client) GetNotificationChannelByID(ctx context.Context, id int) (nc NotificationChannel, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getNotificationChannelURL(id), nil)
	if err != nil {
		return NotificationChannel{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode == http.StatusNotFound {
		return NotificationChannel{}, ErrNotificationChannelNotFound
	}
	if response.StatusCode != http.StatusOK {
		return NotificationChannel{}, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[notificationChannelWrapper](response.Body)
	if err != nil {
		return NotificationChannel{}, err
	}

	if wrapper.NotificationChannel.ID == 0 {
		return NotificationChannel{}, fmt.Errorf("notificationChannel with ID: %d does not exists", id)
	}

	return wrapper.NotificationChannel, nil
}

func (c *Client) GetNotificationChannelByName(ctx context.Context, name string) (nc NotificationChannel, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getNotificationChannelsURL(), nil)
	if err != nil {
		return NotificationChannel{}, nil
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return NotificationChannel{}, c.ErrorFromResponse(response)
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

func (c *Client) CreateNotificationChannel(ctx context.Context, channel NotificationChannel) (nc NotificationChannel, err error) {
	payload, err := Marshal(notificationChannelWrapper{
		NotificationChannel: channel,
	})
	if err != nil {
		return NotificationChannel{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.getNotificationChannelsURL(), payload)
	if err != nil {
		return NotificationChannel{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return NotificationChannel{}, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[notificationChannelWrapper](response.Body)
	if err != nil {
		return NotificationChannel{}, err
	}

	return wrapper.NotificationChannel, nil
}

func (c *Client) UpdateNotificationChannel(ctx context.Context, channel NotificationChannel) (nc NotificationChannel, err error) {
	payload, err := Marshal(notificationChannelWrapper{
		NotificationChannel: channel,
	})
	if err != nil {
		return NotificationChannel{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.getNotificationChannelURL(channel.ID), payload)
	if err != nil {
		return NotificationChannel{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return NotificationChannel{}, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[notificationChannelWrapper](response.Body)
	if err != nil {
		return NotificationChannel{}, err
	}

	return wrapper.NotificationChannel, nil
}

func (c *Client) DeleteNotificationChannel(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.getNotificationChannelURL(id), nil)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return c.ErrorFromResponse(response)
	}

	return nil
}

func (c *Client) getNotificationChannelsURL() string {
	return fmt.Sprintf(getNotificationChannels, c.config.url)
}

func (c *Client) getNotificationChannelURL(id int) string {
	return fmt.Sprintf(getNotificationChannel, c.config.url, id)
}
