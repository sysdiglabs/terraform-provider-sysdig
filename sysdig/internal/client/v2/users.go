package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	getUserByUsernamePath = "%s/api/users/%s"
	getUserPath           = "%s/api/users/%d"
	createUserPath        = "%s/api/user/provisioning/"
	updateUserPath        = "%s/api/users/%d"
	deleteUserPath        = "%s/api/users/%d"
	getCurrentUserPath    = "%s/api/users/me"
)

type UserInterface interface {
	Base
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, id int) error
	GetCurrentUser(ctx context.Context) (u *User, err error)
}

func (c *Client) GetUserByID(ctx context.Context, id int) (user *User, error error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getUserURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[userWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (c *Client) GetUserByUsername(ctx context.Context, username string) (user *User, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getUserByUsernameURL(username), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[userWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return c.GetUserByUsername(ctx, email)
}

func (c *Client) CreateUser(ctx context.Context, user *User) (createdUser *User, err error) {
	payload, err := Marshal(user)
	if err != nil {
		return nil, err
	}
	response, err := c.requester.Request(ctx, http.MethodPost, c.createUsersURL(), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[userWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (c *Client) UpdateUser(ctx context.Context, user *User) (updated *User, err error) {
	payload, err := Marshal(user)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateUserURL(user.ID), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[userWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (c *Client) DeleteUser(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteUserURL(id), nil)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return c.ErrorFromResponse(response)
	}

	return nil
}

func (c *Client) GetCurrentUser(ctx context.Context) (u *User, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getCurrentUserURL(), nil)
	if err != nil {
		return
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		err = c.ErrorFromResponse(response)
		return
	}

	wrapper, err := Unmarshal[userWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (c *Client) getUserURL(id int) string {
	return fmt.Sprintf(getUserPath, c.config.url, id)
}

func (c *Client) getUserByUsernameURL(username string) string {
	return fmt.Sprintf(getUserByUsernamePath, c.config.url, username)
}

func (c *Client) createUsersURL() string {
	return fmt.Sprintf(createUserPath, c.config.url)
}

func (c *Client) updateUserURL(id int) string {
	return fmt.Sprintf(updateUserPath, c.config.url, id)
}

func (c *Client) deleteUserURL(id int) string {
	return fmt.Sprintf(deleteUserPath, c.config.url, id)
}

func (c *Client) getCurrentUserURL() string {
	return fmt.Sprintf(getCurrentUserPath, c.config.url)
}
