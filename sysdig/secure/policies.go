package secure

import (
	"errors"
	"io/ioutil"
	"strconv"
)

func (client *sysdigSecureClient) CreatePolicy(policyRequest Policy) (Policy, error) {
	response, _ := client.doSysdigSecureRequest("POST", client.policiesURL(), policyRequest.ToJSON())
	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return Policy{}, errors.New(string(body))
	}

	defer response.Body.Close()

	return PolicyFromJSON(body), nil
}

func (client *sysdigSecureClient) policiesURL() string {
	return client.URL + "/api/policies"
}

func (client *sysdigSecureClient) DeletePolicy(policyID int) error {
	response, err := client.doSysdigSecureRequest("DELETE", client.policyURL(policyID), nil)

	defer response.Body.Close()

	return err
}

func (client *sysdigSecureClient) policyURL(policyID int) string {
	return client.URL + "/api/policies/" + strconv.Itoa(policyID)
}

func (client *sysdigSecureClient) UpdatePolicy(policyRequest Policy) (Policy, error) {
	response, _ := client.doSysdigSecureRequest("PUT", client.policyURL(policyRequest.ID), policyRequest.ToJSON())
	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return Policy{}, errors.New(string(body))
	}

	defer response.Body.Close()

	return PolicyFromJSON(body), nil
}

func (client *sysdigSecureClient) GetPolicyById(policyID int) (Policy, error) {
	response, _ := client.doSysdigSecureRequest("GET", client.policyURL(policyID), nil)
	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return Policy{}, errors.New(string(body))
	}

	defer response.Body.Close()

	return PolicyFromJSON(body), nil
}
