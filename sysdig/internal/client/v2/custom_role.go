package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var CustomRoleNotFound = errors.New("custom role not found")

const (
	CustomRolePath       = "%s/api/roles"
	UpdateCustomRolePath = "%s/api/roles/%d"
	DeleteCustomRolePath = "%s/api/roles/%d"
	GetCustomRolePath    = "%s/api/roles/%d"
)

type CustomRoleInterface interface {
	Base
	CreateCustomRole(ctx context.Context, cr *CustomRole) (*CustomRole, error)
	UpdateCustomRole(ctx context.Context, cr *CustomRole, id int) (*CustomRole, error)
	DeleteCustomRole(ctx context.Context, id int) error
	GetCustomRole(ctx context.Context, id int) (*CustomRole, error)
	GetCustomRoleByName(ctx context.Context, name string) (CustomRole, error)
}

func (client *Client) CreateCustomRole(ctx context.Context, cr *CustomRole) (*CustomRole, error) {
	payload, err := Marshal(cr)
	if err != nil {
		return nil, err
	}
	response, err := client.requester.Request(ctx, http.MethodPost, client.CreateCustomRoleURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	created, err := Unmarshal[CustomRole](response.Body)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (client *Client) UpdateCustomRole(ctx context.Context, cr *CustomRole, id int) (*CustomRole, error) {
	payload, err := Marshal(cr)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.UpdateCustomRoleURL(id), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	updated, err := Unmarshal[CustomRole](response.Body)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (client *Client) DeleteCustomRole(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.DeleteCustomRoleURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) GetCustomRole(ctx context.Context, id int) (*CustomRole, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetCustomRoleURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, CustomRoleNotFound
		}
		return nil, client.ErrorFromResponse(response)
	}

	cr, err := Unmarshal[CustomRole](response.Body)
	if err != nil {
		return nil, err
	}

	return &cr, nil
}

func (client *Client) GetCustomRoleByName(ctx context.Context, name string) (CustomRole, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetCustomRolesURL(), nil)
	if err != nil {
		return CustomRole{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return CustomRole{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[customRoleListWrapper](response.Body)

	if err != nil {
		return CustomRole{}, err
	}

	for _, customRole := range wrapper.Roles {
		if customRole.Name == name {
			return customRole, nil
		}
	}

	return CustomRole{}, fmt.Errorf("custom role with name: %s does not exist", name)

}

func (client *Client) CreateCustomRoleURL() string {
	return fmt.Sprintf(CustomRolePath, client.config.url)
}

func (client *Client) UpdateCustomRoleURL(id int) string {
	return fmt.Sprintf(UpdateCustomRolePath, client.config.url, id)
}

func (client *Client) DeleteCustomRoleURL(id int) string {
	return fmt.Sprintf(DeleteCustomRolePath, client.config.url, id)
}

func (client *Client) GetCustomRoleURL(id int) string {
	return fmt.Sprintf(GetCustomRolePath, client.config.url, id)
}

func (client *Client) GetCustomRolesURL() string {
	return fmt.Sprintf(CustomRolePath, client.config.url)
}
