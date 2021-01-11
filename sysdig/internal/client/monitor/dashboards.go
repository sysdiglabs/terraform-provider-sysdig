package monitor

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor/model"
)

func (client *sysdigMonitorClient) GetDashboardByID(ctx context.Context, ID int) (*model.Dashboard, error) {
	res, err := client.doSysdigMonitorRequest(ctx, http.MethodGet, client.getDashboardUrl(ID), nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	return model.DashboardFromJSON(body), nil
}

func (client *sysdigMonitorClient) CreateDashboard(ctx context.Context, dashboard *model.Dashboard) (*model.Dashboard, error) {
	res, err := client.doSysdigMonitorRequest(ctx, http.MethodPost, client.getDashboardsUrl(), dashboard.ToJSON())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf(string(body))
	}

	return model.DashboardFromJSON(body), nil
}

func (client *sysdigMonitorClient) UpdateDashboard(ctx context.Context, dashboard *model.Dashboard) (*model.Dashboard, error) {
	res, err := client.doSysdigMonitorRequest(ctx, http.MethodPut, client.getDashboardUrl(dashboard.ID), dashboard.ToJSON())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf(string(body))
	}

	return model.DashboardFromJSON(body), nil
}

func (client *sysdigMonitorClient) DeleteDashboard(ctx context.Context, ID int) error {
	res, err := client.doSysdigMonitorRequest(ctx, http.MethodDelete, client.getDashboardUrl(ID), nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf(string(body))
	}

	return nil
}

func (client *sysdigMonitorClient) getDashboardsUrl() string {
	return fmt.Sprintf("%s/api/v3/dashboards", client.URL)
}

func (client *sysdigMonitorClient) getDashboardUrl(id int) string {
	return fmt.Sprintf("%s/api/v3/dashboards/%d", client.URL, id)
}
