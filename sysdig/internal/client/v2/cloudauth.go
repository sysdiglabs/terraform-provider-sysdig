package v2

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
)

const (
	cloudauthAccountsPath   = "%s/api/cloudauth/v1/accounts"
	cloudauthAccountPath    = "%s/api/cloudauth/v1/accounts/%s"
	getCloudauthAccountPath = "%s/api/cloudauth/v1/accounts/%s?decrypt=%s"
)

type CloudauthAccountSecureInterface interface {
	Base
	CreateCloudauthAccountSecure(ctx context.Context, cloudAccount *CloudauthAccountSecure) (*CloudauthAccountSecure, error)
	GetCloudauthAccountSecure(ctx context.Context, accountID string) (*CloudauthAccountSecure, string, error)
	DeleteCloudauthAccountSecure(ctx context.Context, accountID string) (string, error)
	UpdateCloudauthAccountSecure(ctx context.Context, accountID string, cloudAccount *CloudauthAccountSecure) (*CloudauthAccountSecure, string, error)
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

func (client *Client) GetCloudauthAccountSecure(ctx context.Context, accountID string) (*CloudauthAccountSecure, string, error) {
	// get the cloud account with decrypt query param true to fetch decrypted details on the cloud account
	response, err := client.requester.Request(ctx, http.MethodGet, client.getCloudauthAccountURL(accountID, "true"), nil)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	cloudauthAccount, err := client.unmarshalProto(response.Body)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccount, "", nil
}

func (client *Client) DeleteCloudauthAccountSecure(ctx context.Context, accountID string) (string, error) {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.cloudauthAccountURL(accountID), nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorAndStatusFromResponse(response)
	}
	return "", nil
}

func (client *Client) UpdateCloudauthAccountSecure(ctx context.Context, accountID string, cloudAccount *CloudauthAccountSecure) (
	*CloudauthAccountSecure, string, error) {
	payload, err := client.marshalProto(cloudAccount)
	if err != nil {
		return nil, "", err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.cloudauthAccountURL(accountID), payload)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	cloudauthAccount, err := client.unmarshalProto(response.Body)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccount, "", nil
}

func (client *Client) cloudauthAccountsURL() string {
	return fmt.Sprintf(cloudauthAccountsPath, client.config.url)
}

func (client *Client) cloudauthAccountURL(accountID string) string {
	return fmt.Sprintf(cloudauthAccountPath, client.config.url, accountID)
}

func (client *Client) getCloudauthAccountURL(accountID string, decrypt string) string {
	return fmt.Sprintf(getCloudauthAccountPath, client.config.url, accountID, decrypt)
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

	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(body, result)
	return result, err
}

func (client *Client) ErrorAndStatusFromResponse(response *http.Response) (string, error) {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return response.Status, err
	}
	return response.Status, errors.New(string(b))
}
