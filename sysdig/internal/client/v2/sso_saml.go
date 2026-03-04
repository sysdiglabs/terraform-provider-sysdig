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

	createSystemSSOSamlPath = "%s/platform/v1/system-sso-settings/"
	getSystemSSOSamlPath    = "%s/platform/v1/system-sso-settings/%d"
	updateSystemSSOSamlPath = "%s/platform/v1/system-sso-settings/%d"
	deleteSystemSSOSamlPath = "%s/platform/v1/system-sso-settings/%d"
)

type SSOSamlInterface interface {
	Base
	CreateSSOSaml(ctx context.Context, isSystem bool, sso *SSOSaml) (*SSOSaml, error)
	GetSSOSaml(ctx context.Context, isSystem bool, id int) (*SSOSaml, error)
	UpdateSSOSaml(ctx context.Context, isSystem bool, id int, sso *SSOSaml) (*SSOSaml, error)
	DeleteSSOSaml(ctx context.Context, isSystem bool, id int) error
}

func (c *Client) CreateSSOSaml(ctx context.Context, isSystem bool, sso *SSOSaml) (result *SSOSaml, err error) {
	payload, err := Marshal(sso)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createSSOSamlURL(isSystem), payload)
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

func (c *Client) GetSSOSaml(ctx context.Context, isSystem bool, id int) (result *SSOSaml, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getSSOSamlURL(isSystem, id), nil)
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

func (c *Client) UpdateSSOSaml(ctx context.Context, isSystem bool, id int, sso *SSOSaml) (result *SSOSaml, err error) {
	payload, err := Marshal(sso)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateSSOSamlURL(isSystem, id), payload)
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

func (c *Client) DeleteSSOSaml(ctx context.Context, isSystem bool, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteSSOSamlURL(isSystem, id), nil)
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

func (c *Client) createSSOSamlURL(isSystem bool) string {
	path := createSSOSamlPath
	if isSystem {
		path = createSystemSSOSamlPath
	}
	return fmt.Sprintf(path, c.config.url)
}

func (c *Client) getSSOSamlURL(isSystem bool, id int) string {
	path := getSSOSamlPath
	if isSystem {
		path = getSystemSSOSamlPath
	}
	return fmt.Sprintf(path, c.config.url, id)
}

func (c *Client) updateSSOSamlURL(isSystem bool, id int) string {
	path := updateSSOSamlPath
	if isSystem {
		path = updateSystemSSOSamlPath
	}
	return fmt.Sprintf(path, c.config.url, id)
}

func (c *Client) deleteSSOSamlURL(isSystem bool, id int) string {
	path := deleteSSOSamlPath
	if isSystem {
		path = deleteSystemSSOSamlPath
	}
	return fmt.Sprintf(path, c.config.url, id)
}
