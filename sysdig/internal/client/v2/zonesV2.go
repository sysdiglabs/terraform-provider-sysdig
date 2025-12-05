package v2

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	platformZonesPathV2 = "%s/platform/v2/zones"
	platformZonePathV2  = "%s/platform/v2/zones/%d"
)

type ZoneV2Interface interface {
	Base
	GetZonesV2(ctx context.Context, name string) ([]ZoneV2, error)
	GetZoneV2(ctx context.Context, id int) (*ZoneV2, error)
	CreateZoneV2(ctx context.Context, zone *ZoneV2) (*ZoneV2, error)
	UpdateZoneV2(ctx context.Context, zone *ZoneV2) (*ZoneV2, error)
	DeleteZoneV2(ctx context.Context, id int) error
}

func (c *Client) GetZonesV2(ctx context.Context, name string) (zones []ZoneV2, err error) {
	zonesURL := c.getZonesV2URL()
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
		return nil, c.APIErrorFromResponse(response)
	}
	wrapper, err := Unmarshal[ZonesV2Wrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return wrapper.Zones, nil
}

func (c *Client) GetZoneV2(ctx context.Context, id int) (zone *ZoneV2, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getZoneV2URL(id), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, c.APIErrorFromResponse(response)
	}

	return Unmarshal[*ZoneV2](response.Body)
}

func (c *Client) CreateZoneV2(ctx context.Context, zone *ZoneV2) (createdZone *ZoneV2, err error) {
	payload, err := Marshal(zone)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.getZonesV2URL(), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, c.APIErrorFromResponse(response)
	}

	return Unmarshal[*ZoneV2](response.Body)
}

func (c *Client) UpdateZoneV2(ctx context.Context, zone *ZoneV2) (updatedZone *ZoneV2, err error) {
	payload, err := Marshal(zone)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.getZoneV2URL(zone.ID), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, c.APIErrorFromResponse(response)
	}

	return Unmarshal[*ZoneV2](response.Body)
}

func (c *Client) DeleteZoneV2(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.getZoneV2URL(id), nil)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return c.APIErrorFromResponse(response)
	}

	return nil
}

func (c *Client) getZonesV2URL() string {
	return fmt.Sprintf(platformZonesPathV2, c.config.url)
}

func (c *Client) getZoneV2URL(id int) string {
	return fmt.Sprintf(platformZonePathV2, c.config.url, id)
}
