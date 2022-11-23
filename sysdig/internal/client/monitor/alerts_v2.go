package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *sysdigMonitorClient) alertsV2URL() string {
	return fmt.Sprintf("%s/api/v2/alerts", c.URL)
}

func (c *sysdigMonitorClient) alertV2URL(alertID int) string {
	return fmt.Sprintf("%s/api/v2/alerts/%d", c.URL, alertID)
}

func (c *sysdigMonitorClient) labelDescriptorURL(label string) string {
	return fmt.Sprintf("%s/api/v3/labels/descriptors/%s", c.URL, label)
}

// prometheus

func (c *sysdigMonitorClient) CreateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (createdAlert AlertV2Prometheus, err error) {
	body, err := c.createAlertV2(ctx, alert.ToJSON())
	if err != nil {
		return
	}
	createdAlert = AlertV2PrometheusFromJSON(body)

	// this fixes the APIs bug of not setting the default group on the response of the create method
	if createdAlert.Group == "" {
		createdAlert.Group = "default"
	}
	return
}

func (c *sysdigMonitorClient) UpdateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (updatedAlert AlertV2Prometheus, err error) {
	body, err := c.updateAlertV2(ctx, alert.ID, alert.ToJSON())
	if err != nil {
		return
	}

	updatedAlert = AlertV2PrometheusFromJSON(body)
	return
}

func (c *sysdigMonitorClient) GetAlertV2PrometheusById(ctx context.Context, alertID int) (alert AlertV2Prometheus, err error) {
	body, err := c.getAlertV2ById(ctx, alertID)
	if err != nil {
		return
	}

	alert = AlertV2PrometheusFromJSON(body)
	return
}

func (c *sysdigMonitorClient) DeleteAlertV2Prometheus(ctx context.Context, alertID int) (err error) {
	return c.deleteAlertV2(ctx, alertID)
}

// event

func (c *sysdigMonitorClient) CreateAlertV2Event(ctx context.Context, alert AlertV2Event) (createdAlert AlertV2Event, err error) {
	body, err := c.createAlertV2(ctx, alert.ToJSON())
	if err != nil {
		return
	}
	createdAlert = AlertV2EventFromJSON(body)

	// this fixes the APIs bug of not setting the default group on the response of the create method
	if createdAlert.Group == "" {
		createdAlert.Group = "default"
	}
	return
}

func (c *sysdigMonitorClient) UpdateAlertV2Event(ctx context.Context, alert AlertV2Event) (updatedAlert AlertV2Event, err error) {
	body, err := c.updateAlertV2(ctx, alert.ID, alert.ToJSON())
	if err != nil {
		return
	}

	updatedAlert = AlertV2EventFromJSON(body)
	return
}

func (c *sysdigMonitorClient) GetAlertV2EventById(ctx context.Context, alertID int) (alert AlertV2Event, err error) {
	body, err := c.getAlertV2ById(ctx, alertID)
	if err != nil {
		return
	}

	alert = AlertV2EventFromJSON(body)
	return
}

func (c *sysdigMonitorClient) DeleteAlertV2Event(ctx context.Context, alertID int) (err error) {
	return c.deleteAlertV2(ctx, alertID)
}

// metric

func (c *sysdigMonitorClient) CreateAlertV2Metric(ctx context.Context, alert AlertV2Metric) (createdAlert AlertV2Metric, err error) {
	body, err := c.createAlertV2(ctx, alert.ToJSON())
	if err != nil {
		return
	}
	createdAlert = AlertV2MetricFromJSON(body)

	// this fixes the APIs bug of not setting the default group on the response of the create method
	if createdAlert.Group == "" {
		createdAlert.Group = "default"
	}
	return
}

func (c *sysdigMonitorClient) UpdateAlertV2Metric(ctx context.Context, alert AlertV2Metric) (updatedAlert AlertV2Metric, err error) {
	body, err := c.updateAlertV2(ctx, alert.ID, alert.ToJSON())
	if err != nil {
		return
	}

	updatedAlert = AlertV2MetricFromJSON(body)
	return
}

func (c *sysdigMonitorClient) GetAlertV2MetricById(ctx context.Context, alertID int) (alert AlertV2Metric, err error) {
	body, err := c.getAlertV2ById(ctx, alertID)
	if err != nil {
		return
	}

	alert = AlertV2MetricFromJSON(body)
	return
}

func (c *sysdigMonitorClient) DeleteAlertV2Metric(ctx context.Context, alertID int) (err error) {
	return c.deleteAlertV2(ctx, alertID)
}

// helpers

func (c *sysdigMonitorClient) createAlertV2(ctx context.Context, alertJson io.Reader) (responseBody []byte, err error) {
	response, err := c.doSysdigMonitorRequest(ctx, http.MethodPost, c.alertsV2URL(), alertJson)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, err := io.ReadAll(response.Body)
	return body, err
}

func (c *sysdigMonitorClient) deleteAlertV2(ctx context.Context, alertID int) error {
	response, err := c.doSysdigMonitorRequest(ctx, http.MethodDelete, c.alertV2URL(alertID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return errorFromResponse(response)
	}

	return err
}

func (c *sysdigMonitorClient) updateAlertV2(ctx context.Context, alertID int, alertJson io.Reader) (responseBody []byte, err error) {
	response, err := c.doSysdigMonitorRequest(ctx, http.MethodPut, c.alertV2URL(alertID), alertJson)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, err := io.ReadAll(response.Body)
	return body, err
}

func (c *sysdigMonitorClient) getAlertV2ById(ctx context.Context, alertID int) (respBody []byte, err error) {
	response, err := c.doSysdigMonitorRequest(ctx, http.MethodGet, c.alertV2URL(alertID), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return
	}

	body, err := io.ReadAll(response.Body)
	return body, err
}

func (c *sysdigMonitorClient) GetLabelDescriptor(ctx context.Context, label string) (LabelDescriptorV3, error) {
	var alertDescriptior LabelDescriptorV3

	// always returns 200, even if the label does not exist
	response, err := c.doSysdigMonitorRequest(ctx, http.MethodGet, c.labelDescriptorURL(label), nil)
	if err != nil {
		return alertDescriptior, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return alertDescriptior, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return alertDescriptior, err
	}

	err = json.Unmarshal(body, &alertDescriptior)
	if err != nil {
		return alertDescriptior, err
	}

	return alertDescriptior, nil
}
