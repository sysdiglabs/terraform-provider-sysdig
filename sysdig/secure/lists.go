package secure

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigSecureClient) CreateList(listRequest List) (list List, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPost, client.GetListsUrl(), listRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errors.New(response.Status)
		return
	}

	list, err = ListFromJSON(body)
	return
}

func (client *sysdigSecureClient) GetListById(id int) (list List, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodGet, client.GetListUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	list, err = ListFromJSON(body)
	if err != nil {
		return
	}

	if list.Version == 0 {
		err = fmt.Errorf("List with ID: %d does not exists", id)
		return
	}
	return
}

func (client *sysdigSecureClient) UpdateList(listRequest List) (list List, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPut, client.GetListUrl(listRequest.ID), listRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	return ListFromJSON(body)
}

func (client *sysdigSecureClient) DeleteList(id int) error {
	response, err := client.doSysdigSecureRequest(http.MethodDelete, client.GetListUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}
	return nil
}

func (client *sysdigSecureClient) GetListsUrl() string {
	return fmt.Sprintf("%s/api/secure/falco/lists", client.URL)
}

func (client *sysdigSecureClient) GetListUrl(id int) string {
	return fmt.Sprintf("%s/api/secure/falco/lists/%d", client.URL, id)
}
