package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var GroupMappingConfigNotFound = errors.New("group mapping configuration not found")

const (
	CreateGroupMappingConfigPath = "%s/api/groupmappings/settings"
	UpdateGroupMappingConfigPath = "%s/api/groupmappings/settings"
	GetGroupMappingConfigPath    = "%s/api/groupmappings/settings"
)

type GroupMappingConfigInterface interface {
	Base
	CreateGroupMappingConfig(ctx context.Context, gmc *GroupMappingConfig) (*GroupMappingConfig, error)
	UpdateGroupMappingConfig(ctx context.Context, gmc *GroupMappingConfig) (*GroupMappingConfig, error)
	GetGroupMappingConfig(ctx context.Context) (*GroupMappingConfig, error)
}

func (client *Client) CreateGroupMappingConfig(ctx context.Context, gmc *GroupMappingConfig) (*GroupMappingConfig, error) {
	payload, err := Marshal(gmc)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.CreateGroupMappingConfigURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	created, err := Unmarshal[GroupMappingConfig](response.Body)
	if err != nil {
		return nil, err
	}

	return &created, nil
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

func (client *Client) CreateGroupMappingConfigURL() string {
	return fmt.Sprintf(CreateGroupMappingConfigPath, client.config.url)
}

func (client *Client) UpdateGroupMappingConfigURL() string {
	return fmt.Sprintf(UpdateGroupMappingConfigPath, client.config.url)
}

func (client *Client) GetGroupMappingConfigURL() string {
	return fmt.Sprintf(GetGroupMappingConfigPath, client.config.url)
}
