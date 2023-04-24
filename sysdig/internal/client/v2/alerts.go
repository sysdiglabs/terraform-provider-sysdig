package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	CreateAlertPath = "%s/api/alerts"
	DeleteAlertPath = "%s/api/alerts/%d"
	UpdateAlertPath = "%s/api/alerts/%d"
	GetAlertPath    = "%s/api/alerts/%d"
)

type AlertInterface interface {
	Base
	CreateAlert(ctx context.Context, alert Alert) (Alert, error)
	DeleteAlert(ctx context.Context, alertID int) error
	UpdateAlert(ctx context.Context, alert Alert) (Alert, error)
	GetAlertByID(ctx context.Context, alertID int) (Alert, error)
}

func (client *Client) CreateAlert(ctx context.Context, alert Alert) (Alert, error) {
	payload, err := Marshal[alertWrapper](alertWrapper{Alert: alert})
	if err != nil {
		return Alert{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.createAlertURL(), payload)
	if err != nil {
		return Alert{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Alert{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[alertWrapper](response.Body)
	if err != nil {
		return Alert{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) DeleteAlert(ctx context.Context, alertID int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.deleteAlertURL(alertID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return err
}

func (client *Client) UpdateAlert(ctx context.Context, alert Alert) (Alert, error) {
	payload, err := Marshal[alertWrapper](alertWrapper{Alert: alert})
	if err != nil {
		return Alert{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.updateAlertURL(alert.ID), payload)
	if err != nil {
		return Alert{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Alert{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[alertWrapper](response.Body)
	if err != nil {
		return Alert{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) GetAlertByID(ctx context.Context, alertID int) (Alert, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetAlertByIDURL(alertID), nil)
	if err != nil {
		return Alert{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Alert{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[alertWrapper](response.Body)
	if err != nil {
		return Alert{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) createAlertURL() string {
	return fmt.Sprintf(CreateAlertPath, client.config.url)
}

func (client *Client) deleteAlertURL(alertID int) string {
	return fmt.Sprintf(DeleteAlertPath, client.config.url, alertID)
}

func (client *Client) updateAlertURL(alertID int) string {
	return fmt.Sprintf(UpdateAlertPath, client.config.url, alertID)
}

func (client *Client) GetAlertByIDURL(alertID int) string {
	return fmt.Sprintf(GetAlertPath, client.config.url, alertID)
}
