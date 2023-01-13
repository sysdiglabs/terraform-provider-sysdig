package monitor

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (client *sysdigMonitorClient) GetCloudAccountById(ctx context.Context, id int) (*CloudAccount, error) {
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
	return CloudAccountFromJSON(body), nil
}

func (client *sysdigMonitorClient) CreateCloudAccount(ctx context.Context, provider *CloudAccount) (*CloudAccount, error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodPost, client.getProvidersUrl(), CloudAccountToJSON(provider))
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
	return CloudAccountFromJSON(body), nil
}

func (client *sysdigMonitorClient) UpdateCloudAccount(ctx context.Context, id int, provider *CloudAccount) (*CloudAccount, error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodPut, client.getProviderUrl(id), CloudAccountToJSON(provider))
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
	return CloudAccountFromJSON(body), nil
}

func (client *sysdigMonitorClient) DeleteCloudAccountById(ctx context.Context, id int) error {
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
