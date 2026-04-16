package v2

import (
	"context"
	"fmt"
	"net/http"
)

const zonePolicyAssignmentPath = "%s/api/cspm/v1/zones/%d/policies"

type ZonePolicyAssignmentInterface interface {
	Base
	GetZonePolicyAssignment(ctx context.Context, zoneID int) (*ZonePolicyAssignment, error)
	CreateZonePolicyAssignment(ctx context.Context, zoneID int, req *ZonePolicyAssignmentRequest) (*ZonePolicyAssignment, error)
	UpdateZonePolicyAssignment(ctx context.Context, zoneID int, req *ZonePolicyAssignmentRequest) (*ZonePolicyAssignment, error)
	DeleteZonePolicyAssignment(ctx context.Context, zoneID int) error
}

func (c *Client) GetZonePolicyAssignment(ctx context.Context, zoneID int) (result *ZonePolicyAssignment, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getZonePolicyAssignmentURL(zoneID), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, c.APIErrorFromResponse(response)
	}

	return Unmarshal[*ZonePolicyAssignment](response.Body)
}

func (c *Client) CreateZonePolicyAssignment(ctx context.Context, zoneID int, req *ZonePolicyAssignmentRequest) (result *ZonePolicyAssignment, err error) {
	payload, err := Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.getZonePolicyAssignmentURL(zoneID), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, c.APIErrorFromResponse(response)
	}

	return Unmarshal[*ZonePolicyAssignment](response.Body)
}

func (c *Client) UpdateZonePolicyAssignment(ctx context.Context, zoneID int, req *ZonePolicyAssignmentRequest) (result *ZonePolicyAssignment, err error) {
	payload, err := Marshal(req)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.getZonePolicyAssignmentURL(zoneID), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, c.APIErrorFromResponse(response)
	}

	return Unmarshal[*ZonePolicyAssignment](response.Body)
}

func (c *Client) DeleteZonePolicyAssignment(ctx context.Context, zoneID int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.getZonePolicyAssignmentURL(zoneID), nil)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return c.APIErrorFromResponse(response)
	}

	return nil
}

func (c *Client) getZonePolicyAssignmentURL(zoneID int) string {
	return fmt.Sprintf(zonePolicyAssignmentPath, c.config.url, zoneID)
}
