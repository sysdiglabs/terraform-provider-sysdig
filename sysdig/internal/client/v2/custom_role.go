package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrCustomRoleNotFound = errors.New("custom role not found")

const (
	customRolesPath = "%s/api/roles"
	customRolePath  = "%s/api/roles/%d"
)

type CustomRoleInterface interface {
	Base
	CreateCustomRole(ctx context.Context, cr *CustomRole) (*CustomRole, error)
	UpdateCustomRole(ctx context.Context, cr *CustomRole, id int) (*CustomRole, error)
	DeleteCustomRole(ctx context.Context, id int) error
	GetCustomRoleByID(ctx context.Context, id int) (*CustomRole, error)
	GetCustomRoleByName(ctx context.Context, name string) (*CustomRole, error)
}

func (c *Client) CreateCustomRole(ctx context.Context, cr *CustomRole) (customRole *CustomRole, err error) {
	payload, err := Marshal(cr)
	if err != nil {
		return nil, err
	}
	response, err := c.requester.Request(ctx, http.MethodPost, c.createCustomRoleURL(), payload)
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

	return Unmarshal[*CustomRole](response.Body)
}

func (c *Client) UpdateCustomRole(ctx context.Context, cr *CustomRole, id int) (customRole *CustomRole, err error) {
	payload, err := Marshal(cr)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateCustomRoleURL(id), payload)
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

	return Unmarshal[*CustomRole](response.Body)
}

func (c *Client) DeleteCustomRole(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteCustomRoleURL(id), nil)
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

func (c *Client) GetCustomRoleByID(ctx context.Context, id int) (customRole *CustomRole, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getCustomRoleURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, ErrCustomRoleNotFound
		}
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*CustomRole](response.Body)
}

func (c *Client) GetCustomRoleByName(ctx context.Context, name string) (customRole *CustomRole, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getCustomRolesURL(), nil)
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

	wrapper, err := Unmarshal[customRoleListWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	for _, customRole := range wrapper.Roles {
		if customRole.Name == name {
			return &customRole, nil
		}
	}

	return nil, fmt.Errorf("custom role with name, %s does not exist: %w", name, ErrCustomRoleNotFound)
}

func (c *Client) createCustomRoleURL() string {
	return fmt.Sprintf(customRolesPath, c.config.url)
}

func (c *Client) updateCustomRoleURL(id int) string {
	return fmt.Sprintf(customRolePath, c.config.url, id)
}

func (c *Client) deleteCustomRoleURL(id int) string {
	return fmt.Sprintf(customRolePath, c.config.url, id)
}

func (c *Client) getCustomRoleURL(id int) string {
	return fmt.Sprintf(customRolePath, c.config.url, id)
}

func (c *Client) getCustomRolesURL() string {
	return fmt.Sprintf(customRolesPath, c.config.url)
}
