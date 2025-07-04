package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrGroupMappingConfigNotFound = errors.New("group mapping configuration not found")

const (
	groupMappingConfigPath = "%s/api/groupmappings/settings"
)

type GroupMappingConfigInterface interface {
	Base
	UpdateGroupMappingConfig(ctx context.Context, gmc *GroupMappingConfig) (*GroupMappingConfig, error)
	GetGroupMappingConfig(ctx context.Context) (*GroupMappingConfig, error)
}

func (c *Client) UpdateGroupMappingConfig(ctx context.Context, gmc *GroupMappingConfig) (mapping *GroupMappingConfig, err error) {
	payload, err := Marshal(gmc)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateGroupMappingConfigURL(), payload)
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

	return Unmarshal[*GroupMappingConfig](response.Body)
}

func (c *Client) GetGroupMappingConfig(ctx context.Context) (mapping *GroupMappingConfig, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getGroupMappingConfigURL(), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrGroupMappingConfigNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*GroupMappingConfig](response.Body)
}

func (c *Client) updateGroupMappingConfigURL() string {
	return fmt.Sprintf(groupMappingConfigPath, c.config.url)
}

func (c *Client) getGroupMappingConfigURL() string {
	return fmt.Sprintf(groupMappingConfigPath, c.config.url)
}
