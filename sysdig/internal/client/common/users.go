package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigCommonClient) GetUserById(ctx context.Context, id int) (u *User, err error) {
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

	user := UserFromJSON(body)
	return &user, nil
}

func (client *sysdigCommonClient) GetUserByEmail(ctx context.Context, email string) (u *User, err error) {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodGet, client.GetUsersUrl(), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	var userList struct {
		Users []User `json:"users"`
	}

	err = json.Unmarshal(body, &userList)
	if err != nil {
		return
	}

	for _, user := range userList.Users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user not found for the given email")
}

func (client *sysdigCommonClient) CreateUser(ctx context.Context, uRequest *User) (u *User, err error) {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodPost, client.CreateUsersUrl(), uRequest.ToJSON())

	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errors.New(response.Status)
		return
	}

	user := UserFromJSON(body)
	return &user, nil
}

func (client *sysdigCommonClient) UpdateUser(ctx context.Context, uRequest *User) (u *User, err error) {
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

	user := UserFromJSON(body)
	return &user, nil
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

func (client *sysdigCommonClient) GetCurrentUser(ctx context.Context) (u *User, err error) {
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

	user := UserFromJSON(body)
	return &user, nil
}

func (client *sysdigCommonClient) CreateUsersUrl() string {
	return fmt.Sprintf("%s/api/user/provisioning/", client.URL)
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
