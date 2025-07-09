package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	cloudauthAccountFeaturePath = "%s/api/cloudauth/v1/accounts/%s/feature/%s"
)

type CloudauthAccountFeatureSecureInterface interface {
	Base
	CreateOrUpdateCloudauthAccountFeatureSecure(ctx context.Context, accountID, featureType string, cloudAccountFeature *CloudauthAccountFeatureSecure) (*CloudauthAccountFeatureSecure, string, error)
	GetCloudauthAccountFeatureSecure(ctx context.Context, accountID, featureType string) (*CloudauthAccountFeatureSecure, string, error)
	DeleteCloudauthAccountFeatureSecure(ctx context.Context, accountID, featureType string) (string, error)
}

func (c *Client) CreateOrUpdateCloudauthAccountFeatureSecure(ctx context.Context, accountID, featureType string, cloudAccountFeature *CloudauthAccountFeatureSecure) (feature *CloudauthAccountFeatureSecure, statusCode string, err error) {
	payload, err := c.marshalCloudauthProto(cloudAccountFeature)
	if err != nil {
		return nil, "", err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.cloudauthAccountFeatureURL(accountID, featureType), payload)
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

	cloudauthAccountFeature := &CloudauthAccountFeatureSecure{}
	err = c.unmarshalCloudauthProto(response.Body, cloudauthAccountFeature)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccountFeature, "", nil
}

func (c *Client) GetCloudauthAccountFeatureSecure(ctx context.Context, accountID, featureType string) (feature *CloudauthAccountFeatureSecure, statusCode string, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.cloudauthAccountFeatureURL(accountID, featureType), nil)
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

	cloudauthAccountFeature := &CloudauthAccountFeatureSecure{}
	err = c.unmarshalCloudauthProto(response.Body, cloudauthAccountFeature)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccountFeature, "", nil
}

func (c *Client) DeleteCloudauthAccountFeatureSecure(ctx context.Context, accountID, featureType string) (statusCode string, err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.cloudauthAccountFeatureURL(accountID, featureType), nil)
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

func (c *Client) cloudauthAccountFeatureURL(accountID string, featureType string) string {
	return fmt.Sprintf(cloudauthAccountFeaturePath, c.config.url, accountID, featureType)
}
