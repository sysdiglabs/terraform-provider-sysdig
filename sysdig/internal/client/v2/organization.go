package v2

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	organizationsPath = "%s/api/cloudauth/v1/organizations"
	organizationPath  = "%s/api/cloudauth/v1/organizations/%s"

	// identifying information only: the default verbosity makes the backend read every member
	// account, and the deletion probe only needs to know whether the organization row is there
	organizationProbeVerbosity = "VERBOSITY_IDENT"
)

type OrganizationSecureInterface interface {
	Base
	OrganizationExistsSecure(ctx context.Context, orgID string) (bool, string, error)
	OrgAPIAsyncEnabled() bool
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

	acknowledged := c.config.secureOrgAPIAsync && response.StatusCode == http.StatusAccepted
	if !acknowledged && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	// the create still has to read the body: without the id there is nothing to track
	organization = &OrganizationSecure{}
	if _, err = c.unmarshalOrganizationBody(response, organization, c.config.secureOrgAPIAsync); err != nil {
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
	if _, err = c.unmarshalOrganizationBody(response, organization, false); err != nil {
		return nil, "", err
	}
	return organization, "", nil
}

// OrgAPIAsyncEnabled reports whether the organization API is being called asynchronously. The
// caller needs it to know whether a delete was merely accepted; reading it here keeps the answer
// off the provider's shared ResourceData, which other resources read without a lock.
func (c *Client) OrgAPIAsyncEnabled() bool {
	return c.config.secureOrgAPIAsync
}

// OrganizationExistsSecure answers whether the organization is still there, and nothing else.
// Identity verbosity keeps the backend from listing every member account, which is megabytes per
// call on the organizations this path exists for, and the transport is told not to retry so the
// caller sees the status the server returned rather than spending its attempt on backoff.
func (c *Client) OrganizationExistsSecure(ctx context.Context, orgID string) (exists bool, errString string, err error) {
	url := fmt.Sprintf("%s?verbosity=%s", c.organizationURL(orgID), organizationProbeVerbosity)
	response, err := c.requester.Request(WithoutTransportRetries(ctx), http.MethodGet, url, nil)
	if err != nil {
		return false, "", err
	}
	defer func() { _ = response.Body.Close() }()

	switch response.StatusCode {
	// the body is deliberately not read: an answer that does not parse, a proxy error page say,
	// still proves the organization is there
	case http.StatusOK:
		return true, "", nil
	case http.StatusNotFound, http.StatusGone:
		return false, "", nil
	}
	errStatus, err := c.ErrorAndStatusFromResponse(response)
	return false, errStatus, err
}

func (c *Client) DeleteOrganizationSecure(ctx context.Context, orgID string) (errString string, err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.withAsync(c.organizationURL(orgID)), nil)
	if err != nil {
		return "", err
	}
	// a close failure says nothing the caller can act on, and returning it would discard a
	// successful call: a created organization would exist server side with nothing tracking it
	defer func() { _ = response.Body.Close() }()

	// an async delete answers 204 as well, with the cascade still running in the background, so
	// there is no acknowledgement status to accept here
	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
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

	acknowledged := c.config.secureOrgAPIAsync && response.StatusCode == http.StatusAccepted
	if !acknowledged && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	organization = &OrganizationSecure{}
	var decoded bool
	decoded, err = c.unmarshalOrganizationBody(response, organization, c.config.secureOrgAPIAsync)
	if err != nil {
		return nil, "", err
	}
	// an acknowledgement says the update was accepted, not what was persisted: reporting no
	// organization keeps the caller on what it planned instead of on a possibly partial body
	if acknowledged || !decoded {
		return nil, "", nil
	}
	return organization, "", nil
}

// allowEmptyAck belongs to the caller, not to the status: an async mutation may be acknowledged
// with nothing whichever 2xx carries it, while a read always needs the organization itself. The
// first return says whether an organization was decoded at all.
func (c *Client) unmarshalOrganizationBody(response *http.Response, organization *OrganizationSecure, allowEmptyAck bool) (bool, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}
	if len(bytes.TrimSpace(body)) == 0 {
		if allowEmptyAck {
			return false, nil
		}
		return false, errors.New("the organization API answered with no body")
	}
	if err := c.unmarshalCloudauthProto(io.NopCloser(bytes.NewReader(body)), organization); err != nil {
		return false, err
	}
	return true, nil
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
