package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (client *sysdigCommonClient) GetUserById(ctx context.Context, id int) (u *User, err error) {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodGet, client.GetUserUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	user := UserFromJSON(body)
	return &user, nil
}

func (client *sysdigCommonClient) GetUserByEmail(ctx context.Context, email string) (u *User, err error) {
	urlString := fmt.Sprintf("%s%s", client.GetUsersUrl(), url.PathEscape(email))

	response, err := client.doSysdigCommonRequest(ctx, http.MethodGet, urlString, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return nil, err
	}

	var user User
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (client *sysdigCommonClient) CreateUser(ctx context.Context, uRequest *User) (u *User, err error) {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodPost, client.CreateUsersUrl(), uRequest.ToJSON())

	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	user := UserFromJSON(body)
	return &user, nil
}

func (client *sysdigCommonClient) UpdateUser(ctx context.Context, uRequest *User) (u *User, err error) {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodPut, client.GetUserUrl(uRequest.ID), uRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	user := UserFromJSON(body)
	return &user, nil
}

func (client *sysdigCommonClient) DeleteUser(ctx context.Context, id int) error {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodDelete, client.GetUserUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return errorFromResponse(response)
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
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	user := UserFromJSON(body)
	return &user, nil
}

func (client *sysdigCommonClient) CreateUsersUrl() string {
	return fmt.Sprintf("%s/api/user/provisioning/", client.URL)
}

func (client *sysdigCommonClient) GetUsersUrl() string {
	return fmt.Sprintf("%s/api/users/", client.URL)
}

func (client *sysdigCommonClient) GetUserUrl(id int) string {
	return fmt.Sprintf("%s/api/users/%d", client.URL, id)
}

func (client *sysdigCommonClient) GetCurrentUserUrl() string {
	return fmt.Sprintf("%s/api/users/me", client.URL)
}
