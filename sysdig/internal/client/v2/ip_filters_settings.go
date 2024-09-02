package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	IPFiltersSettingsPath = "%s/platform/v1/ip-filters-settings"
)

type IPFiltersSettingsInterface interface {
	Base
	GetIPFiltersSettings(ctx context.Context) (*IPFiltersSettings, error)
	UpdateIPFiltersSettings(ctx context.Context, ipFiltersSettings *IPFiltersSettings) (*IPFiltersSettings, error)
}

func (client *Client) GetIPFiltersSettings(ctx context.Context) (*IPFiltersSettings, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetIPFiltersSettingsURL(), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	ipFiltersSettings, err := Unmarshal[IPFiltersSettings](response.Body)
	if err != nil {
		return nil, err
	}

	return &ipFiltersSettings, nil
}

func (client *Client) UpdateIPFiltersSettings(ctx context.Context, ipFiltersSettings *IPFiltersSettings) (*IPFiltersSettings, error) {
	payload, err := Marshal(ipFiltersSettings)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.GetIPFiltersSettingsURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	updated, err := Unmarshal[IPFiltersSettings](response.Body)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (client *Client) GetIPFiltersSettingsURL() string {
	return fmt.Sprintf(IPFiltersSettingsPath, client.config.url)
}
