package v2

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	PlatformZonesPath = "%s/platform/v1/zones"
	PlatformZonePath  = "%s/platform/v1/zones/%d"
)

type ZoneInterface interface {
	Base
	GetZones(ctx context.Context, name string) ([]Zone, error)
	GetZoneById(ctx context.Context, id int) (*Zone, error)
	CreateZone(ctx context.Context, zone *ZoneRequest) (*Zone, error)
	UpdateZone(ctx context.Context, zone *ZoneRequest) (*Zone, error)
	DeleteZone(ctx context.Context, id int) error
}

func (client *Client) GetZones(ctx context.Context, name string) ([]Zone, error) {
	zonesURL := client.getZonesURL()
	zonesURL = fmt.Sprintf("%s?filter=name:%s", zonesURL, url.QueryEscape(name))

	response, err := client.requester.Request(ctx, http.MethodGet, zonesURL, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	wrapper, err := Unmarshal[ZonesWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return wrapper.Zones, nil
}

func (client *Client) GetZoneById(ctx context.Context, id int) (*Zone, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.getZoneURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	zone, err := Unmarshal[Zone](response.Body)
	if err != nil {
		return nil, err
	}

	return &zone, nil
}

func (client *Client) CreateZone(ctx context.Context, zone *ZoneRequest) (*Zone, error) {
	payload, err := Marshal(zone)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.getZonesURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, client.ErrorFromResponse(response)
	}

	createdZone, err := Unmarshal[Zone](response.Body)
	if err != nil {
		return nil, err
	}

	return &createdZone, nil
}

func (client *Client) UpdateZone(ctx context.Context, zone *ZoneRequest) (*Zone, error) {
	payload, err := Marshal(zone)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.getZoneURL(zone.ID), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, client.ErrorFromResponse(response)
	}

	updatedZone, err := Unmarshal[Zone](response.Body)
	if err != nil {
		return nil, err
	}

	return &updatedZone, nil
}

func (client *Client) DeleteZone(ctx context.Context, id int) error {
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

func (client *Client) getZonesURL() string {
	return fmt.Sprintf(PlatformZonesPath, client.config.url)
}

func (client *Client) getZoneURL(id int) string {
	return fmt.Sprintf(PlatformZonePath, client.config.url, id)
}
