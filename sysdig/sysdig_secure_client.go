package sysdig

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type SysdigSecureClient interface {
	CreatePolicy(Policy) (Policy, error)
	DeletePolicy(int) error
	UpdatePolicy(Policy) (Policy, error)
	GetPolicyById(int) (Policy, error)

	GetUserRulesFile() (UserRulesFile, error)
	UpdateUserRulesFile(UserRulesFile) (UserRulesFile, error)
}

func NewSysdigSecureClient(sysdigSecureAPIToken string, url string) SysdigSecureClient {
	return &sysdigSecureClient{
		SysdigSecureAPIToken: sysdigSecureAPIToken,
		URL:                  url,
		httpClient:           http.DefaultClient,
	}
}

type sysdigSecureClient struct {
	SysdigSecureAPIToken string
	URL                  string
	httpClient           *http.Client
}

func (client *sysdigSecureClient) CreatePolicy(policyRequest Policy) (Policy, error) {
	response, _ := client.doSysdigSecureRequest("POST", client.policiesURL(), policyRequest.ToJSON())
	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return Policy{}, errors.New(string(body))
	}

	defer response.Body.Close()

	return PolicyFromJSON(body), nil
}

func (client *sysdigSecureClient) doSysdigSecureRequest(method string, url string, payload io.Reader) (*http.Response, error) {
	request, _ := http.NewRequest(method, url, payload)
	request.Header.Set("Authorization", "Bearer "+client.SysdigSecureAPIToken)
	request.Header.Set("Content-Type", "application/json")

	return client.httpClient.Do(request)
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

func (client *sysdigSecureClient) GetUserRulesFile() (UserRulesFile, error) {
	response, _ := client.doSysdigSecureRequest("GET", client.userRulesFileURL(), nil)
	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return UserRulesFile{}, errors.New(string(body))
	}

	defer response.Body.Close()

	return UserRulesFileFromJSON(body), nil
}

func (client *sysdigSecureClient) userRulesFileURL() string {
	return client.URL + "/api/settings/falco/userRulesFile"
}

func (client *sysdigSecureClient) UpdateUserRulesFile(userRulesFileRequest UserRulesFile) (UserRulesFile, error) {
	response, _ := client.doSysdigSecureRequest("PUT", client.userRulesFileURL(), userRulesFileRequest.ToJSON())
	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return UserRulesFile{}, errors.New(string(body))
	}

	defer response.Body.Close()

	return UserRulesFileFromJSON(body), nil
}
