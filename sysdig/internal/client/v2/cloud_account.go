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
	trustedCloudIdentityPath        = "%s/api/cloud/v2/%s/trustedIdentity"
)

type CloudAccountSecureInterface interface {
	Base
	CreateCloudAccountSecure(ctx context.Context, cloudAccount *CloudAccountSecure) (*CloudAccountSecure, error)
	GetCloudAccountSecure(ctx context.Context, accountID string) (*CloudAccountSecure, error)
	DeleteCloudAccountSecure(ctx context.Context, accountID string) error
	UpdateCloudAccountSecure(ctx context.Context, accountID string, cloudAccount *CloudAccountSecure) (*CloudAccountSecure, error)
	GetTrustedCloudIdentitySecure(ctx context.Context, provider string) (string, error)
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

func (client *Client) GetTrustedCloudIdentitySecure(ctx context.Context, provider string) (string, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.trustedCloudIdentityURL(provider), nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", client.ErrorFromResponse(response)
	}

	return Unmarshal[string](response.Body)
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

func (client *Client) trustedCloudIdentityURL(provider string) string {
	return fmt.Sprintf(trustedCloudIdentityPath, client.config.url, provider)
}
