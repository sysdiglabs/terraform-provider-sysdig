package monitor

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *sysdigMonitorClient) CreateAlert(ctx context.Context, alert Alert) (createdAlert Alert, err error) {
	response, err := c.doSysdigMonitorRequest(ctx, http.MethodPost, c.alertsURL(), alert.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return AlertFromJSON(body), nil
}

func (c *sysdigMonitorClient) DeleteAlert(ctx context.Context, alertID int) error {
	response, err := c.doSysdigMonitorRequest(ctx, http.MethodDelete, c.alertURL(alertID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}

	return err
}

func (c *sysdigMonitorClient) UpdateAlert(ctx context.Context, alert Alert) (updatedAlert Alert, err error) {
	response, err := c.doSysdigMonitorRequest(ctx, http.MethodPut, c.alertURL(alert.ID), alert.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = errorFromResponse(response)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return AlertFromJSON(body), nil
}

func (c *sysdigMonitorClient) GetAlertById(ctx context.Context, alertID int) (alert Alert, err error) {
	response, err := c.doSysdigMonitorRequest(ctx, http.MethodGet, c.alertURL(alertID), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = errorFromResponse(response)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	return AlertFromJSON(body), nil
}

func (c *sysdigMonitorClient) alertsURL() string {
	return fmt.Sprintf("%s/api/alerts", c.URL)
}

func (c *sysdigMonitorClient) alertURL(alertID int) string {
	return fmt.Sprintf("%s/api/alerts/%d", c.URL, alertID)

}
