package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	cloudAccountsPathV2 = "%s/api/cloudauth/v1/accounts"
	cloudAccountPathV2  = "%s/api/cloudauth/v1/accounts/%s"
)

type CloudAccountSecureInterfaceV2 interface {
	Base
	CreateCloudAccountSecureV2(ctx context.Context, cloudAccount *CloudAccountSecureV2) (*CloudAccountSecureV2, error)
	GetCloudAccountSecureV2(ctx context.Context, accountID string) (*CloudAccountSecureV2, error)
	DeleteCloudAccountSecureV2(ctx context.Context, accountID string) error
	UpdateCloudAccountSecureV2(ctx context.Context, accountID string, cloudAccount *CloudAccountSecureV2) (*CloudAccountSecureV2, error)
}

func (client *Client) CreateCloudAccountSecureV2(ctx context.Context, cloudAccount *CloudAccountSecureV2) (*CloudAccountSecureV2, error) {
	payload, err := Marshal(cloudAccount)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.cloudAccountsV2URL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	return Unmarshal[*CloudAccountSecureV2](response.Body)
}

func (client *Client) GetCloudAccountSecureV2(ctx context.Context, accountID string) (*CloudAccountSecureV2, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.cloudAccountV2URL(accountID), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return Unmarshal[*CloudAccountSecureV2](response.Body)
}

func (client *Client) DeleteCloudAccountSecureV2(ctx context.Context, accountID string) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.cloudAccountV2URL(accountID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorFromResponse(response)
	}
	return nil
}

func (client *Client) UpdateCloudAccountSecureV2(ctx context.Context, accountID string, cloudAccount *CloudAccountSecureV2) (*CloudAccountSecureV2, error) {
	payload, err := Marshal(cloudAccount)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.cloudAccountV2URL(accountID), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	return Unmarshal[*CloudAccountSecureV2](response.Body)
}

func (client *Client) cloudAccountsV2URL() string {
	return fmt.Sprintf(cloudAccountsPathV2, client.config.url)
}

func (client *Client) cloudAccountV2URL(accountID string) string {
	return fmt.Sprintf(cloudAccountPathV2, client.config.url, accountID)
}
