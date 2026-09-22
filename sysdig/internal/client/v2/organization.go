package v2

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
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

	response, err := c.requester.Request(ctx, http.MethodPost, c.withAsync(c.organizationsURL()), payload)
	if err != nil {
		return nil, "", err
	}
	// a close failure says nothing the caller can act on, and returning it would discard a
	// successful call: a created organization would exist server side with nothing tracking it
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	organization = &OrganizationSecure{}
	if err = c.unmarshalOrganizationBody(response, organization); err != nil {
		return nil, "", err
	}
	return organization, "", nil
}

func (c *Client) GetOrganizationSecure(ctx context.Context, orgID string) (organization *OrganizationSecure, errString string, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.organizationURL(orgID), nil)
	if err != nil {
		return nil, "", err
	}
	// a close failure says nothing the caller can act on, and returning it would discard a
	// successful call: a created organization would exist server side with nothing tracking it
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	organization = &OrganizationSecure{}
	err = c.unmarshalCloudauthProto(response.Body, organization)
	if err != nil {
		return nil, "", err
	}
	return organization, "", nil
}

func (c *Client) DeleteOrganizationSecure(ctx context.Context, orgID string) (errString string, err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.withAsync(c.organizationURL(orgID)), nil)
	if err != nil {
		return "", err
	}
	// a close failure says nothing the caller can act on, and returning it would discard a
	// successful call: a created organization would exist server side with nothing tracking it
	defer func() { _ = response.Body.Close() }()

	// 202 is only expected when async was requested: accepting it otherwise would let a destroy
	// return on an acknowledgement while the resource, seeing the flag off, skips the wait
	acknowledged := c.config.secureOrgAPIAsync && response.StatusCode == http.StatusAccepted
	if !acknowledged && response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return errStatus, err
	}
	return "", nil
}

func (c *Client) UpdateOrganizationSecure(ctx context.Context, orgID string, org *OrganizationSecure) (organization *OrganizationSecure, errString string, err error) {
	payload, err := c.marshalCloudauthProto(org)
	if err != nil {
		return nil, "", err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.withAsync(c.organizationURL(orgID)), payload)
	if err != nil {
		return nil, "", err
	}
	// a close failure says nothing the caller can act on, and returning it would discard a
	// successful call: a created organization would exist server side with nothing tracking it
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	organization = &OrganizationSecure{}
	if err := c.unmarshalOrganizationBody(response, organization); err != nil {
		return nil, "", err
	}
	return organization, "", nil
}

// an async acknowledgement may carry no body at all; anything else still has to decode cleanly,
// so a malformed payload is not mistaken for an empty one
func (c *Client) unmarshalOrganizationBody(response *http.Response, organization *OrganizationSecure) error {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusAccepted && len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	return c.unmarshalCloudauthProto(io.NopCloser(bytes.NewReader(body)), organization)
}

// deliberately not applied to the read: GetOrganizationSecure accepts only 200
func (c *Client) withAsync(url string) string {
	if !c.config.secureOrgAPIAsync {
		return url
	}
	return url + "?async=true"
}

func (c *Client) organizationsURL() string {
	return fmt.Sprintf(organizationsPath, c.config.url)
}

func (c *Client) organizationURL(orgID string) string {
	return fmt.Sprintf(organizationPath, c.config.url, orgID)
}
