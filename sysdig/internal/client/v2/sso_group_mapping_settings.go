package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrSSOGroupMappingSettingsNotFound = errors.New("SSO group mapping settings not found")

const (
	ssoGroupMappingSettingsPath = "%s/platform/v1/group-mappings-settings"
)

type SSOGroupMappingSettingsInterface interface {
	Base
	GetSSOGroupMappingSettings(ctx context.Context) (*SSOGroupMappingSettings, error)
	UpdateSSOGroupMappingSettings(ctx context.Context, settings *SSOGroupMappingSettings) (*SSOGroupMappingSettings, error)
}

func (c *Client) GetSSOGroupMappingSettings(ctx context.Context) (result *SSOGroupMappingSettings, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getSSOGroupMappingSettingsURL(), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrSSOGroupMappingSettingsNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*SSOGroupMappingSettings](response.Body)
}

func (c *Client) UpdateSSOGroupMappingSettings(ctx context.Context, settings *SSOGroupMappingSettings) (result *SSOGroupMappingSettings, err error) {
	payload, err := Marshal(settings)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateSSOGroupMappingSettingsURL(), payload)
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

	return Unmarshal[*SSOGroupMappingSettings](response.Body)
}

func (c *Client) getSSOGroupMappingSettingsURL() string {
	return fmt.Sprintf(ssoGroupMappingSettingsPath, c.config.url)
}

func (c *Client) updateSSOGroupMappingSettingsURL() string {
	return fmt.Sprintf(ssoGroupMappingSettingsPath, c.config.url)
}
