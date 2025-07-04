package v2

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

const (
	organizationsPath = "%s/api/cloudauth/v1/organizations"
	organizationPath  = "%s/api/cloudauth/v1/organizations/%s"
)

type OrganizationSecureInterface interface {
	Base
	CreateOrganizationSecure(ctx context.Context, org *OrganizationSecure) (*OrganizationSecure, string, error)
	GetOrganizationSecure(ctx context.Context, orgID string) (*OrganizationSecure, string, error)
	DeleteOrganizationSecure(ctx context.Context, orgID string) (string, error)
	UpdateOrganizationSecure(ctx context.Context, orgID string, org *OrganizationSecure) (*OrganizationSecure, string, error)
}

func (c *Client) CreateOrganizationSecure(ctx context.Context, org *OrganizationSecure) (organization *OrganizationSecure, errString string, err error) {
	payload, err := c.marshalCloudauthProto(org)
	if err != nil {
		return nil, "", err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.organizationsURL(), payload)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	err = c.unmarshalCloudauthProto(response.Body, organization)
	if err != nil {
		return nil, "", err
	}
	return organization, "", nil
}

func (c *Client) GetOrganizationSecure(ctx context.Context, orgID string) (organization *OrganizationSecure, errString string, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.organizationURL(orgID), nil)
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

	err = c.unmarshalCloudauthProto(response.Body, organization)
	if err != nil {
		return nil, "", err
	}
	return organization, "", nil
}

func (c *Client) DeleteOrganizationSecure(ctx context.Context, orgID string) (errString string, err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.organizationURL(orgID), nil)
	if err != nil {
		return "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return errStatus, err
	}
	return "", nil
}

func (c *Client) UpdateOrganizationSecure(ctx context.Context, orgID string, org *OrganizationSecure) (organization *OrganizationSecure, errString string, err error) {
	payload, err := Marshal(org)
	if err != nil {
		return nil, "", err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.organizationURL(orgID), payload)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	err = c.unmarshalCloudauthProto(response.Body, organization)
	if err != nil {
		return nil, "", err
	}
	return organization, "", nil
}

func (c *Client) organizationsURL() string {
	url := fmt.Sprintf(organizationsPath, c.config.url)
	if os.Getenv("SYSDIG_ORG_API_ASYNC") == "true" {
		url += "?async=true"
	}
	return url
}

func (c *Client) organizationURL(orgID string) string {
	url := fmt.Sprintf(organizationPath, c.config.url, orgID)
	if os.Getenv("SYSDIG_ORG_API_ASYNC") == "true" {
		url += "?async=true"
	}
	return url
}
