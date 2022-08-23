package secure

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (client *sysdigSecureClient) getUserIdbyEmail(ctx context.Context, userRoles []UserRoles) ([]UserRoles, error) {
	// Get UsersList from API
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.getUsersListUrl(), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return nil, err
	}

	body, _ := io.ReadAll(response.Body)
	// Set User Id to UserRoles struct
	usersList := UsersListFromJSON(body)
	usersMap := make(map[string]int)
	for _, u := range usersList {
		usersMap[u.Email] = u.ID
	}

	var modifiedUserRoles []UserRoles

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

func (client *sysdigSecureClient) GetTeamById(ctx context.Context, id int) (t Team, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.GetTeamUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	t = TeamFromJSON(body)

	return
}

func (client *sysdigSecureClient) CreateTeam(ctx context.Context, tRequest Team) (t Team, err error) {
	tRequest.UserRoles, err = client.getUserIdbyEmail(ctx, tRequest.UserRoles)
	if err != nil {
		return
	}
	tRequest.Products = []string{"SDS"}

	response, err := client.doSysdigSecureRequest(ctx, http.MethodPost, client.GetTeamsUrl(), tRequest.ToJSON())

	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	t = TeamFromJSON(body)
	return
}

func (client *sysdigSecureClient) UpdateTeam(ctx context.Context, tRequest Team) (t Team, err error) {
	tRequest.UserRoles, err = client.getUserIdbyEmail(ctx, tRequest.UserRoles)
	if err != nil {
		return
	}
	tRequest.Products = []string{"SDS"}

	response, err := client.doSysdigSecureRequest(ctx, http.MethodPut, client.GetTeamUrl(tRequest.ID), tRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	t = TeamFromJSON(body)
	return
}

func (client *sysdigSecureClient) DeleteTeam(ctx context.Context, id int) error {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodDelete, client.GetTeamUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}
	return nil
}

func (client *sysdigSecureClient) getUsersListUrl() string {
	return fmt.Sprintf("%s/api/users/light", client.URL)
}

func (client *sysdigSecureClient) GetTeamsUrl() string {
	return fmt.Sprintf("%s/api/teams", client.URL)
}

func (client *sysdigSecureClient) GetTeamUrl(id int) string {
	return fmt.Sprintf("%s/api/teams/%d", client.URL, id)
}
