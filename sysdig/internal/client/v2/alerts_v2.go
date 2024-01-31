package v2

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

var AlertV2NotFound = errors.New("alert not found")

type (
	AlertV2Type     string
	AlertV2Severity string
	AlertLinkV2Type string
)

const (
	alertsV2Path            = "%s/api/v2/alerts"
	alertV2Path             = "%s/api/v2/alerts/%d"
	labelsV3Path            = "%s/api/v3/labels/?limit=6000"
	labelsV3DescriptorsPath = "%s/api/v3/labels/descriptors/%s"

	AlertV2TypePrometheus          AlertV2Type = "PROMETHEUS"
	AlertV2TypeManual              AlertV2Type = "MANUAL"
	AlertV2TypeEvent               AlertV2Type = "EVENT"
	AlertV2TypeChange              AlertV2Type = "PERCENTAGE_OF_CHANGE"
	AlertV2TypeFormBasedPrometheus AlertV2Type = "FORM_BASED_PROMETHEUS"
	AlertV2TypeGroupOutlier        AlertV2Type = "GROUP_OUTLIERS"

	AlertV2SeverityHigh   AlertV2Severity = "high"
	AlertV2SeverityMedium AlertV2Severity = "medium"
	AlertV2SeverityLow    AlertV2Severity = "low"
	AlertV2SeverityInfo   AlertV2Severity = "info"

	AlertLinkV2TypeDashboard AlertLinkV2Type = "dashboard"
	AlertLinkV2TypeRunbook   AlertLinkV2Type = "runbook"
)

var labelCache struct {
	sync.Mutex

	labels []LabelDescriptorV3
}

type AlertV2Interface interface {
	AlertV2PrometheusInterface
	AlertV2EventInterface
	AlertV2MetricInterface
	AlertV2DowntimeInterface
	AlertV2ChangeInterface
	AlertV2FormBasedPrometheusInterface
	AlertV2GroupOutlierInterface
}

type AlertV2PrometheusInterface interface {
	Base
	CreateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error)
	UpdateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error)
	GetAlertV2Prometheus(ctx context.Context, alertID int) (AlertV2Prometheus, error)
	DeleteAlertV2Prometheus(ctx context.Context, alertID int) error
}

type AlertV2EventInterface interface {
	Base
	CreateAlertV2Event(ctx context.Context, alert AlertV2Event) (AlertV2Event, error)
	UpdateAlertV2Event(ctx context.Context, alert AlertV2Event) (AlertV2Event, error)
	GetAlertV2Event(ctx context.Context, alertID int) (AlertV2Event, error)
	DeleteAlertV2Event(ctx context.Context, alertID int) error
}

type AlertV2MetricInterface interface {
	Base
	CreateAlertV2Metric(ctx context.Context, alert AlertV2Metric) (AlertV2Metric, error)
	UpdateAlertV2Metric(ctx context.Context, alert AlertV2Metric) (AlertV2Metric, error)
	GetAlertV2Metric(ctx context.Context, alertID int) (AlertV2Metric, error)
	DeleteAlertV2Metric(ctx context.Context, alertID int) error
}

type AlertV2ChangeInterface interface {
	Base
	CreateAlertV2Change(ctx context.Context, alert AlertV2Change) (AlertV2Change, error)
	UpdateAlertV2Change(ctx context.Context, alert AlertV2Change) (AlertV2Change, error)
	GetAlertV2Change(ctx context.Context, alertID int) (AlertV2Change, error)
	DeleteAlertV2Change(ctx context.Context, alertID int) error
}

type AlertV2FormBasedPrometheusInterface interface {
	Base
	CreateAlertV2FormBasedPrometheus(ctx context.Context, alert AlertV2FormBasedPrometheus) (AlertV2FormBasedPrometheus, error)
	UpdateAlertV2FormBasedPrometheus(ctx context.Context, alert AlertV2FormBasedPrometheus) (AlertV2FormBasedPrometheus, error)
	GetAlertV2FormBasedPrometheus(ctx context.Context, alertID int) (AlertV2FormBasedPrometheus, error)
	DeleteAlertV2FormBasedPrometheus(ctx context.Context, alertID int) error
}

type AlertV2GroupOutlierInterface interface {
	Base
	CreateAlertV2GroupOutlier(ctx context.Context, alert AlertV2GroupOutlier) (AlertV2GroupOutlier, error)
	UpdateAlertV2GroupOutlier(ctx context.Context, alert AlertV2GroupOutlier) (AlertV2GroupOutlier, error)
	GetAlertV2GroupOutlier(ctx context.Context, alertID int) (AlertV2GroupOutlier, error)
	DeleteAlertV2GroupOutlier(ctx context.Context, alertID int) error
}

type AlertV2DowntimeInterface interface {
	Base
	CreateAlertV2Downtime(ctx context.Context, alert AlertV2Downtime) (AlertV2Downtime, error)
	UpdateAlertV2Downtime(ctx context.Context, alert AlertV2Downtime) (AlertV2Downtime, error)
	GetAlertV2Downtime(ctx context.Context, alertID int) (AlertV2Downtime, error)
	DeleteAlertV2Downtime(ctx context.Context, alertID int) error
}

func (client *Client) CreateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	payload, err := Marshal(alertV2PrometheusWrapper{Alert: alert})
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	body, err := client.createAlertV2(ctx, payload)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	wrapper, err := Unmarshal[alertV2PrometheusWrapper](body)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) UpdateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	payload, err := Marshal(alertV2PrometheusWrapper{Alert: alert})
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	body, err := client.updateAlertV2(ctx, alert.ID, payload)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	wrapper, err := Unmarshal[alertV2PrometheusWrapper](body)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) GetAlertV2Prometheus(ctx context.Context, alertID int) (AlertV2Prometheus, error) {
	body, err := client.getAlertV2(ctx, alertID)
	if err != nil {
		return AlertV2Prometheus{}, err
	}
	wrapper, err := Unmarshal[alertV2PrometheusWrapper](body)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) DeleteAlertV2Prometheus(ctx context.Context, alertID int) error {
	return client.deleteAlertV2(ctx, alertID)
}

func (client *Client) CreateAlertV2Event(ctx context.Context, alert AlertV2Event) (AlertV2Event, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Event{}, err
	}

	err = client.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Event{}, err
	}

	payload, err := Marshal(alertV2EventWrapper{Alert: alert})
	if err != nil {
		return AlertV2Event{}, err
	}

	body, err := client.createAlertV2(ctx, payload)
	if err != nil {
		return AlertV2Event{}, err
	}

	wrapper, err := Unmarshal[alertV2EventWrapper](body)
	if err != nil {
		return AlertV2Event{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) UpdateAlertV2Event(ctx context.Context, alert AlertV2Event) (AlertV2Event, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Event{}, err
	}

	err = client.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Event{}, err
	}

	payload, err := Marshal(alertV2EventWrapper{Alert: alert})
	if err != nil {
		return AlertV2Event{}, err
	}

	body, err := client.updateAlertV2(ctx, alert.ID, payload)
	if err != nil {
		return AlertV2Event{}, err
	}

	wrapper, err := Unmarshal[alertV2EventWrapper](body)
	if err != nil {
		return AlertV2Event{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) GetAlertV2Event(ctx context.Context, alertID int) (AlertV2Event, error) {
	body, err := client.getAlertV2(ctx, alertID)
	if err != nil {
		return AlertV2Event{}, err
	}

	wrapper, err := Unmarshal[alertV2EventWrapper](body)
	if err != nil {
		return AlertV2Event{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) DeleteAlertV2Event(ctx context.Context, alertID int) error {
	return client.deleteAlertV2(ctx, alertID)
}

func (client *Client) CreateAlertV2Metric(ctx context.Context, alert AlertV2Metric) (AlertV2Metric, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Metric{}, err
	}

	err = client.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Metric{}, err
	}

	payload, err := Marshal(alertV2MetricWrapper{Alert: alert})
	if err != nil {
		return AlertV2Metric{}, err
	}

	body, err := client.createAlertV2(ctx, payload)
	if err != nil {
		return AlertV2Metric{}, err
	}

	wrapper, err := Unmarshal[alertV2MetricWrapper](body)
	if err != nil {
		return AlertV2Metric{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) UpdateAlertV2Metric(ctx context.Context, alert AlertV2Metric) (AlertV2Metric, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Metric{}, err
	}

	err = client.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Metric{}, err
	}

	payload, err := Marshal(alertV2MetricWrapper{Alert: alert})
	if err != nil {
		return AlertV2Metric{}, err
	}

	body, err := client.updateAlertV2(ctx, alert.ID, payload)
	if err != nil {
		return AlertV2Metric{}, err
	}

	wrapper, err := Unmarshal[alertV2MetricWrapper](body)
	if err != nil {
		return AlertV2Metric{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) GetAlertV2Metric(ctx context.Context, alertID int) (AlertV2Metric, error) {
	body, err := client.getAlertV2(ctx, alertID)
	if err != nil {
		return AlertV2Metric{}, err
	}

	wrapper, err := Unmarshal[alertV2MetricWrapper](body)
	if err != nil {
		return AlertV2Metric{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) DeleteAlertV2Metric(ctx context.Context, alertID int) error {
	return client.deleteAlertV2(ctx, alertID)
}

func (client *Client) CreateAlertV2Downtime(ctx context.Context, alert AlertV2Downtime) (AlertV2Downtime, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	err = client.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	payload, err := Marshal(alertV2DowntimeWrapper{Alert: alert})
	if err != nil {
		return AlertV2Downtime{}, err
	}

	body, err := client.createAlertV2(ctx, payload)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	wrapper, err := Unmarshal[alertV2DowntimeWrapper](body)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	return wrapper.Alert, err
}

func (client *Client) UpdateAlertV2Downtime(ctx context.Context, alert AlertV2Downtime) (AlertV2Downtime, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	err = client.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	payload, err := Marshal(alertV2DowntimeWrapper{Alert: alert})
	if err != nil {
		return AlertV2Downtime{}, err
	}

	body, err := client.updateAlertV2(ctx, alert.ID, payload)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	wrapper, err := Unmarshal[alertV2DowntimeWrapper](body)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	return wrapper.Alert, err
}

func (client *Client) GetAlertV2Downtime(ctx context.Context, alertID int) (AlertV2Downtime, error) {
	body, err := client.getAlertV2(ctx, alertID)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	wrapper, err := Unmarshal[alertV2DowntimeWrapper](body)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) DeleteAlertV2Downtime(ctx context.Context, alertID int) error {
	return client.deleteAlertV2(ctx, alertID)
}

func (client *Client) CreateAlertV2Change(ctx context.Context, alert AlertV2Change) (AlertV2Change, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Change{}, err
	}

	err = client.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Change{}, err
	}

	payload, err := Marshal(alertV2ChangeWrapper{Alert: alert})
	if err != nil {
		return AlertV2Change{}, err
	}

	body, err := client.createAlertV2(ctx, payload)
	if err != nil {
		return AlertV2Change{}, err
	}

	wrapper, err := Unmarshal[alertV2ChangeWrapper](body)
	if err != nil {
		return AlertV2Change{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) UpdateAlertV2Change(ctx context.Context, alert AlertV2Change) (AlertV2Change, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Change{}, err
	}

	err = client.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Change{}, err
	}

	payload, err := Marshal(alertV2ChangeWrapper{Alert: alert})
	if err != nil {
		return AlertV2Change{}, err
	}

	body, err := client.updateAlertV2(ctx, alert.ID, payload)
	if err != nil {
		return AlertV2Change{}, err
	}

	wrapper, err := Unmarshal[alertV2ChangeWrapper](body)
	if err != nil {
		return AlertV2Change{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) GetAlertV2Change(ctx context.Context, alertID int) (AlertV2Change, error) {
	body, err := client.getAlertV2(ctx, alertID)
	if err != nil {
		return AlertV2Change{}, err
	}

	wrapper, err := Unmarshal[alertV2ChangeWrapper](body)
	if err != nil {
		return AlertV2Change{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) DeleteAlertV2Change(ctx context.Context, alertID int) error {
	return client.deleteAlertV2(ctx, alertID)
}

func (client *Client) CreateAlertV2FormBasedPrometheus(ctx context.Context, alert AlertV2FormBasedPrometheus) (AlertV2FormBasedPrometheus, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	err = client.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	payload, err := Marshal(alertV2FormBasedPrometheusWrapper{Alert: alert})
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	body, err := client.createAlertV2(ctx, payload)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	wrapper, err := Unmarshal[alertV2FormBasedPrometheusWrapper](body)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) UpdateAlertV2FormBasedPrometheus(ctx context.Context, alert AlertV2FormBasedPrometheus) (AlertV2FormBasedPrometheus, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	err = client.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	payload, err := Marshal(alertV2FormBasedPrometheusWrapper{Alert: alert})
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	body, err := client.updateAlertV2(ctx, alert.ID, payload)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	wrapper, err := Unmarshal[alertV2FormBasedPrometheusWrapper](body)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) GetAlertV2FormBasedPrometheus(ctx context.Context, alertID int) (AlertV2FormBasedPrometheus, error) {
	body, err := client.getAlertV2(ctx, alertID)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	wrapper, err := Unmarshal[alertV2FormBasedPrometheusWrapper](body)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) DeleteAlertV2FormBasedPrometheus(ctx context.Context, alertID int) error {
	return client.deleteAlertV2(ctx, alertID)
}

func (client *Client) CreateAlertV2GroupOutlier(ctx context.Context, alert AlertV2GroupOutlier) (AlertV2GroupOutlier, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	err = client.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	payload, err := Marshal(alertV2GroupOutlierWrapper{Alert: alert})
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	body, err := client.createAlertV2(ctx, payload)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	wrapper, err := Unmarshal[alertV2GroupOutlierWrapper](body)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) UpdateAlertV2GroupOutlier(ctx context.Context, alert AlertV2GroupOutlier) (AlertV2GroupOutlier, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	err = client.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	payload, err := Marshal(alertV2GroupOutlierWrapper{Alert: alert})
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	body, err := client.updateAlertV2(ctx, alert.ID, payload)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	wrapper, err := Unmarshal[alertV2GroupOutlierWrapper](body)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) GetAlertV2GroupOutlier(ctx context.Context, alertID int) (AlertV2GroupOutlier, error) {
	body, err := client.getAlertV2(ctx, alertID)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	wrapper, err := Unmarshal[alertV2GroupOutlierWrapper](body)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	return wrapper.Alert, nil
}

func (client *Client) DeleteAlertV2GroupOutlier(ctx context.Context, alertID int) error {
	return client.deleteAlertV2(ctx, alertID)
}

func (client *Client) createAlertV2(ctx context.Context, alertJson io.Reader) (io.ReadCloser, error) {
	response, err := client.requester.Request(ctx, http.MethodPost, client.alertsV2URL(), alertJson)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return response.Body, nil
}

func (client *Client) updateAlertV2(ctx context.Context, alertID int, alertJson io.Reader) (io.ReadCloser, error) {
	response, err := client.requester.Request(ctx, http.MethodPut, client.alertV2URL(alertID), alertJson)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return response.Body, nil
}

func (client *Client) deleteAlertV2(ctx context.Context, alertID int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.alertV2URL(alertID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) getAlertV2(ctx context.Context, alertID int) (io.ReadCloser, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.alertV2URL(alertID), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, AlertV2NotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return response.Body, nil
}

func (client *Client) addNotificationChannelType(ctx context.Context, notificationChannelConfigList []NotificationChannelConfigV2) error {
	// on put/posts the api wants the type of the channel even if it can be inferred
	for i, n := range notificationChannelConfigList {
		nc, err := client.GetNotificationChannelById(ctx, n.ChannelID)
		if err != nil {
			return fmt.Errorf("error getting info for notification channel %d: %w", n.ChannelID, err)
		}
		notificationChannelConfigList[i].Type = nc.Type
	}
	return nil
}

func (client *Client) translateScopeSegmentLabels(ctx context.Context, scopedSegmentedConfig *ScopedSegmentedConfig) error {
	// the operand of the scope must be in dot notation
	if scopedSegmentedConfig.Scope != nil {
		for i, e := range scopedSegmentedConfig.Scope.Expressions {
			labelDescriptorV3, err := client.getLabelDescriptor(ctx, e.Operand)
			if err != nil {
				return fmt.Errorf("error getting descriptor for label %s: %w", e.Operand, err)
			}
			scopedSegmentedConfig.Scope.Expressions[i].Operand = labelDescriptorV3.ID
		}
	}

	// the label descriptor id must be in dot notation
	for i, d := range scopedSegmentedConfig.SegmentBy {
		labelDescriptorV3, err := client.getLabelDescriptor(ctx, d.ID)
		if err != nil {
			return fmt.Errorf("error getting descriptor for label %s: %w", d.ID, err)
		}
		scopedSegmentedConfig.SegmentBy[i].ID = labelDescriptorV3.ID
	}

	return nil
}

func (client *Client) getLabelDescriptor(ctx context.Context, label string) (LabelDescriptorV3, error) {
	var alertDescriptor LabelDescriptorV3

	labelCache.Lock()
	defer labelCache.Unlock()

	if len(labelCache.labels) == 0 {
		log.Printf("[DEBUG] GetLabel for %s: fetching all labels", label)
		labelDescriptors, err := client.getLabels(ctx)
		if err != nil {
			return alertDescriptor, err
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
	return client.buildLabelDescriptor(ctx, label)
}

// buildLabelDescriptor gets the descriptor of a label in public notation from the v3/labels/descriptors api
// this is not a general solution to get the descriptor for a public notation label since custom labels will not be properly translated
func (client *Client) buildLabelDescriptor(ctx context.Context, label string) (LabelDescriptorV3, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.labelsDescriptorsV3URL(label), nil)
	if err != nil {
		return LabelDescriptorV3{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return LabelDescriptorV3{}, err
	}

	descriptor, err := Unmarshal[labelDescriptorV3](response.Body)
	if err != nil {
		return LabelDescriptorV3{}, err
	}

	return descriptor.LabelDescriptor, nil
}

func (client *Client) getLabels(ctx context.Context) ([]LabelDescriptorV3, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.labelsV3URL(), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	wrapper, err := Unmarshal[labelsDescriptorV3](response.Body)
	if err != nil {
		return nil, err
	}
	return wrapper.AllLabels, nil
}

func (client *Client) alertsV2URL() string {
	return fmt.Sprintf(alertsV2Path, client.config.url)
}

func (client *Client) alertV2URL(alertID int) string {
	return fmt.Sprintf(alertV2Path, client.config.url, alertID)
}

func (client *Client) labelsV3URL() string {
	return fmt.Sprintf(labelsV3Path, client.config.url)
}

func (client *Client) labelsDescriptorsV3URL(label string) string {
	return fmt.Sprintf(labelsV3DescriptorsPath, client.config.url, label)
}
