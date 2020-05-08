package monitor

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *sysdigMonitorClient) CreateAlert(alert Alert) (createdAlert Alert, err error) {
	response, err := c.doSysdigMonitorRequest(http.MethodPost, c.alertsURL(), alert.ToJSON())
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		err = errors.New(string(body))
		return
	}

	defer response.Body.Close()

	return AlertFromJSON(body), nil
}

func (c *sysdigMonitorClient) DeleteAlert(alertID int) error {
	response, err := c.doSysdigMonitorRequest(http.MethodDelete, c.alertURL(alertID), nil)

	defer response.Body.Close()

	return err
}

func (c *sysdigMonitorClient) UpdateAlert(alert Alert) (updatedAlert Alert, err error) {
	response, err := c.doSysdigMonitorRequest(http.MethodPut, c.alertURL(alert.ID), alert.ToJSON())
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		err = errors.New(string(body))
		return
	}

	defer response.Body.Close()

	return AlertFromJSON(body), nil
}

func (c *sysdigMonitorClient) GetAlertById(alertID int) (alert Alert, err error) {
	response, err := c.doSysdigMonitorRequest(http.MethodGet, c.alertURL(alertID), nil)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		err = errors.New(string(body))
		return
	}

	defer response.Body.Close()

	return AlertFromJSON(body), nil
}

func (c *sysdigMonitorClient) alertsURL() string {
	return fmt.Sprintf("%s/api/alerts", c.URL)
}

func (c *sysdigMonitorClient) alertURL(alertID int) string {
	return fmt.Sprintf("%s/api/alerts/%d", c.URL, alertID)

}
