package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrSSOOpenIDNotFound = errors.New("SSO OpenID configuration not found")

const (
	createSSOOpenIDPath = "%s/platform/v1/sso-settings/"
	getSSOOpenIDPath    = "%s/platform/v1/sso-settings/%d"
	updateSSOOpenIDPath = "%s/platform/v1/sso-settings/%d"
	deleteSSOOpenIDPath = "%s/platform/v1/sso-settings/%d"

	createSystemSSOOpenIDPath = "%s/platform/v1/system-sso-settings/"
	getSystemSSOOpenIDPath    = "%s/platform/v1/system-sso-settings/%d"
	updateSystemSSOOpenIDPath = "%s/platform/v1/system-sso-settings/%d"
	deleteSystemSSOOpenIDPath = "%s/platform/v1/system-sso-settings/%d"
)

type SSOOpenIDInterface interface {
	Base
	CreateSSOOpenID(ctx context.Context, isSystem bool, sso *SSOOpenID) (*SSOOpenID, error)
	GetSSOOpenID(ctx context.Context, isSystem bool, id int) (*SSOOpenID, error)
	UpdateSSOOpenID(ctx context.Context, isSystem bool, id int, sso *SSOOpenID) (*SSOOpenID, error)
	DeleteSSOOpenID(ctx context.Context, isSystem bool, id int) error
}

func (c *Client) CreateSSOOpenID(ctx context.Context, isSystem bool, sso *SSOOpenID) (result *SSOOpenID, err error) {
	payload, err := Marshal(sso)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createSSOOpenIDURL(isSystem), payload)
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

	return Unmarshal[*SSOOpenID](response.Body)
}

func (c *Client) GetSSOOpenID(ctx context.Context, isSystem bool, id int) (result *SSOOpenID, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getSSOOpenIDURL(isSystem, id), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrSSOOpenIDNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*SSOOpenID](response.Body)
}

func (c *Client) UpdateSSOOpenID(ctx context.Context, isSystem bool, id int, sso *SSOOpenID) (result *SSOOpenID, err error) {
	payload, err := Marshal(sso)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateSSOOpenIDURL(isSystem, id), payload)
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

	return Unmarshal[*SSOOpenID](response.Body)
}

func (c *Client) DeleteSSOOpenID(ctx context.Context, isSystem bool, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteSSOOpenIDURL(isSystem, id), nil)
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

func (c *Client) createSSOOpenIDURL(isSystem bool) string {
	path := createSSOOpenIDPath
	if isSystem {
		path = createSystemSSOOpenIDPath
	}
	return fmt.Sprintf(path, c.config.url)
}

func (c *Client) getSSOOpenIDURL(isSystem bool, id int) string {
	path := getSSOOpenIDPath
	if isSystem {
		path = getSystemSSOOpenIDPath
	}
	return fmt.Sprintf(path, c.config.url, id)
}

func (c *Client) updateSSOOpenIDURL(isSystem bool, id int) string {
	path := updateSSOOpenIDPath
	if isSystem {
		path = updateSystemSSOOpenIDPath
	}
	return fmt.Sprintf(path, c.config.url, id)
}

func (c *Client) deleteSSOOpenIDURL(isSystem bool, id int) string {
	path := deleteSSOOpenIDPath
	if isSystem {
		path = deleteSystemSSOOpenIDPath
	}
	return fmt.Sprintf(path, c.config.url, id)
}
