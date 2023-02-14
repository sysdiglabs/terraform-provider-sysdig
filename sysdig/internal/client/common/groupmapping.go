package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var GroupMappingNotFound = errors.New("group mapping not found")

func (client *sysdigCommonClient) CreateGroupMapping(ctx context.Context, request *GroupMapping) (*GroupMapping, error) {
	payload, err := request.ToJSON()
	if err != nil {
		return nil, err
	}

	response, err := client.doSysdigCommonRequest(ctx, http.MethodPost, client.CreateGroupMappingUrl(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var groupMapping GroupMapping
	err = json.Unmarshal(body, &groupMapping)
	if err != nil {
		return nil, err
	}

	return &groupMapping, nil
}

func (client *sysdigCommonClient) UpdateGroupMapping(ctx context.Context, request *GroupMapping, id int) (*GroupMapping, error) {
	payload, err := request.ToJSON()
	if err != nil {
		return nil, err
	}

	response, err := client.doSysdigCommonRequest(ctx, http.MethodPut, client.UpdateGroupMappingUrl(id), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var groupMapping GroupMapping
	err = json.Unmarshal(body, &groupMapping)
	if err != nil {
		return nil, err
	}

	return &groupMapping, nil
}

func (client *sysdigCommonClient) DeleteGroupMapping(ctx context.Context, id int) error {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodDelete, client.DeleteGroupMappingUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return errorFromResponse(response)
	}

	return nil
}

func (client *sysdigCommonClient) GetGroupMapping(ctx context.Context, id int) (*GroupMapping, error) {
	response, err := client.doSysdigCommonRequest(ctx, http.MethodGet, client.GetGroupMappingUrl(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, GroupMappingNotFound
		}
		return nil, errorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var groupMapping GroupMapping
	err = json.Unmarshal(body, &groupMapping)
	if err != nil {
		return nil, err
	}

	return &groupMapping, nil
}

func (client *sysdigCommonClient) GetGroupMappingUrl(id int) string {
	return fmt.Sprintf("%s/api/groupmappings/%d", client.URL, id)
}

func (client *sysdigCommonClient) CreateGroupMappingUrl() string {
	return fmt.Sprintf("%s/api/groupmappings", client.URL)
}

func (client *sysdigCommonClient) UpdateGroupMappingUrl(id int) string {
	return fmt.Sprintf("%s/api/groupmappings/%d", client.URL, id)
}

func (client *sysdigCommonClient) DeleteGroupMappingUrl(id int) string {
	return fmt.Sprintf("%s/api/groupmappings/%d", client.URL, id)
}
