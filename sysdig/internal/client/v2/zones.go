package v2

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	platformZonesPath = "%s/platform/v1/zones"
	platformZonePath  = "%s/platform/v1/zones/%d"
)

type ZoneInterface interface {
	Base
	GetZones(ctx context.Context, name string) ([]Zone, error)
	GetZoneByID(ctx context.Context, id int) (*Zone, error)
	CreateZone(ctx context.Context, zone *ZoneRequest) (*Zone, error)
	UpdateZone(ctx context.Context, zone *ZoneRequest) (*Zone, error)
	DeleteZone(ctx context.Context, id int) error
}

func (c *Client) GetZones(ctx context.Context, name string) (zones []Zone, err error) {
	zonesURL := c.getZonesURL()
	zonesURL = fmt.Sprintf("%s?filter=name:%s", zonesURL, url.QueryEscape(name))

	response, err := c.requester.Request(ctx, http.MethodGet, zonesURL, nil)
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
	wrapper, err := Unmarshal[ZonesWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return wrapper.Zones, nil
}

func (c *Client) GetZoneByID(ctx context.Context, id int) (zone *Zone, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getZoneURL(id), nil)
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

	return Unmarshal[*Zone](response.Body)
}

func (c *Client) CreateZone(ctx context.Context, zone *ZoneRequest) (createdZone *Zone, err error) {
	payload, err := Marshal(zone)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.getZonesURL(), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*Zone](response.Body)
}

func (c *Client) UpdateZone(ctx context.Context, zone *ZoneRequest) (updatedZone *Zone, err error) {
	payload, err := Marshal(zone)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.getZoneURL(zone.ID), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*Zone](response.Body)
}

func (c *Client) DeleteZone(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.getZoneURL(id), nil)
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

func (c *Client) getZonesURL() string {
	return fmt.Sprintf(platformZonesPath, c.config.url)
}

func (c *Client) getZoneURL(id int) string {
	return fmt.Sprintf(platformZonePath, c.config.url, id)
}
