package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	ZonesPath = "%s/api/cspm/v1/policy/zones"
	ZonePath  = "%s/api/cspm/v1/policy/zones/%d"
)

type PostureZoneInterface interface {
	Base
	CreateOrUpdatePostureZone(ctx context.Context, z *PostureZoneRequest) (*PostureZone, error)
	GetPostureZone(ctx context.Context, id int) (*PostureZone, error)
	DeletePostureZone(ctx context.Context, id int) error
}

func (client *Client) CreateOrUpdatePostureZone(ctx context.Context, r *PostureZoneRequest) (*PostureZone, error) {
	payload, err := Marshal(r)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.createZoneURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	zone, err := Unmarshal[PostureZone](response.Body)
	if err != nil {
		return nil, err
	}

	return &zone, nil
}

func (client *Client) GetPostureZone(ctx context.Context, id int) (*PostureZone, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.getZoneURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	zone, err := Unmarshal[PostureZone](response.Body)
	if err != nil {
		return nil, err
	}

	return &zone, nil
}

func (client *Client) DeletePostureZone(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.getZoneURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) createZoneURL() string {
	return fmt.Sprintf(ZonesPath, client.config.url)
}

func (client *Client) getZoneURL(id int) string {
	return fmt.Sprintf(ZonePath, client.config.url, id)
}
