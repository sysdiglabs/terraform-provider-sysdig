package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	providersPath         = "%v/api/v2/providers"
	costCloudAccountPath  = "%s/api/cloudaccount"
	costProviderURL       = "%s/api/cloudaccount/features/cost/account?id=%d"
	updateCostProviderURL = "%s/api/cloudaccount/features/cost"
)

type CloudAccountMonitorInterface interface {
	Base
	CreateCloudAccountMonitor(ctx context.Context, provider *CloudAccountMonitor) (*CloudAccountMonitor, error)
	CreateCloudAccountMonitorForCost(ctx context.Context, provider *CloudAccountMonitorForCost) (*CloudAccountCreatedForCost, error)
	UpdateCloudAccountMonitor(ctx context.Context, id int, provider *CloudAccountMonitor) (*CloudAccountMonitor, error)
	UpdateCloudAccountMonitorForCost(ctx context.Context, provider *CloudAccountCostProvider) (*CloudAccountCostProvider, error)
	GetCloudAccountMonitorByID(ctx context.Context, id int) (*CloudAccountMonitor, error)
	GetCloudAccountMonitorForCostByID(ctx context.Context, id int) (*CloudAccountCostProvider, error)
	DeleteCloudAccountMonitor(ctx context.Context, id int) error
}

func (c *Client) CreateCloudAccountMonitor(ctx context.Context, provider *CloudAccountMonitor) (createdProvider *CloudAccountMonitor, err error) {
	payload, err := Marshal(provider)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.getProvidersURL(), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusCreated {
		return nil, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[cloudAccountWrapperMonitor](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.CloudAccount, nil
}

func (c *Client) CreateCloudAccountMonitorForCost(ctx context.Context, provider *CloudAccountMonitorForCost) (createdProvider *CloudAccountCreatedForCost, err error) {
	payload, err := Marshal(provider)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.getCostProvidersURL(), payload)
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

	wrapper, err := Unmarshal[CloudAccountCreatedForCost](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper, nil
}

func (c *Client) UpdateCloudAccountMonitor(ctx context.Context, id int, provider *CloudAccountMonitor) (updatedProvider *CloudAccountMonitor, err error) {
	payload, err := Marshal(provider)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.getProviderURL(id), payload)
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

	wrapper, err := Unmarshal[cloudAccountWrapperMonitor](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.CloudAccount, nil
}

func (c *Client) UpdateCloudAccountMonitorForCost(ctx context.Context, provider *CloudAccountCostProvider) (updatedProvider *CloudAccountCostProvider, err error) {
	payload, err := Marshal(provider)
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.getUpdateCostProviderURL(), payload)
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

	wrapper, err := Unmarshal[CloudAccountCostProviderWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.CloudAccountCostProvider, nil
}

func (c *Client) GetCloudAccountMonitorByID(ctx context.Context, id int) (account *CloudAccountMonitor, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getProviderURL(id), nil)
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

	wrapper, err := Unmarshal[cloudAccountWrapperMonitor](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.CloudAccount, nil
}

func (c *Client) GetCloudAccountMonitorForCostByID(ctx context.Context, id int) (provider *CloudAccountCostProvider, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getCostProviderURL(id), nil)
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

	wrapper, err := Unmarshal[CloudAccountCostProviderWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.CloudAccountCostProvider, nil
}

func (c *Client) DeleteCloudAccountMonitor(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.getProviderURL(id), nil)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return c.ErrorFromResponse(response)
	}

	return nil
}

func (c *Client) getProviderURL(id int) string {
	return fmt.Sprintf("%v/%v", c.getProvidersURL(), id)
}

func (c *Client) getProvidersURL() string {
	return fmt.Sprintf(providersPath, c.config.url)
}

func (c *Client) getCostProvidersURL() string {
	return fmt.Sprintf(costCloudAccountPath, c.config.url)
}

func (c *Client) getCostProviderURL(id int) string {
	return fmt.Sprintf(costProviderURL, c.config.url, id)
}

func (c *Client) getUpdateCostProviderURL() string {
	return fmt.Sprintf(updateCostProviderURL, c.config.url)
}
