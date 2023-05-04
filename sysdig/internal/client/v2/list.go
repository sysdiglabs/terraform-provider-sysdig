package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	CreateListPath = "%s/api/secure/falco/lists"
	GetListPath    = "%s/api/secure/falco/lists/%d"
	UpdateListPath = "%s/api/secure/falco/lists/%d"
	DeleteListPath = "%s/api/secure/falco/lists/%d"
)

type ListInterface interface {
	CreateList(ctx context.Context, list List) (List, error)
	GetListByID(ctx context.Context, id int) (List, error)
	UpdateList(ctx context.Context, list List) (List, error)
	DeleteList(ctx context.Context, id int) error
}

func (client *Client) CreateList(ctx context.Context, list List) (List, error) {
	payload, err := Marshal[List](list)
	if err != nil {
		return List{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.CreateListURL(), payload)
	if err != nil {
		return List{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return List{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[List](response.Body)
}

func (client *Client) GetListByID(ctx context.Context, id int) (List, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetListURL(id), nil)
	if err != nil {
		return List{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return List{}, client.ErrorFromResponse(response)
	}

	list, err := Unmarshal[List](response.Body)
	if err != nil {
		return List{}, err
	}

	if list.Version == 0 {
		return List{}, fmt.Errorf("list with ID: %d does not exists", id)
	}

	return list, nil
}

func (client *Client) UpdateList(ctx context.Context, list List) (List, error) {
	payload, err := Marshal[List](list)
	if err != nil {
		return List{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.UpdateListURL(list.ID), payload)
	if err != nil {
		return List{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return List{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[List](response.Body)
}

func (client *Client) DeleteList(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.DeleteListURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) CreateListURL() string {
	return fmt.Sprintf(CreateListPath, client.config.url)
}

func (client *Client) GetListURL(id int) string {
	return fmt.Sprintf(GetListPath, client.config.url, id)
}

func (client *Client) UpdateListURL(id int) string {
	return fmt.Sprintf(UpdateListPath, client.config.url, id)
}

func (client *Client) DeleteListURL(id int) string {
	return fmt.Sprintf(DeleteListPath, client.config.url, id)
}
