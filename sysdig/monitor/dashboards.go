package monitor

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigMonitorClient) GetDashboardByID(ctx context.Context, ID int) (Dashboard, error) {
	res, err := client.doSysdigMonitorRequest(ctx, http.MethodGet, client.getDashboardUrl(ID), nil)
	if err != nil {
		return Dashboard{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Dashboard{}, nil
	}

	if res.StatusCode != http.StatusOK {
		return Dashboard{}, fmt.Errorf(string(body))
	}

	return DashboardFromJSON(body), nil
}

func (client *sysdigMonitorClient) CreateDashboard(ctx context.Context, dashboard Dashboard) (Dashboard, error) {
	res, err := client.doSysdigMonitorRequest(ctx, http.MethodPost, client.getDashboardsUrl(), dashboard.ToJSON())
	if err != nil {
		return Dashboard{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Dashboard{}, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return Dashboard{}, fmt.Errorf(string(body))
	}

	return DashboardFromJSON(body), nil
}

func (client *sysdigMonitorClient) UpdateDashboard(ctx context.Context, dashboard Dashboard) (Dashboard, error) {
	res, err := client.doSysdigMonitorRequest(ctx, http.MethodPut, client.getDashboardUrl(dashboard.ID), dashboard.ToJSON())
	if err != nil {
		return Dashboard{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Dashboard{}, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return Dashboard{}, fmt.Errorf(string(body))
	}

	return DashboardFromJSON(body), nil
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
