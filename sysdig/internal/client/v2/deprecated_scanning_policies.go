package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	deprecatedScanningPoliciesPath        = "%s/api/scanning/v1/policies"
	deprecatedScanningPolicyPath          = "%s/api/scanning/v1/policies/%s"
	deprecatedScanningPolicyAssigmentPath = "%s/api/scanning/v1/mappings?bundleId=default"
)

type DeprecatedScanningPolicyInterface interface {
	Base
	CreateDeprecatedScanningPolicy(ctx context.Context, scanningPolicy DeprecatedScanningPolicy) (DeprecatedScanningPolicy, error)
	GetDeprecatedScanningPolicyByID(ctx context.Context, scanningPolicyID string) (DeprecatedScanningPolicy, error)
	UpdateDeprecatedScanningPolicy(ctx context.Context, scanningPolicy DeprecatedScanningPolicy) (DeprecatedScanningPolicy, error)
	DeleteDeprecatedScanningPolicyByID(ctx context.Context, scanningPolicyID string) error
}

type DeprecatedScanningPolicyAssignmentInterface interface {
	Base
	CreateDeprecatedScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentRequest DeprecatedScanningPolicyAssignmentList) (DeprecatedScanningPolicyAssignmentList, error)
	DeleteDeprecatedScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentList DeprecatedScanningPolicyAssignmentList) error
	GetDeprecatedScanningPolicyAssignmentList(ctx context.Context) (DeprecatedScanningPolicyAssignmentList, error)
}

func (c *Client) CreateDeprecatedScanningPolicy(ctx context.Context, scanningPolicy DeprecatedScanningPolicy) (policy DeprecatedScanningPolicy, err error) {
	payload, err := Marshal(scanningPolicy)
	if err != nil {
		return DeprecatedScanningPolicy{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.deprecatedScanningPoliciesURL(), payload)
	if err != nil {
		return DeprecatedScanningPolicy{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return DeprecatedScanningPolicy{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[DeprecatedScanningPolicy](response.Body)
}

func (c *Client) GetDeprecatedScanningPolicyByID(ctx context.Context, scanningPolicyID string) (policy DeprecatedScanningPolicy, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.deprecatedScanningPolicyURL(scanningPolicyID), nil)
	if err != nil {
		return DeprecatedScanningPolicy{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return DeprecatedScanningPolicy{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[DeprecatedScanningPolicy](response.Body)
}

func (c *Client) UpdateDeprecatedScanningPolicy(ctx context.Context, scanningPolicy DeprecatedScanningPolicy) (policy DeprecatedScanningPolicy, err error) {
	payload, err := Marshal(scanningPolicy)
	if err != nil {
		return DeprecatedScanningPolicy{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.deprecatedScanningPolicyURL(scanningPolicy.ID), payload)
	if err != nil {
		return DeprecatedScanningPolicy{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return DeprecatedScanningPolicy{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[DeprecatedScanningPolicy](response.Body)
}

func (c *Client) DeleteDeprecatedScanningPolicyByID(ctx context.Context, scanningPolicyID string) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deprecatedScanningPolicyURL(scanningPolicyID), nil)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return c.ErrorFromResponse(response)
	}

	return err
}

func (c *Client) CreateDeprecatedScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentList DeprecatedScanningPolicyAssignmentList) (list DeprecatedScanningPolicyAssignmentList, err error) {
	payload, err := Marshal(scanningPolicyAssignmentList)
	if err != nil {
		return DeprecatedScanningPolicyAssignmentList{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.scanningPolicyAssignmentURL(), payload)
	if err != nil {
		return DeprecatedScanningPolicyAssignmentList{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return DeprecatedScanningPolicyAssignmentList{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[DeprecatedScanningPolicyAssignmentList](response.Body)
}

func (c *Client) DeleteDeprecatedScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentList DeprecatedScanningPolicyAssignmentList) (err error) {
	payload, err := Marshal(scanningPolicyAssignmentList)
	if err != nil {
		return err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.scanningPolicyAssignmentURL(), payload)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return c.ErrorFromResponse(response)
	}

	return err
}

func (c *Client) GetDeprecatedScanningPolicyAssignmentList(ctx context.Context) (list DeprecatedScanningPolicyAssignmentList, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.scanningPolicyAssignmentURL(), nil)
	if err != nil {
		return DeprecatedScanningPolicyAssignmentList{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return DeprecatedScanningPolicyAssignmentList{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[DeprecatedScanningPolicyAssignmentList](response.Body)
}

func (c *Client) deprecatedScanningPoliciesURL() string {
	return fmt.Sprintf(deprecatedScanningPoliciesPath, c.config.url)
}

func (c *Client) deprecatedScanningPolicyURL(scanningPolicyID string) string {
	return fmt.Sprintf(deprecatedScanningPolicyPath, c.config.url, scanningPolicyID)
}

func (c *Client) scanningPolicyAssignmentURL() string {
	return fmt.Sprintf(deprecatedScanningPolicyAssigmentPath, c.config.url)
}
