package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var ErrBuiltinRoleNotFound = errors.New("builtin role not found")

const builtinRolePath = "%s/platform/v1/default-roles/%s"

type BuiltinRoleInterface interface {
	Base
	GetBuiltinRole(ctx context.Context, name string) (*BuiltinRole, error)
}

func (c *Client) GetBuiltinRole(ctx context.Context, name string) (builtinRole *BuiltinRole, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getBuiltinRoleURL(name), nil)
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
			return nil, ErrBuiltinRoleNotFound
		}
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*BuiltinRole](response.Body)
}

func (c *Client) getBuiltinRoleURL(name string) string {
	return fmt.Sprintf(builtinRolePath, c.config.url, url.PathEscape(name))
}
