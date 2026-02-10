package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrSSOGlobalSettingsNotFound = errors.New("SSO global settings not found")

const (
	ssoGlobalSettingsPath = "%s/platform/v1/global-sso-settings/%s"
)

type SSOGlobalSettingsInterface interface {
	Base
	GetSSOGlobalSettings(ctx context.Context, product string) (*SSOGlobalSettings, error)
	UpdateSSOGlobalSettings(ctx context.Context, product string, settings *SSOGlobalSettings) (*SSOGlobalSettings, error)
}

func (c *Client) GetSSOGlobalSettings(ctx context.Context, product string) (result *SSOGlobalSettings, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.ssoGlobalSettingsURL(product), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrSSOGlobalSettingsNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*SSOGlobalSettings](response.Body)
}

func (c *Client) UpdateSSOGlobalSettings(ctx context.Context, product string, settings *SSOGlobalSettings) (result *SSOGlobalSettings, err error) {
	payload, err := Marshal(settings)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.ssoGlobalSettingsURL(product), payload)
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

	return Unmarshal[*SSOGlobalSettings](response.Body)
}

func (c *Client) ssoGlobalSettingsURL(product string) string {
	return fmt.Sprintf(ssoGlobalSettingsPath, c.config.url, product)
}
