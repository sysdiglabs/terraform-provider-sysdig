package v2

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	cloudauthAccountsPath   = "%s/api/cloudauth/v1/accounts"
	cloudauthAccountPath    = "%s/api/cloudauth/v1/accounts/%s"
	getCloudauthAccountPath = "%s/api/cloudauth/v1/accounts/%s?decrypt=%s"
)

type CloudauthAccountSecureInterface interface {
	Base
	CreateCloudauthAccountSecure(ctx context.Context, cloudAccount *CloudauthAccountSecure) (*CloudauthAccountSecure, string, error)
	GetCloudauthAccountSecureByID(ctx context.Context, accountID string) (*CloudauthAccountSecure, string, error)
	DeleteCloudauthAccountSecure(ctx context.Context, accountID string) (string, error)
	UpdateCloudauthAccountSecure(ctx context.Context, accountID string, cloudAccount *CloudauthAccountSecure) (*CloudauthAccountSecure, string, error)
}

func (c *Client) CreateCloudauthAccountSecure(ctx context.Context, cloudAccount *CloudauthAccountSecure) (account *CloudauthAccountSecure, errStatus string, err error) {
	payload, err := c.marshalCloudauthProto(cloudAccount)
	if err != nil {
		return nil, "", err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.cloudauthAccountsURL(), payload)
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

	cloudauthAccount := &CloudauthAccountSecure{}
	err = c.unmarshalCloudauthProto(response.Body, cloudauthAccount)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccount, "", nil
}

func (c *Client) GetCloudauthAccountSecureByID(ctx context.Context, accountID string) (account *CloudauthAccountSecure, errStatus string, err error) {
	// get the cloud account with decrypt query param true to fetch decrypted details on the cloud account
	response, err := c.requester.Request(ctx, http.MethodGet, c.getCloudauthAccountURL(accountID, "true"), nil)
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

	cloudauthAccount := &CloudauthAccountSecure{}
	err = c.unmarshalCloudauthProto(response.Body, cloudauthAccount)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccount, "", nil
}

func (c *Client) DeleteCloudauthAccountSecure(ctx context.Context, accountID string) (errStatus string, err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.cloudauthAccountURL(accountID), nil)
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

func (c *Client) UpdateCloudauthAccountSecure(ctx context.Context, accountID string, cloudAccount *CloudauthAccountSecure) (account *CloudauthAccountSecure, errString string, err error) {
	payload, err := c.marshalCloudauthProto(cloudAccount)
	if err != nil {
		return nil, "", err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.cloudauthAccountURL(accountID), payload)
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

	cloudauthAccount := &CloudauthAccountSecure{}
	err = c.unmarshalCloudauthProto(response.Body, cloudauthAccount)
	if err != nil {
		return nil, "", err
	}
	return cloudauthAccount, "", nil
}

func (c *Client) cloudauthAccountsURL() string {
	return fmt.Sprintf(cloudauthAccountsPath, c.config.url)
}

func (c *Client) cloudauthAccountURL(accountID string) string {
	return fmt.Sprintf(cloudauthAccountPath, c.config.url, accountID)
}

func (c *Client) getCloudauthAccountURL(accountID string, decrypt string) string {
	return fmt.Sprintf(getCloudauthAccountPath, c.config.url, accountID, decrypt)
}

// common func for protojson based marshal/unmarshal of any cloudauth proto
func (c *Client) marshalCloudauthProto(message proto.Message) (io.Reader, error) {
	payload, err := protojson.Marshal(message)
	return bytes.NewBuffer(payload), err
}

func (c *Client) unmarshalCloudauthProto(data io.ReadCloser, message proto.Message) error {
	body, err := io.ReadAll(data)
	if err != nil {
		return err
	}

	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(body, message)
	return err
}

func (c *Client) ErrorAndStatusFromResponse(response *http.Response) (string, error) {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return response.Status, err
	}
	return response.Status, errors.New(string(b))
}
