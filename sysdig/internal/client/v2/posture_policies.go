package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	posturePolicyListPath   = "%s/api/cspm/v1/policy/policies/list"
	posturePolicyCreatePath = "%s/api/cspm/v1/policy"
	posturePolicyGetPath    = "%s/api/cspm/v1/policy/posture/policies/%d?include_controls=true"
	posturePolicyDeletePath = "%s/api/cspm/v1/policy/policies/%d"
)

type PosturePolicyInterface interface {
	Base
	ListPosturePolicies(ctx context.Context) ([]PosturePolicy, error)
	CreateOrUpdatePosturePolicy(ctx context.Context, p *CreatePosturePolicy) (*FullPosturePolicy, string, error)
	GetPosturePolicyByID(ctx context.Context, id int64) (*FullPosturePolicy, error)
	DeletePosturePolicy(ctx context.Context, id int64) error
}

func (c *Client) ListPosturePolicies(ctx context.Context) (policies []PosturePolicy, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getPosturePolicyURL(posturePolicyListPath), nil)
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

	resp, err := Unmarshal[PostureZonePolicyListResponse](response.Body)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (c *Client) CreateOrUpdatePosturePolicy(ctx context.Context, p *CreatePosturePolicy) (policy *FullPosturePolicy, errString string, err error) {
	payload, err := Marshal(p)
	if err != nil {
		return nil, "", err
	}
	response, err := c.requester.Request(ctx, http.MethodPost, c.getPosturePolicyURL(posturePolicyCreatePath), payload)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}
	resp, err := Unmarshal[FullPosturePolicyResponse](response.Body)
	if err != nil {
		return nil, "", err
	}
	return &resp.Data, "", nil
}

func (c *Client) GetPosturePolicyByID(ctx context.Context, id int64) (policy *FullPosturePolicy, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getPosturePolicyURLForID(id), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	wrapper, err := Unmarshal[FullPosturePolicyResponse](response.Body)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

func (c *Client) DeletePosturePolicy(ctx context.Context, id int64) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deletePosturePolicyURLForID(id), nil)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return c.ErrorFromResponse(response)
	}

	return nil
}

func (c *Client) getPosturePolicyURLForID(id int64) string {
	return fmt.Sprintf(posturePolicyGetPath, c.config.url, id)
}

func (c *Client) deletePosturePolicyURLForID(id int64) string {
	return fmt.Sprintf(posturePolicyDeletePath, c.config.url, id)
}

func (c *Client) getPosturePolicyURL(path string) string {
	return fmt.Sprintf(path, c.config.url)
}
