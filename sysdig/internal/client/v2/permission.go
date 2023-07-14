package v2

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const PermissionsURL = "%s/api/permissions/%s?requestedPermissions=%s&withDependencies=true"

type PermissionInterface interface {
	Base

	GetPermissionsWithDependencies(ctx context.Context, product Product, permissions []string) ([]Permission, error)
}

func (client *Client) GetPermissionsWithDependencies(ctx context.Context, product Product, permissions []string) ([]Permission, error) {
	segments := map[Product]string{MonitorProduct: "monitor", SecureProduct: "secure"}
	url := fmt.Sprintf(PermissionsURL, client.config.url, segments[product], strings.Join(permissions, ","))

	return client.getPermissionsWithDependencies(ctx, url)
}

func (client *Client) getPermissionsWithDependencies(ctx context.Context, url string) ([]Permission, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []Permission{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []Permission{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[permissionListWrapper](response.Body)

	if err != nil {
		return []Permission{}, err
	}

	return wrapper.Permissions, nil
}
