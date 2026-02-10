package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrSSOGroupMappingNotFound = errors.New("SSO group mapping not found")

const (
	createSSOGroupMappingPath = "%s/platform/v1/group-mappings"
	getSSOGroupMappingPath    = "%s/platform/v1/group-mappings/%d"
	updateSSOGroupMappingPath = "%s/platform/v1/group-mappings/%d"
	deleteSSOGroupMappingPath = "%s/platform/v1/group-mappings/%d"
)

type SSOGroupMappingInterface interface {
	Base
	CreateSSOGroupMapping(ctx context.Context, gm *SSOGroupMapping) (*SSOGroupMapping, error)
	GetSSOGroupMapping(ctx context.Context, id int) (*SSOGroupMapping, error)
	UpdateSSOGroupMapping(ctx context.Context, id int, gm *SSOGroupMapping) (*SSOGroupMapping, error)
	DeleteSSOGroupMapping(ctx context.Context, id int) error
}

func (c *Client) CreateSSOGroupMapping(ctx context.Context, gm *SSOGroupMapping) (result *SSOGroupMapping, err error) {
	payload, err := Marshal(gm)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createSSOGroupMappingURL(), payload)
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

	return Unmarshal[*SSOGroupMapping](response.Body)
}

func (c *Client) GetSSOGroupMapping(ctx context.Context, id int) (result *SSOGroupMapping, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getSSOGroupMappingURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrSSOGroupMappingNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*SSOGroupMapping](response.Body)
}

func (c *Client) UpdateSSOGroupMapping(ctx context.Context, id int, gm *SSOGroupMapping) (result *SSOGroupMapping, err error) {
	payload, err := Marshal(gm)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateSSOGroupMappingURL(id), payload)
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

	return Unmarshal[*SSOGroupMapping](response.Body)
}

func (c *Client) DeleteSSOGroupMapping(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteSSOGroupMappingURL(id), nil)
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

func (c *Client) createSSOGroupMappingURL() string {
	return fmt.Sprintf(createSSOGroupMappingPath, c.config.url)
}

func (c *Client) getSSOGroupMappingURL(id int) string {
	return fmt.Sprintf(getSSOGroupMappingPath, c.config.url, id)
}

func (c *Client) updateSSOGroupMappingURL(id int) string {
	return fmt.Sprintf(updateSSOGroupMappingPath, c.config.url, id)
}

func (c *Client) deleteSSOGroupMappingURL(id int) string {
	return fmt.Sprintf(deleteSSOGroupMappingPath, c.config.url, id)
}
