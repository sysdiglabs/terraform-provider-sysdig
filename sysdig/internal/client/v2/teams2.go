package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	GetTeamsPath2 = "%s/api/teams"
	GetTeamPath2  = "%s/api/teams/%d"
)

type TeamInterface2 interface {
	Base
	GetTeamById2(ctx context.Context, id int) (t Team2, err error)
	CreateTeam2(ctx context.Context, tRequest Team2) (t Team2, err error)
	UpdateTeam2(ctx context.Context, tRequest Team2) (t Team2, err error)
	DeleteTeam2(ctx context.Context, id int) error
}

func (client *Client) GetTeamById2(ctx context.Context, id int) (Team2, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetTeamURL2(id), nil)
	if err != nil {
		return Team2{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Team2{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[teamWrapper2](response.Body)
	if err != nil {
		return Team2{}, client.ErrorFromResponse(response)
	}

	return wrapper.Team, err
}

func (client *Client) CreateTeam2(ctx context.Context, team Team2) (Team2, error) {
	var err error

	payload, err := Marshal(team)
	if err != nil {
		return Team2{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.GetTeamsURL2(), payload)
	if err != nil {
		return Team2{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return Team2{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[teamWrapper2](response.Body)
	if err != nil {
		return Team2{}, err
	}

	return wrapper.Team, nil
}

func (client *Client) UpdateTeam2(ctx context.Context, team Team2) (Team2, error) {
	var err error

	payload, err := Marshal(team)
	if err != nil {
		return Team2{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.GetTeamURL2(team.ID), payload)
	if err != nil {
		return Team2{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Team2{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[teamWrapper2](response.Body)
	if err != nil {
		return Team2{}, err
	}

	return wrapper.Team, nil
}

func (client *Client) DeleteTeam2(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.GetTeamURL2(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) GetTeamsURL2() string {
	return fmt.Sprintf(GetTeamsPath2, client.config.url)
}

func (client *Client) GetTeamURL2(id int) string {
	return fmt.Sprintf(GetTeamPath2, client.config.url, id)
}
