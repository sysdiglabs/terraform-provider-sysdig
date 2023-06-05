package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var GroupMappingConfigNotFound = errors.New("group mapping configuration not found")

const (
	GroupMappingConfigPath = "%s/api/groupmappings/settings"
)

type GroupMappingConfigInterface interface {
	Base
	UpdateGroupMappingConfig(ctx context.Context, gmc *GroupMappingConfig) (*GroupMappingConfig, error)
	GetGroupMappingConfig(ctx context.Context) (*GroupMappingConfig, error)
}

func (client *Client) UpdateGroupMappingConfig(ctx context.Context, gmc *GroupMappingConfig) (*GroupMappingConfig, error) {
	payload, err := Marshal(gmc)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.UpdateGroupMappingConfigURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	updated, err := Unmarshal[GroupMappingConfig](response.Body)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (client *Client) GetGroupMappingConfig(ctx context.Context) (*GroupMappingConfig, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetGroupMappingConfigURL(), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, GroupMappingConfigNotFound
		}
		return nil, client.ErrorFromResponse(response)
	}

	gmc, err := Unmarshal[GroupMappingConfig](response.Body)
	if err != nil {
		return nil, err
	}

	return &gmc, nil
}

func (client *Client) UpdateGroupMappingConfigURL() string {
	return fmt.Sprintf(GroupMappingConfigPath, client.config.url)
}

func (client *Client) GetGroupMappingConfigURL() string {
	return fmt.Sprintf(GroupMappingConfigPath, client.config.url)
}
