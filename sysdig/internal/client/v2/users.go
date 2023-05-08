package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	GetUserByUsernamePath = "%s/api/users/%s"
	GetUserPath           = "%s/api/users/%d"
	GetUsersPath          = "%s/api/users"
	CreateUserPath        = "%s/api/user/provisioning/"
	UpdateUserPath        = "%s/api/users/%d"
	DeleteUserPath        = "%s/api/users/%d"
	GetCurrentUserPath    = "%s/api/users/me"
)

type UserInterface interface {
	Base
	GetUserById(ctx context.Context, id int) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, id int) error
	GetCurrentUser(ctx context.Context) (u *User, err error)
}

func (client *Client) GetUserById(ctx context.Context, id int) (*User, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetUserUrl(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	wrapper, err := Unmarshal[userWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (client *Client) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetUserByUsernameURL(username), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	wrapper, err := Unmarshal[userWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (client *Client) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return client.GetUserByUsername(ctx, email)
}

func (client *Client) CreateUser(ctx context.Context, user *User) (*User, error) {
	payload, err := Marshal(user)
	if err != nil {
		return nil, err
	}
	response, err := client.requester.Request(ctx, http.MethodPost, client.CreateUsersURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	wrapper, err := Unmarshal[userWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (client *Client) UpdateUser(ctx context.Context, user *User) (*User, error) {
	payload, err := Marshal(user)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.UpdateUserURL(user.ID), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	wrapper, err := Unmarshal[userWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (client *Client) DeleteUser(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.DeleteUserURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) GetCurrentUser(ctx context.Context) (u *User, err error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetCurrentUserURL(), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return
	}

	wrapper, err := Unmarshal[userWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (client *Client) GetUserUrl(id int) string {
	return fmt.Sprintf(GetUserPath, client.config.url, id)
}

func (client *Client) GetUsersUrl() string {
	return fmt.Sprintf(GetUsersPath, client.config.url)
}

func (client *Client) GetUserByUsernameURL(username string) string {
	return fmt.Sprintf(GetUserByUsernamePath, client.config.url, username)
}

func (client *Client) CreateUsersURL() string {
	return fmt.Sprintf(CreateUserPath, client.config.url)
}

func (client *Client) UpdateUserURL(id int) string {
	return fmt.Sprintf(UpdateUserPath, client.config.url, id)
}

func (client *Client) DeleteUserURL(id int) string {
	return fmt.Sprintf(DeleteUserPath, client.config.url, id)
}

func (client *Client) GetCurrentUserURL() string {
	return fmt.Sprintf(GetCurrentUserPath, client.config.url)
}
