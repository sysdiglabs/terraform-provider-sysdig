package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

func (c *sysdigMonitorClient) alertsV2URL() string {
	return fmt.Sprintf("%s/api/v2/alerts", c.URL)
}

func (c *sysdigMonitorClient) alertV2URL(alertID int) string {
	return fmt.Sprintf("%s/api/v2/alerts/%d", c.URL, alertID)
}

func (c *sysdigMonitorClient) labelsDescriptorsV3URL(label string) string {
	return fmt.Sprintf("%s/api/v3/labels/descriptors/%s", c.URL, label)
}

func (c *sysdigMonitorClient) labelsV3URL() string {
	return fmt.Sprintf("%s/api/v3/labels/?limit=6000", c.URL) //6000 is the maximum number of labels a customer can have
}

func (c *sysdigMonitorClient) addNotificationChannelType(ctx context.Context, notificationChannelConfigList []NotificationChannelConfigV2) error {
	// on put/posts the api wants the type of the channel even if it can be inferred
	for i, n := range notificationChannelConfigList {
		nc, err := c.GetNotificationChannelById(ctx, n.ChannelID)
		if err != nil {
			return fmt.Errorf("error getting info for notification channel %d: %w", n.ChannelID, err)
		}
		notificationChannelConfigList[i].Type = nc.Type
	}
	return nil
}

func (c *sysdigMonitorClient) translateScopeSegmentLabels(ctx context.Context, scopedSegmentedConfig *ScopedSegmentedConfig) error {
	// the operand of the scope must be in dot notation
	if scopedSegmentedConfig.Scope != nil {
		for i, e := range scopedSegmentedConfig.Scope.Expressions {
			labelDescriptorV3, err := c.GetLabelDescriptor(ctx, e.Operand)
			if err != nil {
				return fmt.Errorf("error getting descriptor for label %s: %w", e.Operand, err)
			}
			scopedSegmentedConfig.Scope.Expressions[i].Operand = labelDescriptorV3.ID
		}
	}

	// the label descriptor id must be in dot notation
	for i, d := range scopedSegmentedConfig.SegmentBy {
		labelDescriptorV3, err := c.GetLabelDescriptor(ctx, d.ID)
		if err != nil {
			return fmt.Errorf("error getting descriptor for label %s: %w", d.ID, err)
		}
		scopedSegmentedConfig.SegmentBy[i].ID = labelDescriptorV3.ID
	}

	return nil
}

// prometheus

func (c *sysdigMonitorClient) CreateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (createdAlert AlertV2Prometheus, err error) {
	if err = c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList); err != nil {
		return
	}

	body, err := c.createAlertV2(ctx, alert.ToJSON())
	if err != nil {
		return
	}

	createdAlert = AlertV2PrometheusFromJSON(body)
	return
}

func (c *sysdigMonitorClient) UpdateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (updatedAlert AlertV2Prometheus, err error) {
	if err = c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList); err != nil {
		return
	}

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
	if err = c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList); err != nil {
		return
	}

	if err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig); err != nil {
		return
	}

	body, err := c.createAlertV2(ctx, alert.ToJSON())
	if err != nil {
		return
	}

	createdAlert = AlertV2EventFromJSON(body)
	return
}

func (c *sysdigMonitorClient) UpdateAlertV2Event(ctx context.Context, alert AlertV2Event) (updatedAlert AlertV2Event, err error) {
	if err = c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList); err != nil {
		return
	}

	if err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig); err != nil {
		return
	}

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
	if err = c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList); err != nil {
		return
	}

	if err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig); err != nil {
		return
	}

	body, err := c.createAlertV2(ctx, alert.ToJSON())
	if err != nil {
		return
	}

	createdAlert = AlertV2MetricFromJSON(body)
	return
}

func (c *sysdigMonitorClient) UpdateAlertV2Metric(ctx context.Context, alert AlertV2Metric) (updatedAlert AlertV2Metric, err error) {
	if err = c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList); err != nil {
		return
	}

	if err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig); err != nil {
		return
	}

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

// downtime

func (c *sysdigMonitorClient) CreateAlertV2Downtime(ctx context.Context, alert AlertV2Downtime) (createdAlert AlertV2Downtime, err error) {
	if err = c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList); err != nil {
		return
	}

	if err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig); err != nil {
		return
	}

	body, err := c.createAlertV2(ctx, alert.ToJSON())
	if err != nil {
		return
	}

	createdAlert = AlertV2DowntimeFromJSON(body)
	return
}

func (c *sysdigMonitorClient) UpdateAlertV2Downtime(ctx context.Context, alert AlertV2Downtime) (updatedAlert AlertV2Downtime, err error) {
	if err = c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList); err != nil {
		return
	}

	if err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig); err != nil {
		return
	}

	body, err := c.updateAlertV2(ctx, alert.ID, alert.ToJSON())
	if err != nil {
		return
	}

	updatedAlert = AlertV2DowntimeFromJSON(body)
	return
}

func (c *sysdigMonitorClient) GetAlertV2DowntimeById(ctx context.Context, alertID int) (alert AlertV2Downtime, err error) {
	body, err := c.getAlertV2ById(ctx, alertID)
	if err != nil {
		return
	}

	alert = AlertV2DowntimeFromJSON(body)
	return
}

func (c *sysdigMonitorClient) DeleteAlertV2Downtime(ctx context.Context, alertID int) (err error) {
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

// buildLabelDescriptor gets the descriptor of a label in public notation from the v3/labels/descriptors api
// this is not a general solution to get the descriptor for a public notation label since custom labels will not be properly translated
// e.g. the public notation cloud_provider_tag_k8s_io_role_master will not be translated to the correct cloudProvider.tag.k8s.io/role/master id
func (c *sysdigMonitorClient) buildLabelDescriptor(ctx context.Context, label string) (LabelDescriptorV3, error) {
	// always returns 200, even if the label does not exist for the customer
	response, err := c.doSysdigMonitorRequest(ctx, http.MethodGet, c.labelsDescriptorsV3URL(label), nil)
	if err != nil {
		return LabelDescriptorV3{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return LabelDescriptorV3{}, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return LabelDescriptorV3{}, err
	}

	var labelsDescriptorsV3result struct {
		LabelDescriptorV3 `json:"labelDescriptor"`
	}

	err = json.Unmarshal(body, &labelsDescriptorsV3result)
	if err != nil {
		return LabelDescriptorV3{}, err
	}

	return labelsDescriptorsV3result.LabelDescriptorV3, nil
}

func (c *sysdigMonitorClient) getLabels(ctx context.Context, label string) ([]LabelDescriptorV3, error) {

	var labelsResp struct {
		AllLabels []LabelDescriptorV3 `json:"allLabels"`
	}

	response, err := c.doSysdigMonitorRequest(ctx, http.MethodGet, c.labelsV3URL(), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &labelsResp)
	return labelsResp.AllLabels, err

}

var labelCache struct {
	sync.Mutex

	labels []LabelDescriptorV3
}

// GetLabel gets the descriptor from a label in public notation
func (c *sysdigMonitorClient) GetLabelDescriptor(ctx context.Context, label string) (LabelDescriptorV3, error) {
	var alertDescriptior LabelDescriptorV3

	labelCache.Lock()
	defer labelCache.Unlock()

	if len(labelCache.labels) == 0 {
		log.Printf("[DEBUG] GetLabel for %s: fetching all labels", label)
		labelDescriptors, err := c.getLabels(ctx, label)
		if err != nil {
			return alertDescriptior, err
		}
		labelCache.labels = labelDescriptors
	} else {
		log.Printf("[DEBUG] GetLabel for %s: using cached labels", label)
	}

	for _, l := range labelCache.labels {
		if l.PublicID == label {
			return l, nil
		}
	}

	// if the label did not exist, build the descriptor from /v3/labels/descriptor
	log.Printf("[DEBUG] GetLabel for %s: not found in existing customer labels", label)
	return c.buildLabelDescriptor(ctx, label)

}
