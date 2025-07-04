package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrGroupMappingNotFound = errors.New("group mapping not found")

const (
	createGroupMappingPath = "%s/api/groupmappings"
	updateGroupMappingPath = "%s/api/groupmappings/%d"
	deleteGroupMappingPath = "%s/api/groupmappings/%d"
	getGroupMappingPath    = "%s/api/groupmappings/%d"
)

type GroupMappingInterface interface {
	Base
	CreateGroupMapping(ctx context.Context, gm *GroupMapping) (*GroupMapping, error)
	UpdateGroupMapping(ctx context.Context, gm *GroupMapping, id int) (*GroupMapping, error)
	DeleteGroupMapping(ctx context.Context, id int) error
	GetGroupMapping(ctx context.Context, id int) (*GroupMapping, error)
}

func (c *Client) CreateGroupMapping(ctx context.Context, gm *GroupMapping) (mapping *GroupMapping, err error) {
	payload, err := Marshal(gm)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createGroupMappingURL(), payload)
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

	return Unmarshal[*GroupMapping](response.Body)
}

func (c *Client) UpdateGroupMapping(ctx context.Context, gm *GroupMapping, id int) (mapping *GroupMapping, err error) {
	payload, err := Marshal(gm)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateGroupMappingURL(id), payload)
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

	return Unmarshal[*GroupMapping](response.Body)
}

func (c *Client) DeleteGroupMapping(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteGroupMappingURL(id), nil)
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

func (c *Client) GetGroupMapping(ctx context.Context, id int) (mapping *GroupMapping, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getGroupMappingURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrGroupMappingNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[*GroupMapping](response.Body)
}

func (c *Client) createGroupMappingURL() string {
	return fmt.Sprintf(createGroupMappingPath, c.config.url)
}

func (c *Client) updateGroupMappingURL(id int) string {
	return fmt.Sprintf(updateGroupMappingPath, c.config.url, id)
}

func (c *Client) deleteGroupMappingURL(id int) string {
	return fmt.Sprintf(deleteGroupMappingPath, c.config.url, id)
}

func (c *Client) getGroupMappingURL(id int) string {
	return fmt.Sprintf(getGroupMappingPath, c.config.url, id)
}
