package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrSSOSamlNotFound = errors.New("SSO SAML configuration not found")

const (
	createSSOSamlPath = "%s/platform/v1/sso-settings/"
	getSSOSamlPath    = "%s/platform/v1/sso-settings/%d"
	updateSSOSamlPath = "%s/platform/v1/sso-settings/%d"
	deleteSSOSamlPath = "%s/platform/v1/sso-settings/%d"
)

type SSOSamlInterface interface {
	Base
	CreateSSOSaml(ctx context.Context, sso *SSOSaml) (*SSOSaml, error)
	GetSSOSaml(ctx context.Context, id int) (*SSOSaml, error)
	UpdateSSOSaml(ctx context.Context, id int, sso *SSOSaml) (*SSOSaml, error)
	DeleteSSOSaml(ctx context.Context, id int) error
}

func (c *Client) CreateSSOSaml(ctx context.Context, sso *SSOSaml) (result *SSOSaml, err error) {
	payload, err := Marshal(sso)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createSSOSamlURL(), payload)
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

	return Unmarshal[*SSOSaml](response.Body)
}

func (c *Client) GetSSOSaml(ctx context.Context, id int) (result *SSOSaml, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getSSOSamlURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrSSOSamlNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*SSOSaml](response.Body)
}

func (c *Client) UpdateSSOSaml(ctx context.Context, id int, sso *SSOSaml) (result *SSOSaml, err error) {
	payload, err := Marshal(sso)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateSSOSamlURL(id), payload)
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

	return Unmarshal[*SSOSaml](response.Body)
}

func (c *Client) DeleteSSOSaml(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteSSOSamlURL(id), nil)
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

func (c *Client) createSSOSamlURL() string {
	return fmt.Sprintf(createSSOSamlPath, c.config.url)
}

func (c *Client) getSSOSamlURL(id int) string {
	return fmt.Sprintf(getSSOSamlPath, c.config.url, id)
}

func (c *Client) updateSSOSamlURL(id int) string {
	return fmt.Sprintf(updateSSOSamlPath, c.config.url, id)
}

func (c *Client) deleteSSOSamlURL(id int) string {
	return fmt.Sprintf(deleteSSOSamlPath, c.config.url, id)
}
