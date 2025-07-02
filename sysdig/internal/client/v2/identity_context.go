package v2

import (
	"context"
	"fmt"
	"net/http"
)

const getIdentityContextPath = "%s/api/identity/context"

type IdentityContextInterface interface {
	GetIdentityContext(ctx context.Context) (*IdentityContext, error)
}

func (c *Client) GetIdentityContext(ctx context.Context) (idx *IdentityContext, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getIdentityContextURL(), nil)
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

	return Unmarshal[*IdentityContext](response.Body)
}

func (c *Client) getIdentityContextURL() string {
	return fmt.Sprintf(getIdentityContextPath, c.config.url)
}
