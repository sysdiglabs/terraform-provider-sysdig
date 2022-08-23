package secure

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (client *sysdigSecureClient) CreateList(ctx context.Context, listRequest List) (list List, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPost, client.GetListsUrl(), listRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	list, err = ListFromJSON(body)
	return
}

func (client *sysdigSecureClient) GetListById(ctx context.Context, id int) (list List, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.GetListUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	list, err = ListFromJSON(body)
	if err != nil {
		return
	}

	if list.Version == 0 {
		err = fmt.Errorf("List with ID: %d does not exists", id)
		return
	}
	return
}

func (client *sysdigSecureClient) UpdateList(ctx context.Context, listRequest List) (list List, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPut, client.GetListUrl(listRequest.ID), listRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, _ := io.ReadAll(response.Body)
	return ListFromJSON(body)
}

func (client *sysdigSecureClient) DeleteList(ctx context.Context, id int) error {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodDelete, client.GetListUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}
	return nil
}

func (client *sysdigSecureClient) GetListsUrl() string {
	return fmt.Sprintf("%s/api/secure/falco/lists", client.URL)
}

func (client *sysdigSecureClient) GetListUrl(id int) string {
	return fmt.Sprintf("%s/api/secure/falco/lists/%d", client.URL, id)
}
