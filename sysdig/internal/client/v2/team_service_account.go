package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrTeamServiceAccountNotFound = errors.New("team service account not found")

const (
	serviceAccountsPath      = "%s/api/serviceaccounts/team"
	serviceAccountPath       = "%s/api/serviceaccounts/team/%d"
	serviceAccountDeletePath = "%s/api/serviceaccounts/team/%d/delete"
)

type TeamServiceAccountInterface interface {
	Base
	GetTeamServiceAccountByID(ctx context.Context, id int) (*TeamServiceAccount, error)
	CreateTeamServiceAccount(ctx context.Context, account *TeamServiceAccount) (*TeamServiceAccount, error)
	UpdateTeamServiceAccount(ctx context.Context, account *TeamServiceAccount, id int) (*TeamServiceAccount, error)
	DeleteTeamServiceAccount(ctx context.Context, id int) error
}

func (c *Client) GetTeamServiceAccountByID(ctx context.Context, id int) (team *TeamServiceAccount, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getTeamServiceAccountURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrTeamServiceAccountNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*TeamServiceAccount](response.Body)
}

func (c *Client) CreateTeamServiceAccount(ctx context.Context, account *TeamServiceAccount) (teamAccount *TeamServiceAccount, err error) {
	payload, err := Marshal(account)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createTeamServiceAccountURL(), payload)
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

	return Unmarshal[*TeamServiceAccount](response.Body)
}

func (c *Client) UpdateTeamServiceAccount(ctx context.Context, account *TeamServiceAccount, id int) (serviceAccount *TeamServiceAccount, err error) {
	payload, err := Marshal(account)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateTeamServiceAccountURL(id), payload)
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

	return Unmarshal[*TeamServiceAccount](response.Body)
}

func (c *Client) DeleteTeamServiceAccount(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteTeamServiceAccountURL(id), nil)
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

func (c *Client) getTeamServiceAccountURL(id int) string {
	return fmt.Sprintf(serviceAccountPath, c.config.url, id)
}

func (c *Client) createTeamServiceAccountURL() string {
	return fmt.Sprintf(serviceAccountsPath, c.config.url)
}

func (c *Client) updateTeamServiceAccountURL(id int) string {
	return fmt.Sprintf(serviceAccountPath, c.config.url, id)
}

func (c *Client) deleteTeamServiceAccountURL(id int) string {
	return fmt.Sprintf(serviceAccountDeletePath, c.config.url, id)
}
