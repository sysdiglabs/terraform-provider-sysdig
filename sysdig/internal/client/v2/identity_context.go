package v2

import (
	"context"
	"fmt"
	"net/http"
)

const GetIdentityContextPath = "%s/api/identity/context"

type IdentityContextInterface interface {
	GetIdentityContext(ctx context.Context) (*IdentityContext, error)
}

func (client *Client) GetIdentityContext(ctx context.Context) (*IdentityContext, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetIdentityContextURL(), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return Unmarshal[*IdentityContext](response.Body)
}

func (client *Client) GetIdentityContextURL() string {
	return fmt.Sprintf(GetIdentityContextPath, client.config.url)
}
