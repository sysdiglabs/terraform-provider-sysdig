package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var GroupMappingNotFound = errors.New("group mapping not found")

const (
	CreateGroupMappingPath = "%s/api/groupmappings"
	UpdateGroupMappingPath = "%s/api/groupmappings/%d"
	DeleteGroupMappingPath = "%s/api/groupmappings/%d"
	GetGroupMappingPath    = "%s/api/groupmappings/%d"
)

type GroupMappingInterface interface {
	Base
	CreateGroupMapping(ctx context.Context, gm *GroupMapping) (*GroupMapping, error)
	UpdateGroupMapping(ctx context.Context, gm *GroupMapping, id int) (*GroupMapping, error)
	DeleteGroupMapping(ctx context.Context, id int) error
	GetGroupMapping(ctx context.Context, id int) (*GroupMapping, error)
}

func (client *Client) CreateGroupMapping(ctx context.Context, gm *GroupMapping) (*GroupMapping, error) {
	payload, err := Marshal(gm)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.CreateGroupMappingURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	created, err := Unmarshal[GroupMapping](response.Body)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (client *Client) UpdateGroupMapping(ctx context.Context, gm *GroupMapping, id int) (*GroupMapping, error) {
	payload, err := Marshal(gm)
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.UpdateGroupMappingURL(id), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	updated, err := Unmarshal[GroupMapping](response.Body)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (client *Client) DeleteGroupMapping(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.DeleteGroupMappingURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) GetGroupMapping(ctx context.Context, id int) (*GroupMapping, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetGroupMappingURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, GroupMappingNotFound
		}
		return nil, client.ErrorFromResponse(response)
	}

	gm, err := Unmarshal[GroupMapping](response.Body)
	if err != nil {
		return nil, err
	}

	return &gm, nil
}

func (client *Client) CreateGroupMappingURL() string {
	return fmt.Sprintf(CreateGroupMappingPath, client.config.url)
}

func (client *Client) UpdateGroupMappingURL(id int) string {
	return fmt.Sprintf(UpdateGroupMappingPath, client.config.url, id)
}

func (client *Client) DeleteGroupMappingURL(id int) string {
	return fmt.Sprintf(DeleteGroupMappingPath, client.config.url, id)
}

func (client *Client) GetGroupMappingURL(id int) string {
	return fmt.Sprintf(GetGroupMappingPath, client.config.url, id)
}
