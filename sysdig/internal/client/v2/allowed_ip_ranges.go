package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var AllowedIpRangeNotFound = errors.New("IP range not found")

const (
	AllowedIPRangesPath = "%s/platform/v1/allowed-ip-ranges"
	AllowedIPRangePath  = "%s/platform/v1/allowed-ip-ranges/%d"
)

type AllowedIPRangesInterface interface {
	Base
	GetAllowedIpRangeById(ctx context.Context, id int) (*AllowedIpRange, error)
	CreateAllowedIpRange(ctx context.Context, allowedIpRange *AllowedIpRange) (*AllowedIpRange, error)
	UpdateAllowedIpRange(ctx context.Context, allowedIpRange *AllowedIpRange, id int) (*AllowedIpRange, error)
	DeleteAllowedIpRange(ctx context.Context, id int) error
}

func (client *Client) GetAllowedIpRangeById(ctx context.Context, id int) (*AllowedIpRange, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetAllowedIpRangeURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	allowedIpRange, err := Unmarshal[AllowedIpRange](response.Body)
	if err != nil {
		return nil, err
	}

	return &allowedIpRange, nil
}

func (client *Client) CreateAllowedIpRange(ctx context.Context, allowedIpRange *AllowedIpRange) (*AllowedIpRange, error) {
	payload, err := Marshal(allowedIpRange)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.GetAllowedIpRangesURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, client.ErrorFromResponse(response)
	}

	created, err := Unmarshal[AllowedIpRange](response.Body)
	if err != nil {
		return nil, err
	}

	return &created, nil

}

func (client *Client) UpdateAllowedIpRange(ctx context.Context, allowedIpRange *AllowedIpRange, id int) (*AllowedIpRange, error) {
	payload, err := Marshal(allowedIpRange)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.GetAllowedIpRangeURL(id), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	updated, err := Unmarshal[AllowedIpRange](response.Body)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (client *Client) DeleteAllowedIpRange(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.GetAllowedIpRangeURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) GetAllowedIpRangeURL(id int) string {
	return fmt.Sprintf(AllowedIPRangePath, client.config.url, id)
}

func (client *Client) GetAllowedIpRangesURL() string {
	return fmt.Sprintf(AllowedIPRangesPath, client.config.url)
}
