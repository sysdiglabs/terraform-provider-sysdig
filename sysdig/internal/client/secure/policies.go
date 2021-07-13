package secure

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigSecureClient) CreatePolicy(ctx context.Context, policyRequest Policy) (policy Policy, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPost, client.policiesURL(), policyRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return PolicyFromJSON(body), nil
}

func (client *sysdigSecureClient) policiesURL() string {
	return fmt.Sprintf("%s/api/v2/policies", client.URL)
}

func (client *sysdigSecureClient) DeletePolicy(ctx context.Context, policyID int) error {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodDelete, client.policyURL(policyID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}

	return err
}

func (client *sysdigSecureClient) policyURL(policyID int) string {
	return fmt.Sprintf("%s/api/v2/policies/%d", client.URL, policyID)
}

func (client *sysdigSecureClient) UpdatePolicy(ctx context.Context, policyRequest Policy) (policy Policy, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPut, client.policyURL(policyRequest.ID), policyRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Policy{}, errorFromResponse(response)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	return PolicyFromJSON(body), nil
}

func (client *sysdigSecureClient) GetPolicyById(ctx context.Context, policyID int) (policy Policy, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.policyURL(policyID), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Policy{}, errorFromResponse(response)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return PolicyFromJSON(body), nil
}
