package v2

import (
	"context"
	"fmt"
	"net/http"
)

const PosturePolicyListPath = "%s/api/cspm/v1/policy/policies/list"

type PosturePolicyInterface interface {
	Base
	ListPosturePolicies(ctx context.Context) ([]PostureZonePolicySlim, error)
}

func (client *Client) ListPosturePolicies(ctx context.Context) ([]PostureZonePolicySlim, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.getPosturePolicyListURL(), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return Unmarshal[[]PostureZonePolicySlim](response.Body)
}

func (client *Client) getPosturePolicyListURL() string {
	return fmt.Sprintf(PosturePolicyListPath, client.config.url)
}
