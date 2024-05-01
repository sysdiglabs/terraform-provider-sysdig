package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	cloudauthAccountFeaturePath = "%s/api/cloudauth/v1/accounts/%s/feature/%s" // GET, PUT, DEL
)

type CloudauthAccountFeatureSecureInterface interface {
	Base
	CreateOrUpdateCloudauthAccountFeatureSecure(ctx context.Context, accountID, featureType string, cloudAccountFeature *CloudauthAccountFeatureSecure) (*CloudauthAccountFeatureSecure, string, error)
	GetCloudauthAccountFeatureSecure(ctx context.Context, accountID, featureType string) (*CloudauthAccountFeatureSecure, string, error)
	DeleteCloudauthAccountFeatureSecure(ctx context.Context, accountID, featureType string) (string, error)
}

// both create and update makes a PUT call to backend
func (client *Client) CreateOrUpdateCloudauthAccountFeatureSecure(ctx context.Context, accountID, featureType string, cloudAccountFeature *CloudauthAccountFeatureSecure) (
	*CloudauthAccountFeatureSecure, string, error) {
	payload, err := client.marshalCloudauthProto(cloudAccountFeature)
	if err != nil {
		return nil, "", err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.cloudauthAccountFeatureURL(accountID, featureType), payload)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	cloudauthAccountFeature := &CloudauthAccountFeatureSecure{}
	err = client.unmarshalCloudauthProto(response.Body, cloudauthAccountFeature)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccountFeature, "", nil
}

func (client *Client) GetCloudauthAccountFeatureSecure(ctx context.Context, accountID, featureType string) (*CloudauthAccountFeatureSecure, string, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.cloudauthAccountFeatureURL(accountID, featureType), nil)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	cloudauthAccountFeature := &CloudauthAccountFeatureSecure{}
	err = client.unmarshalCloudauthProto(response.Body, cloudauthAccountFeature)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccountFeature, "", nil
}

func (client *Client) DeleteCloudauthAccountFeatureSecure(ctx context.Context, accountID, featureType string) (string, error) {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.cloudauthAccountFeatureURL(accountID, featureType), nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorAndStatusFromResponse(response)
	}
	return "", nil
}

func (client *Client) cloudauthAccountFeatureURL(accountID string, featureType string) string {
	return fmt.Sprintf(cloudauthAccountFeaturePath, client.config.url, accountID, featureType)
}
