package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	ipFiltersSettingsPath = "%s/platform/v1/ip-filters-settings"
)

type IPFilteringSettingsInterface interface {
	Base
	GetIPFilteringSettings(ctx context.Context) (*IPFiltersSettings, error)
	UpdateIPFilteringSettings(ctx context.Context, ipFiltersSettings *IPFiltersSettings) (*IPFiltersSettings, error)
}

func (c *Client) GetIPFilteringSettings(ctx context.Context) (settings *IPFiltersSettings, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.GetIPFiltersSettingsURL(), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		err = c.ErrorFromResponse(response)
		return nil, err
	}

	return Unmarshal[*IPFiltersSettings](response.Body)
}

func (c *Client) UpdateIPFilteringSettings(ctx context.Context, ipFiltersSettings *IPFiltersSettings) (settings *IPFiltersSettings, err error) {
	payload, err := Marshal(ipFiltersSettings)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.GetIPFiltersSettingsURL(), payload)
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

	return Unmarshal[*IPFiltersSettings](response.Body)
}

func (c *Client) GetIPFiltersSettingsURL() string {
	return fmt.Sprintf(ipFiltersSettingsPath, c.config.url)
}
