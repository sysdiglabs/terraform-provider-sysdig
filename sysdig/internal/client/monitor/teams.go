package monitor

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigMonitorClient) getUserIdbyEmail(ctx context.Context, userRoles []UserRoles) ([]UserRoles, error) {
	// Get UsersList from API
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodGet, client.getUsersListUrl(), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return nil, err
	}

	body, _ := ioutil.ReadAll(response.Body)
	// Set User Id to UserRoles struct
	usersList := UsersListFromJSON(body)
	usersMap := make(map[string]int)
	for _, u := range usersList {
		usersMap[u.Email] = u.ID
	}

	modifiedUserRoles := []UserRoles{}
	for _, userRole := range userRoles {
		ur := userRole
		id, ok := usersMap[ur.Email]
		if !ok {
			return nil, errors.New(ur.Email + " doesn't exist.")
		}
		ur.UserId = id
		modifiedUserRoles = append(modifiedUserRoles, ur)
	}

	return modifiedUserRoles, nil
}

func (client *sysdigMonitorClient) GetTeamById(ctx context.Context, id int) (t Team, err error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodGet, client.GetTeamUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	t = TeamFromJSON(body)

	return
}

func (client *sysdigMonitorClient) GetTeamByName(ctx context.Context, name string) (t Team, err error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodGet, client.GetTeamsUrl(), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	teams := TeamListFromJSON(body)

	for _, team := range teams {
		if team.Name == name && team.Products[0] == "SDC" {
			return team, nil
		}
	}
	err = fmt.Errorf("Failed to find team %s", name)
	return
}

func (client *sysdigMonitorClient) CreateTeam(ctx context.Context, tRequest Team) (t Team, err error) {
	tRequest.UserRoles, err = client.getUserIdbyEmail(ctx, tRequest.UserRoles)
	if err != nil {
		return
	}
	tRequest.Origin = "SYSDIG"

	response, err := client.doSysdigMonitorRequest(ctx, http.MethodPost, client.GetTeamsUrl(), tRequest.ToJSON())

	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errorFromResponse(response)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	t = TeamFromJSON(body)
	return
}

func (client *sysdigMonitorClient) UpdateTeam(ctx context.Context, tRequest Team) (t Team, err error) {
	tRequest.UserRoles, err = client.getUserIdbyEmail(ctx, tRequest.UserRoles)
	if err != nil {
		return
	}
	tRequest.Products = []string{"SDC"}

	response, err := client.doSysdigMonitorRequest(ctx, http.MethodPut, client.GetTeamUrl(tRequest.ID), tRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	t = TeamFromJSON(body)
	return
}

func (client *sysdigMonitorClient) DeleteTeam(ctx context.Context, id int) error {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodDelete, client.GetTeamUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return errorFromResponse(response)
	}
	return nil
}

func (client *sysdigMonitorClient) getUsersListUrl() string {
	return fmt.Sprintf("%s/api/users/light", client.URL)
}

func (client *sysdigMonitorClient) GetTeamsUrl() string {
	return fmt.Sprintf("%s/api/teams", client.URL)
}

func (client *sysdigMonitorClient) GetTeamUrl(id int) string {
	return fmt.Sprintf("%s/api/teams/%d", client.URL, id)
}
