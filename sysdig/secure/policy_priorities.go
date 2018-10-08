package secure

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigSecureClient) CreatePoliciesPriority(priorityRequest PoliciesPriority) (priority PoliciesPriority, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPut, client.policiesPriorityURL(), priorityRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	priority = PoliciesPriorityFromJSON(body)
	return

}

// Same behaviour as Create
func (client *sysdigSecureClient) UpdatePoliciesPriority(priorityRequest PoliciesPriority) (priority PoliciesPriority, err error) {
	return client.CreatePoliciesPriority(priorityRequest)
}

func (client *sysdigSecureClient) GetPoliciesPriority() (priority PoliciesPriority, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodGet, client.policiesPriorityURL(), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	priority = PoliciesPriorityFromJSON(body)
	return
}

func (client *sysdigSecureClient) policiesPriorityURL() string {
	return fmt.Sprintf("%s/api/policies/priorities", client.URL)
}
