package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	postureZonesPath = "%s/api/cspm/v1/policy/zones"
	postureZonePath  = "%s/api/cspm/v1/policy/zones/%d"
)

type PostureZoneInterface interface {
	Base
	CreateOrUpdatePostureZone(ctx context.Context, z *PostureZoneRequest) (*PostureZone, string, error)
	GetPostureZoneByID(ctx context.Context, id int) (*PostureZone, error)
	DeletePostureZone(ctx context.Context, id int) error
}

func (c *Client) CreateOrUpdatePostureZone(ctx context.Context, r *PostureZoneRequest) (zone *PostureZone, errStatus string, err error) {
	if r.ID == "" {
		r.ID = "0"
	}

	payload, err := Marshal(r)
	if err != nil {
		return nil, "", err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createPostureZoneURL(), payload)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	wrapper, err := Unmarshal[PostureZoneResponse](response.Body)
	if err != nil {
		return nil, "", err
	}

	return &wrapper.Data, "", nil
}

func (c *Client) GetPostureZoneByID(ctx context.Context, id int) (zone *PostureZone, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getPostureZoneURL(id), nil)
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
	wrapper, err := Unmarshal[PostureZoneResponse](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.Data, nil
}

func (c *Client) DeletePostureZone(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.getPostureZoneURL(id), nil)
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

func (c *Client) createPostureZoneURL() string {
	return fmt.Sprintf(postureZonesPath, c.config.url)
}

func (c *Client) getPostureZoneURL(id int) string {
	return fmt.Sprintf(postureZonePath, c.config.url, id)
}
