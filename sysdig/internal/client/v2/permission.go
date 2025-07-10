package v2

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const permissionsURL = "%s/api/permissions/%s/dependencies?requestedPermissions=%s"

type CustomRolePermissionInterface interface {
	Base

	GetPermissionsDependencies(ctx context.Context, product Product, permissions []string) ([]Dependency, error)
}

func (c *Client) GetPermissionsDependencies(ctx context.Context, product Product, permissions []string) ([]Dependency, error) {
	segments := map[Product]string{MonitorProduct: "monitor", SecureProduct: "secure"}
	url := fmt.Sprintf(permissionsURL, c.config.url, segments[product], strings.Join(permissions, ","))

	return c.getPermissionsDependencies(ctx, url)
}

func (c *Client) getPermissionsDependencies(ctx context.Context, url string) (dependencies []Dependency, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, url, nil)
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

	return Unmarshal[Dependencies](response.Body)
}
