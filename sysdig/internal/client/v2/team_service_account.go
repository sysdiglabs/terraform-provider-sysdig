package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var TeamServiceAccountNotFound = errors.New("team service account not found")

const (
	ServiceAccountsPath      = "%s/api/serviceaccounts/team"
	ServiceAccountPath       = "%s/api/serviceaccounts/team/%d"
	ServiceAccountDeletePath = "%s/api/serviceaccounts/team/%d/delete"
)

type TeamServiceAccountInterface interface {
	Base
	GetTeamServiceAccountByID(ctx context.Context, id int) (*TeamServiceAccount, error)
	CreateTeamServiceAccount(ctx context.Context, account *TeamServiceAccount) (*TeamServiceAccount, error)
	UpdateTeamServiceAccount(ctx context.Context, account *TeamServiceAccount, id int) (*TeamServiceAccount, error)
	DeleteTeamServiceAccount(ctx context.Context, id int) error
}

func (client *Client) GetTeamServiceAccountByID(ctx context.Context, id int) (*TeamServiceAccount, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetTeamServiceAccountURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, TeamServiceAccountNotFound
		}
		return nil, client.ErrorFromResponse(response)
	}

	teamServiceAccount, err := Unmarshal[TeamServiceAccount](response.Body)
	if err != nil {
		return nil, err
	}
	return &teamServiceAccount, nil
}

func (client *Client) CreateTeamServiceAccount(ctx context.Context, account *TeamServiceAccount) (*TeamServiceAccount, error) {
	payload, err := Marshal(account)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.CreateTeamServiceAccountURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	created, err := Unmarshal[TeamServiceAccount](response.Body)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (client *Client) UpdateTeamServiceAccount(ctx context.Context, account *TeamServiceAccount, id int) (*TeamServiceAccount, error) {
	payload, err := Marshal(account)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.UpdateTeamServiceAccountURL(id), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	updated, err := Unmarshal[TeamServiceAccount](response.Body)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (client *Client) DeleteTeamServiceAccount(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.DeleteTeamServiceAccountURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) GetTeamServiceAccountURL(id int) string {
	return fmt.Sprintf(ServiceAccountPath, client.config.url, id)
}

func (client *Client) CreateTeamServiceAccountURL() string {
	return fmt.Sprintf(ServiceAccountsPath, client.config.url)
}

func (client *Client) UpdateTeamServiceAccountURL(id int) string {
	return fmt.Sprintf(ServiceAccountPath, client.config.url, id)
}

func (client *Client) DeleteTeamServiceAccountURL(id int) string {
	return fmt.Sprintf(ServiceAccountDeletePath, client.config.url, id)
}
