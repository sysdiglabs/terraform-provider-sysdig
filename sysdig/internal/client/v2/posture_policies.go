package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	PosturePolicyListPath   = "%s/api/cspm/v1/policy/policies/list"
	PosturePolicyCreatePath = "%s/api/cspm/v1/policy"
	PosturePolicyGetPath    = "%s/api/cspm/v1/policy/%d?include_controls=true"
)

type PosturePolicyInterface interface {
	Base
	ListPosturePolicies(ctx context.Context) ([]PosturePolicy, error)
	CreateOrUpdatePosturePolicy(ctx context.Context, p *CreatePosturePolicy) (*PosturePolicy, string, error)
	GetPosturePolicy(ctx context.Context, id int64) (*PosturePolicy, error)
}

func (client *Client) ListPosturePolicies(ctx context.Context) ([]PosturePolicy, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.getPosturePolicyURL(PosturePolicyListPath), nil)
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

func (client *Client) CreateOrUpdatePosturePolicy(ctx context.Context, p *CreatePosturePolicy) (*PosturePolicy, string, error) {
	payload, err := Marshal(p)
	if err != nil {
		return nil, "", err
	}
	response, err := client.requester.Request(ctx, http.MethodPost, client.getPosturePolicyURL(PosturePolicyCreatePath), payload)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}
	resp, err := Unmarshal[PosturePolicyResponse](response.Body)
	if err != nil {
		return nil, "", err
	}
	return &resp.Data, "", nil
}

func (client *Client) GetPosturePolicy(ctx context.Context, id int64) (*PosturePolicy, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.getPolicyUrl(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	wrapper, err := Unmarshal[PosturePolicyResponse](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.Data, nil
}

func (client *Client) getPolicyUrl(id int64) string {
	return fmt.Sprintf(PosturePolicyGetPath, client.config.url, id)
}

func (client *Client) getPosturePolicyURL(path string) string {
	return fmt.Sprintf(path, client.config.url)
}
