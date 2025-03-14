package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	PostureZonesPath = "%s/api/cspm/v1/policy/zones"
	PostureZonePath  = "%s/api/cspm/v1/policy/zones/%d"
)

type PostureZoneInterface interface {
	Base
	CreateOrUpdatePostureZone(ctx context.Context, z *PostureZoneRequest) (*PostureZone, string, error)
	GetPostureZone(ctx context.Context, id int) (*PostureZone, error)
	DeletePostureZone(ctx context.Context, id int) error
}

func (client *Client) CreateOrUpdatePostureZone(ctx context.Context, r *PostureZoneRequest) (*PostureZone, string, error) {
	if r.ID == "" {
		r.ID = "0"
	}

	payload, err := Marshal(r)
	if err != nil {
		return nil, "", err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.createPostureZoneURL(), payload)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
		errStatus, err := client.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	wrapper, err := Unmarshal[PostureZoneResponse](response.Body)
	if err != nil {
		return nil, "", err
	}

	return &wrapper.Data, "", nil
}

func (client *Client) GetPostureZone(ctx context.Context, id int) (*PostureZone, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.getPostureZoneURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	wrapper, err := Unmarshal[PostureZoneResponse](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.Data, nil
}

func (client *Client) DeletePostureZone(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.getPostureZoneURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) createPostureZoneURL() string {
	return fmt.Sprintf(PostureZonesPath, client.config.url)
}

func (client *Client) getPostureZoneURL(id int) string {
	return fmt.Sprintf(PostureZonePath, client.config.url, id)
}
