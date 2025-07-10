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
	GetDashboardByID(ctx context.Context, ID int) (*Dashboard, error)
	CreateDashboard(ctx context.Context, dashboard *Dashboard) (*Dashboard, error)
	UpdateDashboard(ctx context.Context, dashboard *Dashboard) (*Dashboard, error)
	DeleteDashboard(ctx context.Context, ID int) error
}

func (c *Client) GetDashboardByID(ctx context.Context, ID int) (dashboard *Dashboard, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getDashboardURL(ID), nil)
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

	wrapper, err := Unmarshal[*dashboardWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return wrapper.Dashboard, nil
}

func (c *Client) CreateDashboard(ctx context.Context, dashboard *Dashboard) (createdDashboard *Dashboard, err error) {
	payload, err := Marshal(dashboardWrapper{Dashboard: dashboard})
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.getDashboardsURL(), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[*dashboardWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return wrapper.Dashboard, nil
}

func (c *Client) UpdateDashboard(ctx context.Context, dashboard *Dashboard) (updatedDashboard *Dashboard, err error) {
	payload, err := Marshal(dashboardWrapper{Dashboard: dashboard})
	if err != nil {
		return nil, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.getDashboardURL(dashboard.ID), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, c.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[*dashboardWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return wrapper.Dashboard, nil
}

func (c *Client) DeleteDashboard(ctx context.Context, ID int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.getDashboardURL(ID), nil)
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

func (c *Client) getDashboardsURL() string {
	return fmt.Sprintf(dashboardsPath, c.config.url)
}

func (c *Client) getDashboardURL(id int) string {
	return fmt.Sprintf(dashboardPath, c.config.url, id)
}
