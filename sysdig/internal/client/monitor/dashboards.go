package monitor

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor/model"
)

func (client *sysdigMonitorClient) GetDashboardByID(ctx context.Context, ID int) (*model.Dashboard, error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodGet, client.getDashboardUrl(ID), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errorFromResponse(response)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil
	}
	return model.DashboardFromJSON(body), nil
}

func (client *sysdigMonitorClient) CreateDashboard(ctx context.Context, dashboard *model.Dashboard) (*model.Dashboard, error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodPost, client.getDashboardsUrl(), dashboard.ToJSON())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, errorFromResponse(response)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return model.DashboardFromJSON(body), nil
}

func (client *sysdigMonitorClient) UpdateDashboard(ctx context.Context, dashboard *model.Dashboard) (*model.Dashboard, error) {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodPut, client.getDashboardUrl(dashboard.ID), dashboard.ToJSON())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, errorFromResponse(response)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return model.DashboardFromJSON(body), nil
}

func (client *sysdigMonitorClient) DeleteDashboard(ctx context.Context, ID int) error {
	response, err := client.doSysdigMonitorRequest(ctx, http.MethodDelete, client.getDashboardUrl(ID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}

	return nil
}

func (client *sysdigMonitorClient) getDashboardsUrl() string {
	return fmt.Sprintf("%s/api/v3/dashboards", client.URL)
}

func (client *sysdigMonitorClient) getDashboardUrl(id int) string {
	return fmt.Sprintf("%s/api/v3/dashboards/%d", client.URL, id)
}
