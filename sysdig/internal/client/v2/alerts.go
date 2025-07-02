package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	createAlertPath = "%s/api/alerts"
	deleteAlertPath = "%s/api/alerts/%d"
	updateAlertPath = "%s/api/alerts/%d"
	getAlertPath    = "%s/api/alerts/%d"
)

type AlertInterface interface {
	Base
	CreateAlert(ctx context.Context, alert Alert) (Alert, error)
	GetAlertByID(ctx context.Context, alertID int) (Alert, error)
	UpdateAlert(ctx context.Context, alert Alert) (Alert, error)
	DeleteAlertByID(ctx context.Context, alertID int) error
}

func (c *Client) CreateAlert(ctx context.Context, alert Alert) (createdAlert Alert, err error) {
	payload, err := Marshal(alertWrapper{Alert: alert})
	if err != nil {
		return Alert{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createAlertURL(), payload)
	if err != nil {
		return Alert{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Alert{}, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[alertWrapper](response.Body)
	if err != nil {
		return Alert{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) DeleteAlertByID(ctx context.Context, alertID int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteAlertURL(alertID), nil)
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

	return err
}

func (c *Client) UpdateAlert(ctx context.Context, alert Alert) (updatedAlert Alert, err error) {
	payload, err := Marshal(alertWrapper{Alert: alert})
	if err != nil {
		return Alert{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateAlertURL(alert.ID), payload)
	if err != nil {
		return Alert{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Alert{}, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[alertWrapper](response.Body)
	if err != nil {
		return Alert{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) GetAlertByID(ctx context.Context, alertID int) (alert Alert, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getAlertByIDURL(alertID), nil)
	if err != nil {
		return Alert{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Alert{}, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[alertWrapper](response.Body)
	if err != nil {
		return Alert{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) createAlertURL() string {
	return fmt.Sprintf(createAlertPath, c.config.url)
}

func (c *Client) deleteAlertURL(alertID int) string {
	return fmt.Sprintf(deleteAlertPath, c.config.url, alertID)
}

func (c *Client) updateAlertURL(alertID int) string {
	return fmt.Sprintf(updateAlertPath, c.config.url, alertID)
}

func (c *Client) getAlertByIDURL(alertID int) string {
	return fmt.Sprintf(getAlertPath, c.config.url, alertID)
}
