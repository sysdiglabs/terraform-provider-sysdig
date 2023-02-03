package secure

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (client *sysdigSecureClient) scanningPoliciesURL() string {
	return fmt.Sprintf("%s/api/scanning/v1/policies", client.URL)
}

func (client *sysdigSecureClient) scanningPolicyAssignmentURL() string {
	return fmt.Sprintf("%s/api/scanning/v1/mappings?bundleId=default", client.URL)
}

func (client *sysdigSecureClient) scanningPolicyURL(scanningPolicyId string) string {
	return fmt.Sprintf("%s/api/scanning/v1/policies/%s", client.URL, scanningPolicyId)
}

// Scanning Policies

func (client *sysdigSecureClient) CreateScanningPolicy(ctx context.Context, scanningPolicyRequest ScanningPolicy) (scanningPolicy ScanningPolicy, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPost, client.scanningPoliciesURL(), scanningPolicyRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	return ScanningPolicyFromJSON(body), nil
}

func (client *sysdigSecureClient) GetScanningPolicyById(ctx context.Context, scanningPolicyID string) (scanningPolicy ScanningPolicy, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.scanningPolicyURL(scanningPolicyID), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ScanningPolicy{}, errorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	return ScanningPolicyFromJSON(body), nil
}

func (client *sysdigSecureClient) UpdateScanningPolicyById(ctx context.Context, scanningPolicyRequest ScanningPolicy) (scanningPolicy ScanningPolicy, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPut, client.scanningPolicyURL(scanningPolicyRequest.ID), scanningPolicyRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ScanningPolicy{}, errorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	return ScanningPolicyFromJSON(body), nil
}

func (client *sysdigSecureClient) DeleteScanningPolicyById(ctx context.Context, scanningPolicyID string) error {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodDelete, client.scanningPolicyURL(scanningPolicyID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}

	return err
}

// Scanning Policy Assignments

func (client *sysdigSecureClient) CreateScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentRequest ScanningPolicyAssignmentList) (scanningPolicyAssignmentList ScanningPolicyAssignmentList, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPut, client.scanningPolicyAssignmentURL(), scanningPolicyAssignmentRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	return ScanningPolicyAssignmentFromJSON(body), nil
}

func (client *sysdigSecureClient) DeleteScanningPolicyAssignmentList(ctx context.Context, scanningPolicyAssignmentList ScanningPolicyAssignmentList) error {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPut, client.scanningPolicyAssignmentURL(), scanningPolicyAssignmentList.ToJSON())
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}

	return err
}

func (client *sysdigSecureClient) GetScanningPolicyAssignmentList(ctx context.Context) (scanningPolicyAssignmentList ScanningPolicyAssignmentList, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.scanningPolicyAssignmentURL(), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ScanningPolicyAssignmentList{}, errorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	return ScanningPolicyAssignmentFromJSON(body), nil
}
