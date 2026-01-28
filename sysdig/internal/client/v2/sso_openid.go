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
)

type SSOOpenIDInterface interface {
	Base
	CreateSSOOpenID(ctx context.Context, sso *SSOOpenID) (*SSOOpenID, error)
	GetSSOOpenID(ctx context.Context, id int) (*SSOOpenID, error)
	UpdateSSOOpenID(ctx context.Context, id int, sso *SSOOpenID) (*SSOOpenID, error)
	DeleteSSOOpenID(ctx context.Context, id int) error
}

func (c *Client) CreateSSOOpenID(ctx context.Context, sso *SSOOpenID) (result *SSOOpenID, err error) {
	payload, err := Marshal(sso)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createSSOOpenIDURL(), payload)
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

func (c *Client) GetSSOOpenID(ctx context.Context, id int) (result *SSOOpenID, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getSSOOpenIDURL(id), nil)
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

func (c *Client) UpdateSSOOpenID(ctx context.Context, id int, sso *SSOOpenID) (result *SSOOpenID, err error) {
	payload, err := Marshal(sso)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateSSOOpenIDURL(id), payload)
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

func (c *Client) DeleteSSOOpenID(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteSSOOpenIDURL(id), nil)
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

func (c *Client) createSSOOpenIDURL() string {
	return fmt.Sprintf(createSSOOpenIDPath, c.config.url)
}

func (c *Client) getSSOOpenIDURL(id int) string {
	return fmt.Sprintf(getSSOOpenIDPath, c.config.url, id)
}

func (c *Client) updateSSOOpenIDURL(id int) string {
	return fmt.Sprintf(updateSSOOpenIDPath, c.config.url, id)
}

func (c *Client) deleteSSOOpenIDURL(id int) string {
	return fmt.Sprintf(deleteSSOOpenIDPath, c.config.url, id)
}
