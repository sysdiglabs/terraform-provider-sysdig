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
	UpdateDeprecatedScanningPolicyByID(ctx context.Context, scanningPolicy DeprecatedScanningPolicy) (DeprecatedScanningPolicy, error)
	DeleteDeprecatedScanningPolicyByID(ctx context.Context, scanningPolicyID string) error
}

type DeprecatedScanningPolicyAssignmentInterface interface {
	Base
	CreateDeprecatedScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentRequest DeprecatedScanningPolicyAssignmentList) (DeprecatedScanningPolicyAssignmentList, error)
	DeleteDeprecatedScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentList DeprecatedScanningPolicyAssignmentList) error
	GetDeprecatedScanningPolicyAssignmentList(ctx context.Context) (DeprecatedScanningPolicyAssignmentList, error)
}

func (client *Client) CreateDeprecatedScanningPolicy(ctx context.Context, scanningPolicy DeprecatedScanningPolicy) (DeprecatedScanningPolicy, error) {
	payload, err := Marshal(scanningPolicy)
	if err != nil {
		return DeprecatedScanningPolicy{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.deprecatedScanningPoliciesURL(), payload)
	if err != nil {
		return DeprecatedScanningPolicy{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return DeprecatedScanningPolicy{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[DeprecatedScanningPolicy](response.Body)
}

func (client *Client) GetDeprecatedScanningPolicyByID(ctx context.Context, scanningPolicyID string) (DeprecatedScanningPolicy, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.deprecatedScanningPolicyURL(scanningPolicyID), nil)
	if err != nil {
		return DeprecatedScanningPolicy{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return DeprecatedScanningPolicy{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[DeprecatedScanningPolicy](response.Body)
}

func (client *Client) UpdateDeprecatedScanningPolicyByID(ctx context.Context, scanningPolicy DeprecatedScanningPolicy) (DeprecatedScanningPolicy, error) {
	payload, err := Marshal(scanningPolicy)
	if err != nil {
		return DeprecatedScanningPolicy{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.deprecatedScanningPolicyURL(scanningPolicy.ID), payload)
	if err != nil {
		return DeprecatedScanningPolicy{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return DeprecatedScanningPolicy{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[DeprecatedScanningPolicy](response.Body)
}

func (client *Client) DeleteDeprecatedScanningPolicyByID(ctx context.Context, scanningPolicyID string) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.deprecatedScanningPolicyURL(scanningPolicyID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorFromResponse(response)
	}

	return err
}

func (client *Client) CreateDeprecatedScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentList DeprecatedScanningPolicyAssignmentList) (DeprecatedScanningPolicyAssignmentList, error) {
	payload, err := Marshal(scanningPolicyAssignmentList)
	if err != nil {
		return DeprecatedScanningPolicyAssignmentList{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.scanningPolicyAssignmentURL(), payload)
	if err != nil {
		return DeprecatedScanningPolicyAssignmentList{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return DeprecatedScanningPolicyAssignmentList{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[DeprecatedScanningPolicyAssignmentList](response.Body)
}

func (client *Client) DeleteDeprecatedScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentList DeprecatedScanningPolicyAssignmentList) error {
	payload, err := Marshal(scanningPolicyAssignmentList)
	if err != nil {
		return err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.scanningPolicyAssignmentURL(), payload)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorFromResponse(response)
	}

	return err
}

func (client *Client) GetDeprecatedScanningPolicyAssignmentList(ctx context.Context) (DeprecatedScanningPolicyAssignmentList, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.scanningPolicyAssignmentURL(), nil)
	if err != nil {
		return DeprecatedScanningPolicyAssignmentList{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return DeprecatedScanningPolicyAssignmentList{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[DeprecatedScanningPolicyAssignmentList](response.Body)
}

func (client *Client) deprecatedScanningPoliciesURL() string {
	return fmt.Sprintf(deprecatedScanningPoliciesPath, client.config.url)
}

func (client *Client) deprecatedScanningPolicyURL(scanningPolicyID string) string {
	return fmt.Sprintf(deprecatedScanningPolicyPath, client.config.url, scanningPolicyID)
}

func (client *Client) scanningPolicyAssignmentURL() string {
	return fmt.Sprintf(deprecatedScanningPolicyAssigmentPath, client.config.url)
}
