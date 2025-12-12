package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	getUsersLightPath = "%s/api/users/light"
	getTeamsPath      = "%s/api/teams"
	getTeamPath       = "%s/api/teams/%d"
)

type TeamInterface interface {
	Base
	GetUserIDByEmail(ctx context.Context, userRoles []UserRoles) ([]UserRoles, error)
	GetTeamByID(ctx context.Context, id int) (t Team, statusCode int, err error)
	CreateTeam(ctx context.Context, tRequest Team) (t Team, err error)
	UpdateTeam(ctx context.Context, tRequest Team) (t Team, err error)
	DeleteTeam(ctx context.Context, id int) error
	ListTeams(ctx context.Context) ([]Team, error)
}

type teamsWrapper struct {
	Teams []Team `json:"teams"`
}

func (c *Client) GetUserIDByEmail(ctx context.Context, userRoles []UserRoles) (modifiedUserRoles []UserRoles, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getUsersLightURL(), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		err = c.ErrorFromResponse(response)
		return nil, err
	}

	wrapper, err := Unmarshal[usersWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	usersMap := make(map[string]int)
	for _, u := range wrapper.Users {
		usersMap[u.Email] = u.ID
	}

	modifiedUserRoles = make([]UserRoles, 0)
	for _, userRole := range userRoles {
		ur := userRole
		id, ok := usersMap[ur.Email]
		if !ok {
			return nil, fmt.Errorf("email %s doesn't exist", ur.Email)
		}
		ur.UserID = id
		modifiedUserRoles = append(modifiedUserRoles, ur)
	}

	return modifiedUserRoles, nil
}

func (c *Client) GetTeamByID(ctx context.Context, id int) (team Team, statusCode int, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getTeamURL(id), nil)
	if err != nil {
		return Team{}, 0, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Team{}, response.StatusCode, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[teamWrapper](response.Body)
	if err != nil {
		return Team{}, response.StatusCode, c.ErrorFromResponse(response)
	}

	return wrapper.Team, response.StatusCode, err
}

func (c *Client) CreateTeam(ctx context.Context, team Team) (createdTeam Team, err error) {
	team.UserRoles, err = c.GetUserIDByEmail(ctx, team.UserRoles)
	if err != nil {
		return Team{}, err
	}

	origin := "SYSDIG"
	team.Origin = origin
	payload, err := Marshal(team)
	if err != nil {
		return Team{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.getTeamsURL(), payload)
	if err != nil {
		return Team{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return Team{}, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[teamWrapper](response.Body)
	if err != nil {
		return Team{}, err
	}

	return wrapper.Team, nil
}

func (c *Client) UpdateTeam(ctx context.Context, team Team) (updatedTeam Team, err error) {
	team.UserRoles, err = c.GetUserIDByEmail(ctx, team.UserRoles)
	if err != nil {
		return Team{}, err
	}

	payload, err := Marshal(team)
	if err != nil {
		return Team{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.getTeamURL(team.ID), payload)
	if err != nil {
		return Team{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Team{}, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[teamWrapper](response.Body)
	if err != nil {
		return Team{}, err
	}

	return wrapper.Team, nil
}

func (c *Client) DeleteTeam(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.getTeamURL(id), nil)
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

func (c *Client) ListTeams(ctx context.Context) (teams []Team, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getTeamsURL(), nil)
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

	wrapper, err := Unmarshal[teamsWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return wrapper.Teams, nil
}

func (c *Client) getUsersLightURL() string {
	return fmt.Sprintf(getUsersLightPath, c.config.url)
}

func (c *Client) getTeamsURL() string {
	return fmt.Sprintf(getTeamsPath, c.config.url)
}

func (c *Client) getTeamURL(id int) string {
	return fmt.Sprintf(getTeamPath, c.config.url, id)
}
