package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	dashboardsPath = "%s/api/v3/dashboards"
	dashboardPath  = "%s/api/v3/dashboards/%d"
)

type DashboardInterface interface {
	GetDashboard(ctx context.Context, ID int) (*Dashboard, error)
	CreateDashboard(ctx context.Context, dashboard *Dashboard) (*Dashboard, error)
	UpdateDashboard(ctx context.Context, dashboard *Dashboard) (*Dashboard, error)
	DeleteDashboard(ctx context.Context, ID int) error
}

func (client *Client) GetDashboard(ctx context.Context, ID int) (*Dashboard, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.getDashboardURL(ID), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[*dashboardWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return wrapper.Dashboard, nil
}

func (client *Client) CreateDashboard(ctx context.Context, dashboard *Dashboard) (*Dashboard, error) {
	payload, err := Marshal(dashboardWrapper{Dashboard: dashboard})
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.getDashboardsURL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[*dashboardWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return wrapper.Dashboard, nil
}

func (client *Client) UpdateDashboard(ctx context.Context, dashboard *Dashboard) (*Dashboard, error) {
	payload, err := Marshal(dashboardWrapper{Dashboard: dashboard})
	if err != nil {
		return nil, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.getDashboardURL(dashboard.ID), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[*dashboardWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return wrapper.Dashboard, nil
}

func (client *Client) DeleteDashboard(ctx context.Context, ID int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.getDashboardURL(ID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) getDashboardsURL() string {
	return fmt.Sprintf(dashboardsPath, client.config.url)
}

func (client *Client) getDashboardURL(id int) string {
	return fmt.Sprintf(dashboardPath, client.config.url, id)
}
