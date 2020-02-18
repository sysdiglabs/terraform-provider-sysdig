package secure

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigSecureClient) CreatePolicy(policyRequest Policy) (policy Policy, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPost, client.policiesURL(), policyRequest.ToJSON())
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		return Policy{}, errors.New(string(body))
	}

	defer response.Body.Close()

	return PolicyFromJSON(body), nil
}

func (client *sysdigSecureClient) policiesURL() string {
	return fmt.Sprintf("%s/api/v2/policies", client.URL)
}

func (client *sysdigSecureClient) DeletePolicy(policyID int) error {
	response, err := client.doSysdigSecureRequest(http.MethodDelete, client.policyURL(policyID), nil)

	defer response.Body.Close()

	return err
}

func (client *sysdigSecureClient) policyURL(policyID int) string {
	return fmt.Sprintf("%s/api/v2/policies/%d", client.URL, policyID)

}

func (client *sysdigSecureClient) UpdatePolicy(policyRequest Policy) (policy Policy, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPut, client.policyURL(policyRequest.ID), policyRequest.ToJSON())
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		return Policy{}, errors.New(string(body))
	}

	defer response.Body.Close()

	return PolicyFromJSON(body), nil
}

func (client *sysdigSecureClient) GetPolicyById(policyID int) (policy Policy, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodGet, client.policyURL(policyID), nil)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		return Policy{}, errors.New(string(body))
	}

	defer response.Body.Close()

	return PolicyFromJSON(body), nil
}
