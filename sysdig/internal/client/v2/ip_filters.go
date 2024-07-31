package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var IpFilterNotFound = errors.New("IP filter not found")

const (
	IPFiltersPath = "%s/platform/v1/ip-filters"
	IPFilterPath  = "%s/platform/v1/ip-filters/%d"
)

type IPFiltersInterface interface {
	Base
	GetIPFilterById(ctx context.Context, id int) (*IPFilter, error)
	CreateIPFilter(ctx context.Context, ipFilter *IPFilter) (*IPFilter, error)
	UpdateIPFilter(ctx context.Context, ipFilter *IPFilter, id int) (*IPFilter, error)
	DeleteIPFilter(ctx context.Context, id int) error
}

func (client *Client) GetIPFilterById(ctx context.Context, id int) (*IPFilter, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetIPFilterURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	ipFilter, err := Unmarshal[IPFilter](response.Body)
	if err != nil {
		return nil, err
	}

	return &ipFilter, nil
}

func (client *Client) CreateIPFilter(ctx context.Context, ipFilter *IPFilter) (*IPFilter, error) {
	payload, err := Marshal(ipFilter)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.GetIPFiltersURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, client.ErrorFromResponse(response)
	}

	created, err := Unmarshal[IPFilter](response.Body)
	if err != nil {
		return nil, err
	}

	return &created, nil

}

func (client *Client) UpdateIPFilter(ctx context.Context, ipFilter *IPFilter, id int) (*IPFilter, error) {
	payload, err := Marshal(ipFilter)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.GetIPFilterURL(id), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	updated, err := Unmarshal[IPFilter](response.Body)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (client *Client) DeleteIPFilter(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.GetIPFilterURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) GetIPFilterURL(id int) string {
	return fmt.Sprintf(IPFilterPath, client.config.url, id)
}

func (client *Client) GetIPFiltersURL() string {
	return fmt.Sprintf(IPFiltersPath, client.config.url)
}
