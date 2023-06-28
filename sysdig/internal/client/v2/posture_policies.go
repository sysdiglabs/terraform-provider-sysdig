package v2

import (
	"context"
	"fmt"
	"net/http"
)

const PosturePolicyListPath = "%s/api/cspm/v1/policy/policies/list"

type PosturePolicyInterface interface {
	Base
	ListPosturePolicies(ctx context.Context) ([]PosturePolicy, error)
}

func (client *Client) ListPosturePolicies(ctx context.Context) ([]PosturePolicy, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.getPosturePolicyListURL(), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	resp, err := Unmarshal[PostureZonePolicyListResponse](response.Body)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (client *Client) getPosturePolicyListURL() string {
	return fmt.Sprintf(PosturePolicyListPath, client.config.url)
}
