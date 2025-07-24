package v2

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

const (
	createAutomationPath = "%s/api/platform-automations/v1/automations"
	getAutomationPath    = "%s/api/platform-automations/v1/automations/%s"
	updateAutomationPath = "%s/api/platform-automations/v1/automations/%s"
	deleteAutomationPath = "%s/api/platform-automations/v1/automations/%s"
)

type AutomationInterface interface {
	Base
	CreateAutomation(ctx context.Context, automationJSON []byte) (*AutomationResponse, error)
	GetAutomationByID(ctx context.Context, id string) (*AutomationResponse, error)
	UpdateAutomation(ctx context.Context, id string, automationJSON []byte) (*AutomationResponse, error)
	DeleteAutomation(ctx context.Context, id string) error
}

func (c *Client) CreateAutomation(ctx context.Context, automationJSON []byte) (automation *AutomationResponse, err error) {
	response, err := c.requester.Request(ctx, http.MethodPost, c.createAutomationURL(), bytes.NewReader(automationJSON))
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*AutomationResponse](response.Body)
}

func (c *Client) GetAutomationByID(ctx context.Context, id string) (automation *AutomationResponse, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getAutomationURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*AutomationResponse](response.Body)
}

func (c *Client) UpdateAutomation(ctx context.Context, id string, automationJSON []byte) (automation *AutomationResponse, err error) {
	response, err := c.requester.Request(ctx, http.MethodPut, c.updateAutomationURL(id), bytes.NewReader(automationJSON))
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*AutomationResponse](response.Body)
}

func (c *Client) DeleteAutomation(ctx context.Context, id string) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteAutomationURL(id), nil)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		return c.ErrorFromResponse(response)
	}

	return nil
}

// URL builders
func (c *Client) createAutomationURL() string {
	return fmt.Sprintf(createAutomationPath, c.config.url)
}

func (c *Client) getAutomationURL(id string) string {
	return fmt.Sprintf(getAutomationPath, c.config.url, id)
}

func (c *Client) updateAutomationURL(id string) string {
	return fmt.Sprintf(updateAutomationPath, c.config.url, id)
}

func (c *Client) deleteAutomationURL(id string) string {
	return fmt.Sprintf(deleteAutomationPath, c.config.url, id)
}
