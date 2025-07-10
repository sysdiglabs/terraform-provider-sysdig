package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	postureControlSavePath   = "%s/api/cspm/v1/policy/controls"
	postureControlGetPath    = "%s/api/cspm/v1/policy/controls/view/%d"
	postureControlDeletePath = "%s/api/cspm/v1/policy/controls/%d"
)

type PostureControlInterface interface {
	Base
	CreateOrUpdatePostureControl(ctx context.Context, p *SaveControlRequest) (*PostureControl, string, error)
	GetPostureControlByID(ctx context.Context, id int64) (*PostureControl, error)
	DeletePostureControlByID(ctx context.Context, id int64) error
}

func (c *Client) CreateOrUpdatePostureControl(ctx context.Context, p *SaveControlRequest) (control *PostureControl, status string, err error) {
	payload, err := Marshal(p)
	if err != nil {
		return nil, "", err
	}
	response, err := c.requester.Request(ctx, http.MethodPost, c.getPostureControlURL(postureControlSavePath), payload)
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
	resp, err := Unmarshal[SaveControlResponse](response.Body)
	if err != nil {
		return nil, "", err
	}
	return &resp.Data, "", nil
}

func (c *Client) GetPostureControlByID(ctx context.Context, id int64) (control *PostureControl, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, fmt.Sprintf(postureControlGetPath, c.config.url, id), nil)
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
	wrapper, err := Unmarshal[SaveControlResponse](response.Body)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

func (c *Client) DeletePostureControlByID(ctx context.Context, id int64) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, fmt.Sprintf(postureControlDeletePath, c.config.url, id), nil)
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

func (c *Client) getPostureControlURL(path string) string {
	return fmt.Sprintf(path, c.config.url)
}
