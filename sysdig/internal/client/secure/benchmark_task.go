package secure

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigSecureClient) createBenchmarkTaskURL() string {
	return fmt.Sprintf("%s/api/compliance/v2/tasks", client.URL)
}

func (client *sysdigSecureClient) benchmarkTaskByIdURL(id string) string {
	return fmt.Sprintf("%s/api/compliance/v2/tasks/%s", client.URL, id)
}

func (client *sysdigSecureClient) setBenchmarkTaskEnabledURL(id string, enabled bool) string {
	if enabled {
		return fmt.Sprintf("%s/api/compliance/v2/tasks/%s/enable", client.URL, id)
	}

	return fmt.Sprintf("%s/api/compliance/v2/tasks/%s/disable", client.URL, id)
}

func (client *sysdigSecureClient) CreateBenchmarkTask(ctx context.Context, task *BenchmarkTask) (*BenchmarkTask, error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPost, client.createBenchmarkTaskURL(), task.ToJSON())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errorFromResponse(response)
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(response.Body)
	return BenchmarkTaskFromJSON(bodyBytes), nil
}

func (client *sysdigSecureClient) GetBenchmarkTask(ctx context.Context, id string) (*BenchmarkTask, error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.benchmarkTaskByIdURL(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errorFromResponse(response)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return BenchmarkTaskFromJSON(bodyBytes), nil
}

func (client *sysdigSecureClient) DeleteBenchmarkTask(ctx context.Context, id string) error {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodDelete, client.benchmarkTaskByIdURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}
	return nil
}

func (client *sysdigSecureClient) SetBenchmarkTaskEnabled(ctx context.Context, id string, enabled bool) error {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPut, client.setBenchmarkTaskEnabledURL(id, enabled), nil)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}
	return nil
}
