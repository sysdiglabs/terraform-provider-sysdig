package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	cloudAccountsPath               = "%s/api/cloud/v2/accounts"
	cloudAccountsWithExternalIDPath = "%s/api/cloud/v2/accounts?includeExternalID=true&upsert=true"
	cloudAccountPath                = "%s/api/cloud/v2/accounts/%s"
	cloudAccountWithExternalIDPath  = "%s/api/cloud/v2/accounts/%s?includeExternalID=true"
	providersPath                   = "%v/api/v2/providers"
)

type CloudAccountSecureInterface interface {
	Base
	CreateCloudAccountSecure(ctx context.Context, cloudAccount *CloudAccountSecure) (*CloudAccountSecure, error)
	GetCloudAccountSecure(ctx context.Context, accountID string) (*CloudAccountSecure, error)
	DeleteCloudAccountSecure(ctx context.Context, accountID string) error
	UpdateCloudAccountSecure(ctx context.Context, accountID string, cloudAccount *CloudAccountSecure) (*CloudAccountSecure, error)
}

type CloudAccountMonitorInterface interface {
	Base
	CreateCloudAccountMonitor(ctx context.Context, provider *CloudAccountMonitor) (*CloudAccountMonitor, error)
	UpdateCloudAccountMonitor(ctx context.Context, id int, provider *CloudAccountMonitor) (*CloudAccountMonitor, error)
	GetCloudAccountMonitor(ctx context.Context, id int) (*CloudAccountMonitor, error)
	DeleteCloudAccountMonitor(ctx context.Context, id int) error
}

func (client *Client) CreateCloudAccountSecure(ctx context.Context, cloudAccount *CloudAccountSecure) (*CloudAccountSecure, error) {
	payload, err := Marshal(cloudAccount)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.cloudAccountsURL(true), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	return Unmarshal[*CloudAccountSecure](response.Body)
}

func (client *Client) GetCloudAccountSecure(ctx context.Context, accountID string) (*CloudAccountSecure, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.cloudAccountURL(accountID, true), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return Unmarshal[*CloudAccountSecure](response.Body)
}

func (client *Client) DeleteCloudAccountSecure(ctx context.Context, accountID string) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.cloudAccountURL(accountID, false), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorFromResponse(response)
	}
	return nil
}

func (client *Client) UpdateCloudAccountSecure(ctx context.Context, accountID string, cloudAccount *CloudAccountSecure) (*CloudAccountSecure, error) {
	payload, err := Marshal(cloudAccount)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.cloudAccountURL(accountID, true), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	return Unmarshal[*CloudAccountSecure](response.Body)
}

func (client *Client) cloudAccountsURL(includeExternalID bool) string {
	if includeExternalID {
		return fmt.Sprintf(cloudAccountsWithExternalIDPath, client.config.url)
	}
	return fmt.Sprintf(cloudAccountsPath, client.config.url)
}

func (client *Client) cloudAccountURL(accountID string, includeExternalID bool) string {
	if includeExternalID {
		return fmt.Sprintf(cloudAccountWithExternalIDPath, client.config.url, accountID)
	}
	return fmt.Sprintf(cloudAccountPath, client.config.url, accountID)
}

func (client *Client) CreateCloudAccountMonitor(ctx context.Context, provider *CloudAccountMonitor) (*CloudAccountMonitor, error) {
	payload, err := Marshal(provider)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.getProvidersURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[cloudAccountWrapperMonitor](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.CloudAccount, nil
}

func (client *Client) UpdateCloudAccountMonitor(ctx context.Context, id int, provider *CloudAccountMonitor) (*CloudAccountMonitor, error) {
	payload, err := Marshal(provider)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.getProviderURL(id), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[cloudAccountWrapperMonitor](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.CloudAccount, nil
}

func (client *Client) GetCloudAccountMonitor(ctx context.Context, id int) (*CloudAccountMonitor, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.getProviderURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[cloudAccountWrapperMonitor](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.CloudAccount, nil
}

func (client *Client) DeleteCloudAccountMonitor(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.getProviderURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) getProviderURL(id int) string {
	return fmt.Sprintf("%v/%v", client.getProvidersURL(), id)
}

func (client *Client) getProvidersURL() string {
	return fmt.Sprintf(providersPath, client.config.url)
}
