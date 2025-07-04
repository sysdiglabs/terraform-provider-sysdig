package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	cloudauthAccountComponentsPath = "%s/api/cloudauth/v1/accounts/%s/components"
	cloudauthAccountComponentPath  = "%s/api/cloudauth/v1/accounts/%s/components/%s/%s"
)

type CloudauthAccountComponentSecureInterface interface {
	Base
	CreateCloudauthAccountComponentSecure(ctx context.Context, accountID string, cloudAccountComponent *CloudauthAccountComponentSecure) (*CloudauthAccountComponentSecure, string, error)
	GetCloudauthAccountComponentSecure(ctx context.Context, accountID, componentType, componentInstance string) (*CloudauthAccountComponentSecure, string, error)
	DeleteCloudauthAccountComponentSecure(ctx context.Context, accountID, componentType, componentInstance string) (string, error)
	UpdateCloudauthAccountComponentSecure(ctx context.Context, accountID, componentType, componentInstance string, cloudAccountComponent *CloudauthAccountComponentSecure) (*CloudauthAccountComponentSecure, string, error)
}

func (c *Client) CreateCloudauthAccountComponentSecure(ctx context.Context, accountID string, cloudAccountComponent *CloudauthAccountComponentSecure) (component *CloudauthAccountComponentSecure, errString string, err error) {
	payload, err := c.marshalCloudauthProto(cloudAccountComponent)
	if err != nil {
		return nil, "", err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.cloudauthAccountComponentsURL(accountID), payload)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	cloudauthAccountComponent := &CloudauthAccountComponentSecure{}
	err = c.unmarshalCloudauthProto(response.Body, cloudauthAccountComponent)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccountComponent, "", nil
}

func (c *Client) GetCloudauthAccountComponentSecure(ctx context.Context, accountID, componentType, componentInstance string) (component *CloudauthAccountComponentSecure, errString string, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.cloudauthAccountComponentURL(accountID, componentType, componentInstance), nil)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	cloudauthAccountComponent := &CloudauthAccountComponentSecure{}
	err = c.unmarshalCloudauthProto(response.Body, cloudauthAccountComponent)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccountComponent, "", nil
}

func (c *Client) DeleteCloudauthAccountComponentSecure(ctx context.Context, accountID, componentType, componentInstance string) (errString string, err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.cloudauthAccountComponentURL(accountID, componentType, componentInstance), nil)
	if err != nil {
		return "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return c.ErrorAndStatusFromResponse(response)
	}
	return "", nil
}

func (c *Client) UpdateCloudauthAccountComponentSecure(ctx context.Context, accountID, componentType, componentInstance string, cloudAccountComponent *CloudauthAccountComponentSecure) (component *CloudauthAccountComponentSecure, errString string, err error) {
	payload, err := c.marshalCloudauthProto(cloudAccountComponent)
	if err != nil {
		return nil, "", err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.cloudauthAccountComponentURL(accountID, componentType, componentInstance), payload)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	cloudauthAccountComponent := &CloudauthAccountComponentSecure{}
	err = c.unmarshalCloudauthProto(response.Body, cloudauthAccountComponent)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccountComponent, "", nil
}

func (c *Client) cloudauthAccountComponentsURL(accountID string) string {
	return fmt.Sprintf(cloudauthAccountComponentsPath, c.config.url, accountID)
}

func (c *Client) cloudauthAccountComponentURL(accountID string, componentType string, componentInstance string) string {
	return fmt.Sprintf(cloudauthAccountComponentPath, c.config.url, accountID, componentType, componentInstance)
}
