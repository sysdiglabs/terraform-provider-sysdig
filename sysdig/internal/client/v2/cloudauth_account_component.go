package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	cloudauthAccountComponentsPath = "%s/api/cloudauth/v1/accounts/%s/components"       // POST
	cloudauthAccountComponentPath  = "%s/api/cloudauth/v1/accounts/%s/components/%s/%s" // GET, PUT, DEL
	// getCloudauthAccountPath        = "%s/api/cloudauth/v1/accounts/%s?decrypt=%s" // does GET require decryption?
)

type CloudauthAccountComponentSecureInterface interface {
	Base
	CreateCloudauthAccountComponentSecure(ctx context.Context, accountID string, cloudAccountComponent *CloudauthAccountComponentSecure) (*CloudauthAccountComponentSecure, string, error)
	GetCloudauthAccountComponentSecure(ctx context.Context, accountID, componentType, componentInstance string) (*CloudauthAccountComponentSecure, string, error)
	DeleteCloudauthAccountComponentSecure(ctx context.Context, accountID, componentType, componentInstance string) (string, error)
	UpdateCloudauthAccountComponentSecure(ctx context.Context, accountID, componentType, componentInstance string, cloudAccountComponent *CloudauthAccountComponentSecure) (*CloudauthAccountComponentSecure, string, error)
}

func (client *Client) CreateCloudauthAccountComponentSecure(ctx context.Context, accountID string, cloudAccountComponent *CloudauthAccountComponentSecure) (*CloudauthAccountComponentSecure, string, error) {
	payload, err := client.marshalCloudauthProto(cloudAccountComponent)
	if err != nil {
		return nil, "", err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.cloudauthAccountComponentsURL(accountID), payload)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	cloudauthAccountComponent := &CloudauthAccountComponentSecure{}
	err = client.unmarshalCloudauthProto(response.Body, cloudauthAccountComponent)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccountComponent, "", nil
}

func (client *Client) GetCloudauthAccountComponentSecure(ctx context.Context, accountID, componentType, componentInstance string) (*CloudauthAccountComponentSecure, string, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.cloudauthAccountComponentURL(accountID, componentType, componentInstance), nil)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	cloudauthAccountComponent := &CloudauthAccountComponentSecure{}
	err = client.unmarshalCloudauthProto(response.Body, cloudauthAccountComponent)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccountComponent, "", nil
}

func (client *Client) DeleteCloudauthAccountComponentSecure(ctx context.Context, accountID, componentType, componentInstance string) (string, error) {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.cloudauthAccountComponentURL(accountID, componentType, componentInstance), nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorAndStatusFromResponse(response)
	}
	return "", nil
}

func (client *Client) UpdateCloudauthAccountComponentSecure(ctx context.Context, accountID, componentType, componentInstance string, cloudAccountComponent *CloudauthAccountComponentSecure) (
	*CloudauthAccountComponentSecure, string, error) {
	payload, err := client.marshalCloudauthProto(cloudAccountComponent)
	if err != nil {
		return nil, "", err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.cloudauthAccountComponentURL(accountID, componentType, componentInstance), payload)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	cloudauthAccountComponent := &CloudauthAccountComponentSecure{}
	err = client.unmarshalCloudauthProto(response.Body, cloudauthAccountComponent)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccountComponent, "", nil
}

func (client *Client) cloudauthAccountComponentsURL(accountID string) string {
	return fmt.Sprintf(cloudauthAccountComponentsPath, client.config.url, accountID)
}

func (client *Client) cloudauthAccountComponentURL(accountID string, componentType string, componentInstance string) string {
	return fmt.Sprintf(cloudauthAccountComponentPath, client.config.url, accountID, componentType, componentInstance)
}
