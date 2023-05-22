package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	scanningPoliciesPath        = "%s/api/scanning/v1/policies"
	scanningPolicyPath          = "%s/api/scanning/v1/policies/%s"
	scanningPolicyAssigmentPath = "%s/api/scanning/v1/mappings?bundleId=default"
)

type ScanningPolicyInterface interface {
	Base
	CreateScanningPolicy(ctx context.Context, scanningPolicy ScanningPolicy) (ScanningPolicy, error)
	GetScanningPolicyByID(ctx context.Context, scanningPolicyID string) (ScanningPolicy, error)
	UpdateScanningPolicyByID(ctx context.Context, scanningPolicy ScanningPolicy) (ScanningPolicy, error)
	DeleteScanningPolicyByID(ctx context.Context, scanningPolicyID string) error
}

type ScanningPolicyAssignmentInterface interface {
	Base
	CreateScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentRequest ScanningPolicyAssignmentList) (ScanningPolicyAssignmentList, error)
	DeleteScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentList ScanningPolicyAssignmentList) error
	GetScanningPolicyAssignmentList(ctx context.Context) (ScanningPolicyAssignmentList, error)
}

func (client *Client) CreateScanningPolicy(ctx context.Context, scanningPolicy ScanningPolicy) (ScanningPolicy, error) {
	payload, err := Marshal(scanningPolicy)
	if err != nil {
		return ScanningPolicy{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.scanningPoliciesURL(), payload)
	if err != nil {
		return ScanningPolicy{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ScanningPolicy{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[ScanningPolicy](response.Body)
}

func (client *Client) GetScanningPolicyByID(ctx context.Context, scanningPolicyID string) (ScanningPolicy, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.scanningPolicyURL(scanningPolicyID), nil)
	if err != nil {
		return ScanningPolicy{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ScanningPolicy{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[ScanningPolicy](response.Body)
}

func (client *Client) UpdateScanningPolicyByID(ctx context.Context, scanningPolicy ScanningPolicy) (ScanningPolicy, error) {
	payload, err := Marshal(scanningPolicy)
	if err != nil {
		return ScanningPolicy{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.scanningPolicyURL(scanningPolicy.ID), payload)
	if err != nil {
		return ScanningPolicy{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ScanningPolicy{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[ScanningPolicy](response.Body)
}

func (client *Client) DeleteScanningPolicyByID(ctx context.Context, scanningPolicyID string) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.scanningPolicyURL(scanningPolicyID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorFromResponse(response)
	}

	return err
}

func (client *Client) CreateScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentList ScanningPolicyAssignmentList) (ScanningPolicyAssignmentList, error) {
	payload, err := Marshal(scanningPolicyAssignmentList)
	if err != nil {
		return ScanningPolicyAssignmentList{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.scanningPolicyAssignmentURL(), payload)
	if err != nil {
		return ScanningPolicyAssignmentList{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ScanningPolicyAssignmentList{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[ScanningPolicyAssignmentList](response.Body)
}

func (client *Client) DeleteScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentList ScanningPolicyAssignmentList) error {
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

func (client *Client) GetScanningPolicyAssignmentList(ctx context.Context) (ScanningPolicyAssignmentList, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.scanningPolicyAssignmentURL(), nil)
	if err != nil {
		return ScanningPolicyAssignmentList{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ScanningPolicyAssignmentList{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[ScanningPolicyAssignmentList](response.Body)
}

func (client *Client) scanningPoliciesURL() string {
	return fmt.Sprintf(scanningPoliciesPath, client.config.url)
}

func (client *Client) scanningPolicyURL(scanningPolicyID string) string {
	return fmt.Sprintf(scanningPolicyPath, client.config.url, scanningPolicyID)
}

func (client *Client) scanningPolicyAssignmentURL() string {
	return fmt.Sprintf(scanningPolicyAssigmentPath, client.config.url)
}
