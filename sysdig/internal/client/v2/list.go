package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	createListPath = "%s/api/secure/falco/lists?skipPolicyV2Msg=%t"
	getListPath    = "%s/api/secure/falco/lists/%d"
	updateListPath = "%s/api/secure/falco/lists/%d?skipPolicyV2Msg=%t"
	deleteListPath = "%s/api/secure/falco/lists/%d?skipPolicyV2Msg=%t"
)

type ListInterface interface {
	Base
	CreateList(ctx context.Context, list List) (List, error)
	GetListByID(ctx context.Context, id int) (List, error)
	UpdateList(ctx context.Context, list List) (List, error)
	DeleteList(ctx context.Context, id int) error
}

func (c *Client) CreateList(ctx context.Context, list List) (createdList List, err error) {
	payload, err := Marshal(list)
	if err != nil {
		return List{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createListURL(), payload)
	if err != nil {
		return List{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return List{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[List](response.Body)
}

func (c *Client) GetListByID(ctx context.Context, id int) (list List, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getListURL(id), nil)
	if err != nil {
		return List{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return List{}, c.ErrorFromResponse(response)
	}

	list, err = Unmarshal[List](response.Body)
	if err != nil {
		return List{}, err
	}

	if list.Version == 0 {
		return List{}, fmt.Errorf("list with ID: %d does not exists", id)
	}

	return list, nil
}

func (c *Client) UpdateList(ctx context.Context, list List) (updatedList List, err error) {
	payload, err := Marshal(list)
	if err != nil {
		return List{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateListURL(list.ID), payload)
	if err != nil {
		return List{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return List{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[List](response.Body)
}

func (c *Client) DeleteList(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteListURL(id), nil)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return c.ErrorFromResponse(response)
	}

	return nil
}

func (c *Client) createListURL() string {
	return fmt.Sprintf(createListPath, c.config.url, c.config.secureSkipPolicyV2Msg)
}

func (c *Client) getListURL(id int) string {
	return fmt.Sprintf(getListPath, c.config.url, id)
}

func (c *Client) updateListURL(id int) string {
	return fmt.Sprintf(updateListPath, c.config.url, id, c.config.secureSkipPolicyV2Msg)
}

func (c *Client) deleteListURL(id int) string {
	return fmt.Sprintf(deleteListPath, c.config.url, id, c.config.secureSkipPolicyV2Msg)
}
