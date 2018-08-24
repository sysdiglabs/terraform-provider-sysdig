package sysdig

import (
	"errors"
	"fmt"
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

	CreateNotificationChannel(NotificationChannel) (NotificationChannel, error)
	GetNotificationChannelById(int) (NotificationChannel, error)
	DeleteNotificationChannel(int) error
	UpdateNotificationChannel(NotificationChannel) (NotificationChannel, error)
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

// == NotificationChannel ==============================================================================================

func (client *sysdigSecureClient) GetNotificationChannelById(id int) (nc NotificationChannel, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodGet, client.GetNotificationChannelUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	nc = NotificationChannelFromJSON(body)
	return
}

func (client *sysdigSecureClient) CreateNotificationChannel(ncRequest NotificationChannel) (nc NotificationChannel, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPost, client.GetNotificationChannelsUrl(), ncRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errors.New(response.Status)
		return
	}

	nc = NotificationChannelFromJSON(body)
	return
}

func (client *sysdigSecureClient) UpdateNotificationChannel(ncRequest NotificationChannel) (nc NotificationChannel, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPut, client.GetNotificationChannelUrl(ncRequest.ID), ncRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	nc = NotificationChannelFromJSON(body)
	return

}

func (client *sysdigSecureClient) DeleteNotificationChannel(id int) error {
	response, err := client.doSysdigSecureRequest(http.MethodDelete, client.GetNotificationChannelUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return errors.New(response.Status)
	}
	return nil
}

func (client *sysdigSecureClient) GetNotificationChannelsUrl() string {
	return fmt.Sprintf("%s/api/notificationChannels", client.URL)
}

func (client *sysdigSecureClient) GetNotificationChannelUrl(id int) string {
	return fmt.Sprintf("%s/api/notificationChannels/%d", client.URL, id)
}

// == Policy ===========================================================================================================
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
