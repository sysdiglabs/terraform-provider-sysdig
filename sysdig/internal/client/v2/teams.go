package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	GetUsersPath      = "%s/api/users/light"
	GetTeamsPath      = "%s/api/teams"
	GetTeamPath       = "%s/api/teams/%d"
	GetTeamByNamePath = "%s/api/v2/teams/light/name/%s"
)

type TeamInterface interface {
	GetUserIDByEmail(ctx context.Context, userRoles []UserRoles) ([]UserRoles, error)
	GetTeamById(ctx context.Context, id int) (Team, error)
	GetTeamByName(ctx context.Context, name string) (Team, error)
	CreateTeam(ctx context.Context, tRequest Team) (Team, error)
	UpdateTeam(ctx context.Context, tRequest Team) (Team, error)
	DeleteTeam(ctx context.Context, id int) error
}

func (client *Client) GetUserIDByEmail(ctx context.Context, userRoles []UserRoles) ([]UserRoles, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetUsersLightURL(), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	wrapper, err := Unmarshal[usersListWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	usersMap := make(map[string]int)
	for _, u := range wrapper.UsersList {
		usersMap[u.Email] = u.ID
	}

	modifiedUserRoles := make([]UserRoles, 0)
	for _, userRole := range userRoles {
		ur := userRole
		id, ok := usersMap[ur.Email]
		if !ok {
			return nil, fmt.Errorf("email %s doesn't exist", ur.Email)
		}
		ur.UserId = id
		modifiedUserRoles = append(modifiedUserRoles, ur)
	}

	return modifiedUserRoles, nil
}

func (client *Client) GetTeamById(ctx context.Context, id int) (Team, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetTeamURL(id), nil)
	if err != nil {
		return Team{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Team{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[teamWrapper](response.Body)
	if err != nil {
		return Team{}, client.ErrorFromResponse(response)
	}

	return wrapper.Team, err
}

func (client *Client) GetTeamByName(ctx context.Context, name string) (Team, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetTeamByNameURL(name), nil)
	if err != nil {
		return Team{}, err
	}

	if response.StatusCode != http.StatusOK {
		return Team{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[teamWrapper](response.Body)
	if err != nil {
		return Team{}, client.ErrorFromResponse(response)
	}

	return wrapper.Team, err
}

func (client *Client) CreateTeam(ctx context.Context, team Team) (Team, error) {
	var err error

	team.UserRoles, err = client.GetUserIDByEmail(ctx, team.UserRoles)
	if err != nil {
		return Team{}, err
	}

	origin := "SYSDIG"
	team.Origin = origin
	payload, err := Marshal(team)
	if err != nil {
		return Team{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.GetTeamsURL(), payload)
	if err != nil {
		return Team{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return Team{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[teamWrapper](response.Body)
	if err != nil {
		return Team{}, err
	}

	return wrapper.Team, nil
}

func (client *Client) UpdateTeam(ctx context.Context, team Team) (Team, error) {
	var err error

	team.UserRoles, err = client.GetUserIDByEmail(ctx, team.UserRoles)
	if err != nil {
		return Team{}, err
	}

	payload, err := Marshal(team)
	if err != nil {
		return Team{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.GetTeamURL(team.ID), payload)
	if err != nil {
		return Team{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Team{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[teamWrapper](response.Body)
	if err != nil {
		return Team{}, err
	}

	return wrapper.Team, nil
}

func (client *Client) DeleteTeam(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.GetTeamURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) GetTeamByNameURL(name string) string {
	return fmt.Sprintf(GetTeamByNamePath, client.config.url, name)
}

func (client *Client) GetUsersLightURL() string {
	return fmt.Sprintf(GetUsersPath, client.config.url)
}

func (client *Client) GetTeamsURL() string {
	return fmt.Sprintf(GetTeamsPath, client.config.url)
}

func (client *Client) GetTeamURL(id int) string {
	return fmt.Sprintf(GetTeamPath, client.config.url, id)
}
