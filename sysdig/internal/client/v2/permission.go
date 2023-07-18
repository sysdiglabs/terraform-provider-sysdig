package v2

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const PermissionsURL = "%s/api/permissions/%s/dependencies?requestedPermissions=%s"

type PermissionInterface interface {
	Base

	GetPermissionsDependencies(ctx context.Context, product Product, permissions []string) ([]Dependency, error)
}

func (client *Client) GetPermissionsDependencies(ctx context.Context, product Product, permissions []string) ([]Dependency, error) {
	segments := map[Product]string{MonitorProduct: "monitor", SecureProduct: "secure"}
	url := fmt.Sprintf(PermissionsURL, client.config.url, segments[product], strings.Join(permissions, ","))

	return client.getPermissionsDependencies(ctx, url)
}

func (client *Client) getPermissionsDependencies(ctx context.Context, url string) ([]Dependency, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []Dependency{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []Dependency{}, client.ErrorFromResponse(response)
	}

	dependencies, err := Unmarshal[Dependencies](response.Body)

	if err != nil {
		return []Dependency{}, err
	}

	return dependencies, nil
}
