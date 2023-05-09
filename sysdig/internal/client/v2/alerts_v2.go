package v2

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type AlertsV2Interface interface {
	Base
}

type AlertV2Type string

const (
	alertsV2Path = "%s/api/v2/alerts"
	alertV2Path  = "%s/api/v2/alerts/%d"

	Prometheus AlertV2Type = "PROMETHEUS"
)

type AlertV2PrometheusInterface interface {
	CreateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error)
	UpdateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error)
	GetAlertV2Prometheus(ctx context.Context, alertID int) (AlertV2Prometheus, error)
	DeleteAlertV2Prometheus(ctx context.Context, alertID int) error
}

func (client *Client) CreateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	payload, err := Marshal(alertV2PrometheusWrapper{Alert: alert})
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	body, err := client.createAlertV2(ctx, payload)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	wrapper, err := Unmarshal[alertV2PrometheusWrapper](body)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) UpdateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	payload, err := Marshal(alertV2PrometheusWrapper{Alert: alert})
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	body, err := client.updateAlertV2(ctx, alert.ID, payload)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	wrapper, err := Unmarshal[alertV2PrometheusWrapper](body)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) GetAlertV2Prometheus(ctx context.Context, alertID int) (AlertV2Prometheus, error) {
	body, err := client.getAlertV2(ctx, alertID)
	if err != nil {
		return AlertV2Prometheus{}, err
	}
	wrapper, err := Unmarshal[alertV2PrometheusWrapper](body)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) DeleteAlertV2Prometheus(ctx context.Context, alertID int) error {
	return client.deleteAlertV2(ctx, alertID)
}

func (client *Client) createAlertV2(ctx context.Context, alertJson io.Reader) (io.ReadCloser, error) {
	response, err := client.requester.Request(ctx, http.MethodPost, client.alertsV2URL(), alertJson)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return response.Body, nil
}

func (client *Client) updateAlertV2(ctx context.Context, alertID int, alertJson io.Reader) (io.ReadCloser, error) {
	response, err := client.requester.Request(ctx, http.MethodPut, client.alertV2URL(alertID), alertJson)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return response.Body, nil
}

func (client *Client) deleteAlertV2(ctx context.Context, alertID int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.alertV2URL(alertID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) getAlertV2(ctx context.Context, alertID int) (io.ReadCloser, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.alertV2URL(alertID), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return response.Body, nil
}

func (client *Client) addNotificationChannelType(ctx context.Context, notificationChannelConfigList []NotificationChannelConfigV2) error {
	// on put/posts the api wants the type of the channel even if it can be inferred
	for i, n := range notificationChannelConfigList {
		nc, err := client.GetNotificationChannelById(ctx, n.ChannelID)
		if err != nil {
			return fmt.Errorf("error getting info for notification channel %d: %w", n.ChannelID, err)
		}
		notificationChannelConfigList[i].Type = nc.Type
	}
	return nil
}

func (client *Client) alertsV2URL() string {
	return fmt.Sprintf(alertsV2Path, client.config.url)
}

func (client *Client) alertV2URL(alertID int) string {
	return fmt.Sprintf(alertV2Path, client.config.url, alertID)
}
