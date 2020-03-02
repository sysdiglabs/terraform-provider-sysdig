package secure

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigSecureClient) GetUserById(id int) (u User, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodGet, client.GetUserUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	u = UsersFromJSON(body)

	return
}

func (client *sysdigSecureClient) CreateUser(uRequest User) (u User, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPost, client.GetUsersUrl(), uRequest.ToJSON())

	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errors.New(response.Status)
		return
	}

	u = UsersFromJSON(body)
	return
}

func (client *sysdigSecureClient) UpdateUser(uRequest User) (u User, err error) {
	response, err := client.doSysdigSecureRequest(http.MethodPut, client.GetUserUrl(uRequest.ID), uRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	u = UsersFromJSON(body)
	return
}

func (client *sysdigSecureClient) DeleteUser(id int) error {
	response, err := client.doSysdigSecureRequest(http.MethodDelete, client.GetUserUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}
	return nil
}

func (client *sysdigSecureClient) GetUsersUrl() string {
	return fmt.Sprintf("%s/api/users", client.URL)
}

func (client *sysdigSecureClient) GetUserUrl(id int) string {
	return fmt.Sprintf("%s/api/users/%d", client.URL, id)
}
