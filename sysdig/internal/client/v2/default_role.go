package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var ErrDefaultRoleNotFound = errors.New("default role not found")

const defaultRolePath = "%s/platform/v1/default-roles/%s"

type DefaultRoleInterface interface {
	Base
	GetDefaultRole(ctx context.Context, name string) (*DefaultRole, error)
}

func (c *Client) GetDefaultRole(ctx context.Context, name string) (defaultRole *DefaultRole, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getDefaultRoleURL(name), nil)
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
			return nil, ErrDefaultRoleNotFound
		}
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*DefaultRole](response.Body)
}

func (c *Client) getDefaultRoleURL(name string) string {
	return fmt.Sprintf(defaultRolePath, c.config.url, url.PathEscape(name))
}
