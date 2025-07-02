package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrIPFilterNotFound = errors.New("IP filter not found")

const (
	ipFiltersPath = "%s/platform/v1/ip-filters"
	ipFilterPath  = "%s/platform/v1/ip-filters/%d"
)

type IPFiltersInterface interface {
	Base
	GetIPFilterByID(ctx context.Context, id int) (*IPFilter, error)
	CreateIPFilter(ctx context.Context, ipFilter *IPFilter) (*IPFilter, error)
	UpdateIPFilter(ctx context.Context, ipFilter *IPFilter, id int) (*IPFilter, error)
	DeleteIPFilter(ctx context.Context, id int) error
}

func (c *Client) GetIPFilterByID(ctx context.Context, id int) (ipFilter *IPFilter, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getIPFilterURL(id), nil)
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

	return Unmarshal[*IPFilter](response.Body)
}

func (c *Client) CreateIPFilter(ctx context.Context, ipFilter *IPFilter) (createdFilter *IPFilter, err error) {
	payload, err := Marshal(ipFilter)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.getIPFiltersURL(), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusCreated {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*IPFilter](response.Body)
}

func (c *Client) UpdateIPFilter(ctx context.Context, ipFilter *IPFilter, id int) (updatedFilter *IPFilter, err error) {
	payload, err := Marshal(ipFilter)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.getIPFilterURL(id), payload)
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

	return Unmarshal[*IPFilter](response.Body)
}

func (c *Client) DeleteIPFilter(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.getIPFilterURL(id), nil)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusNotFound {
		return c.ErrorFromResponse(response)
	}

	return nil
}

func (c *Client) getIPFilterURL(id int) string {
	return fmt.Sprintf(ipFilterPath, c.config.url, id)
}

func (c *Client) getIPFiltersURL() string {
	return fmt.Sprintf(ipFiltersPath, c.config.url)
}
