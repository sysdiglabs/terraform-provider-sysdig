package monitor

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (client *sysdigMonitorClient) GetCustomerProviderKeyById(ctx context.Context, id int) (*CustomerProviderKey, error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodGet, client.getProviderUrl(id), nil)
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
	return CustomerProviderKeyFromJSON(body), nil
}

func (client *sysdigMonitorClient) CreateCustomerProviderKey(ctx context.Context, provider *CustomerProviderKey) (*CustomerProviderKey, error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodPost, client.getProvidersUrl(), CustomerProviderKeyToJSON(provider))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, errorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return CustomerProviderKeyFromJSON(body), nil
}

func (client *sysdigMonitorClient) UpdateCustomerProviderKey(ctx context.Context, id int, provider *CustomerProviderKey) (*CustomerProviderKey, error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodPut, client.getProviderUrl(id), CustomerProviderKeyToJSON(provider))
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
	return CustomerProviderKeyFromJSON(body), nil
}

func (client *sysdigMonitorClient) DeleteCustomerProviderKeyById(ctx context.Context, id int) error {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodDelete, client.getProviderUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}

	return nil
}

func (client *sysdigMonitorClient) getProviderUrl(id int) string {
	return fmt.Sprintf("%v/%v", client.getProvidersUrl(), id)
}

func (client *sysdigMonitorClient) getProvidersUrl() string {
	return fmt.Sprintf("%v/api/v2/providers", client.URL)
}
