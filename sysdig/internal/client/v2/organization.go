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
	organizationsPath = "%s/api/cloudauth/v1/organizations"
	organizationPath  = "%s/api/cloudauth/v1/organizations/%s"
)

type OrganizationSecureInterface interface {
	Base
	CreateOrganizationSecure(ctx context.Context, org *OrganizationSecure) (*OrganizationSecure, error)
	GetOrganizationSecure(ctx context.Context, orgID string) (*OrganizationSecure, string, error)
	DeleteOrganizationSecure(ctx context.Context, orgID string) (string, error)
	UpdateOrganizationSecure(ctx context.Context, orgID string, org *OrganizationSecure) (*OrganizationSecure, string, error)
}

func (client *Client) CreateOrganizationSecure(ctx context.Context, org *OrganizationSecure) (*OrganizationSecure, error) {
	payload, err := client.marshalOrg(org)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.organizationsURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	return client.unmarshalOrg(response.Body)
}

func (client *Client) GetOrganizationSecure(ctx context.Context, orgID string) (*OrganizationSecure, string, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.organizationURL(orgID), nil)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	organization, err := client.unmarshalOrg(response.Body)
	if err != nil {
		return nil, "", err
	}
	return organization, "", nil
}

func (client *Client) DeleteOrganizationSecure(ctx context.Context, orgID string) (string, error) {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.organizationURL(orgID), nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return errStatus, err
	}
	return "", nil
}

func (client *Client) UpdateOrganizationSecure(ctx context.Context, orgID string, org *OrganizationSecure) (*OrganizationSecure, string, error) {
	payload, err := Marshal(org)
	if err != nil {
		return nil, "", err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.organizationURL(orgID), payload)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	organization, err := client.unmarshalOrg(response.Body)
	if err != nil {
		return nil, "", err
	}
	return organization, "", nil
}

func (client *Client) organizationsURL() string {
	return fmt.Sprintf(organizationsPath, client.config.url)
}

func (client *Client) organizationURL(orgId string) string {
	return fmt.Sprintf(organizationPath, client.config.url, orgId)
}

// local function for protojson based marshal/unmarshal of organization proto
func (client *Client) marshalOrg(data *OrganizationSecure) (io.Reader, error) {
	payload, err := protojson.Marshal(data)
	return bytes.NewBuffer(payload), err
}

func (client *Client) unmarshalOrg(data io.ReadCloser) (*OrganizationSecure, error) {
	result := &OrganizationSecure{}

	body, err := io.ReadAll(data)
	if err != nil {
		return result, err
	}

	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(body, result)
	return result, err
}
