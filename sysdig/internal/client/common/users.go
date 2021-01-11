package common

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigCommonClient) GetUserById(ctx context.Context, id int) (u User, err error) {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodGet, client.GetUserUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	u = UserFromJSON(body)

	return
}

func (client *sysdigCommonClient) CreateUser(ctx context.Context, uRequest User) (u User, err error) {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodPost, client.GetUsersUrl(), uRequest.ToJSON())

	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errors.New(response.Status)
		return
	}

	u = UserFromJSON(body)
	return
}

func (client *sysdigCommonClient) UpdateUser(ctx context.Context, uRequest User) (u User, err error) {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodPut, client.GetUserUrl(uRequest.ID), uRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	u = UserFromJSON(body)
	return
}

func (client *sysdigCommonClient) DeleteUser(ctx context.Context, id int) error {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodDelete, client.GetUserUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}
	return nil
}

func (client *sysdigCommonClient) GetCurrentUser(ctx context.Context) (u User, err error) {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodGet, client.GetCurrentUserUrl(), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}
	body, _ := ioutil.ReadAll(response.Body)

	u = UserFromJSON(body)
	return
}

func (client *sysdigCommonClient) GetUsersUrl() string {
	return fmt.Sprintf("%s/api/users", client.URL)
}

func (client *sysdigCommonClient) GetUserUrl(id int) string {
	return fmt.Sprintf("%s/api/users/%d", client.URL, id)
}

func (client *sysdigCommonClient) GetCurrentUserUrl() string {
	return fmt.Sprintf("%s/api/users/me", client.URL)
}
