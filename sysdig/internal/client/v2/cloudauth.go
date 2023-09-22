package v2

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
)

const (
	cloudauthAccountsPath = "%s/api/cloudauth/v1/accounts"
	cloudauthAccountPath  = "%s/api/cloudauth/v1/accounts/%s"
)

type CloudauthAccountSecureInterface interface {
	Base
	CreateCloudauthAccountSecure(ctx context.Context, cloudAccount *CloudauthAccountSecure) (*CloudauthAccountSecure, error)
	GetCloudauthAccountSecure(ctx context.Context, accountID string) (*CloudauthAccountSecure, error)
	DeleteCloudauthAccountSecure(ctx context.Context, accountID string) error
	UpdateCloudauthAccountSecure(ctx context.Context, accountID string, cloudAccount *CloudauthAccountSecure) (*CloudauthAccountSecure, error)
}

func (client *Client) CreateCloudauthAccountSecure(ctx context.Context, cloudAccount *CloudauthAccountSecure) (*CloudauthAccountSecure, error) {
	payload, err := client.marshalProto(cloudAccount)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.cloudauthAccountsURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	return client.unmarshalProto(response.Body)
}

func (client *Client) GetCloudauthAccountSecure(ctx context.Context, accountID string) (*CloudauthAccountSecure, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.cloudauthAccountURL(accountID), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return client.unmarshalProto(response.Body)
}

func (client *Client) DeleteCloudauthAccountSecure(ctx context.Context, accountID string) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.cloudauthAccountURL(accountID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorFromResponse(response)
	}
	return nil
}

func (client *Client) UpdateCloudauthAccountSecure(ctx context.Context, accountID string, cloudAccount *CloudauthAccountSecure) (*CloudauthAccountSecure, error) {
	payload, err := client.marshalProto(cloudAccount)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.cloudauthAccountURL(accountID), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	return client.unmarshalProto(response.Body)
}

func (client *Client) cloudauthAccountsURL() string {
	return fmt.Sprintf(cloudauthAccountsPath, client.config.url)
}

func (client *Client) cloudauthAccountURL(accountID string) string {
	return fmt.Sprintf(cloudauthAccountPath, client.config.url, accountID)
}

// local function for protojson based marshal/unmarshal of cloudauthAccount proto
func (client *Client) marshalProto(data *CloudauthAccountSecure) (io.Reader, error) {
	payload, err := protojson.Marshal(data)
	return bytes.NewBuffer(payload), err
}

func (client *Client) unmarshalProto(data io.ReadCloser) (*CloudauthAccountSecure, error) {
	result := &CloudauthAccountSecure{}

	body, err := io.ReadAll(data)
	if err != nil {
		return result, err
	}

	err = protojson.Unmarshal(body, result)
	return result, err
}
