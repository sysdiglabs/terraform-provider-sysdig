package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	PostureControlSavePath   = "%s/api/cspm/v1/policy/controls"
	PostureControlGetPath    = "%s/api/cspm/v1/policy/controls/view/%d"
	PostureControlDeletePath = "%s/api/cspm/v1/policy/controls/%d"
)

type PostureControlInterface interface {
	Base
	CreateOrUpdatePostureControl(ctx context.Context, p *SaveControlRequest) (*PostureControl, string, error)
	GetPostureControl(ctx context.Context, id int64) (*PostureControl, error)
	DeletePostureControl(ctx context.Context, id int64) error
}

func (c *Client) CreateOrUpdatePostureControl(ctx context.Context, p *SaveControlRequest) (*PostureControl, string, error) {
	payload, err := Marshal(p)
	if err != nil {
		return nil, "", err
	}
	response, err := c.requester.Request(ctx, http.MethodPost, c.getPostureControlURL(PostureControlSavePath), payload)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}
	resp, err := Unmarshal[SaveControlResponse](response.Body)
	if err != nil {
		return nil, "", err
	}
	return &resp.Data, "", nil

}

func (c *Client) GetPostureControl(ctx context.Context, id int64) (*PostureControl, error) {
	response, err := c.requester.Request(ctx, http.MethodGet, fmt.Sprintf(PostureControlGetPath, c.config.url, id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	wrapper, err := Unmarshal[SaveControlResponse](response.Body)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

func (c *Client) DeletePostureControl(ctx context.Context, id int64) error {
	response, err := c.requester.Request(ctx, http.MethodDelete, fmt.Sprintf(PostureControlDeletePath, c.config.url, id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return c.ErrorFromResponse(response)
	}

	return nil
}
func (c *Client) getPostureControlURL(path string) string {
	return fmt.Sprintf(path, c.config.url)
}
