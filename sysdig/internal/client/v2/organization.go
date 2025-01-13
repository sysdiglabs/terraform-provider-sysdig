package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	organizationsPath = "%s/api/cloudauth/v1/organizations?async=true"
	organizationPath  = "%s/api/cloudauth/v1/organizations/%s?async=true"
)

type OrganizationSecureInterface interface {
	Base
	CreateOrganizationSecure(ctx context.Context, org *OrganizationSecure) (*OrganizationSecure, string, error)
	GetOrganizationSecure(ctx context.Context, orgID string) (*OrganizationSecure, string, error)
	DeleteOrganizationSecure(ctx context.Context, orgID string) (string, error)
	UpdateOrganizationSecure(ctx context.Context, orgID string, org *OrganizationSecure) (*OrganizationSecure, string, error)
}

func (client *Client) CreateOrganizationSecure(ctx context.Context, org *OrganizationSecure) (*OrganizationSecure, string, error) {
	payload, err := client.marshalCloudauthProto(org)
	if err != nil {
		return nil, "", err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.organizationsURL(), payload)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	organization := &OrganizationSecure{}
	err = client.unmarshalCloudauthProto(response.Body, organization)
	if err != nil {
		return nil, "", err
	}
	return organization, "", nil
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

	organization := &OrganizationSecure{}
	err = client.unmarshalCloudauthProto(response.Body, organization)
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

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	organization := &OrganizationSecure{}
	err = client.unmarshalCloudauthProto(response.Body, organization)
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
