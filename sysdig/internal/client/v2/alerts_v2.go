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

var ErrAlertV2NotFound = errors.New("alert not found")

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
)

const (
	AlertV2TypePrometheus          AlertV2Type = "PROMETHEUS"
	AlertV2TypeManual              AlertV2Type = "MANUAL"
	AlertV2TypeEvent               AlertV2Type = "EVENT"
	AlertV2TypeChange              AlertV2Type = "PERCENTAGE_OF_CHANGE"
	AlertV2TypeFormBasedPrometheus AlertV2Type = "FORM_BASED_PROMETHEUS"
	AlertV2TypeGroupOutlier        AlertV2Type = "GROUP_OUTLIERS"
	AlertV2TypeDowntime            AlertV2Type = "DOWNTIME"
)

const (
	AlertV2SeverityHigh   AlertV2Severity = "high"
	AlertV2SeverityMedium AlertV2Severity = "medium"
	AlertV2SeverityLow    AlertV2Severity = "low"
	AlertV2SeverityInfo   AlertV2Severity = "info"
)

const (
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
	GetAlertV2PrometheusByID(ctx context.Context, alertID int) (AlertV2Prometheus, error)
	DeleteAlertV2Prometheus(ctx context.Context, alertID int) error
}

type AlertV2EventInterface interface {
	Base
	CreateAlertV2Event(ctx context.Context, alert AlertV2Event) (AlertV2Event, error)
	UpdateAlertV2Event(ctx context.Context, alert AlertV2Event) (AlertV2Event, error)
	GetAlertV2EventByID(ctx context.Context, alertID int) (AlertV2Event, error)
	DeleteAlertV2Event(ctx context.Context, alertID int) error
}

type AlertV2MetricInterface interface {
	Base
	CreateAlertV2Metric(ctx context.Context, alert AlertV2Metric) (AlertV2Metric, error)
	UpdateAlertV2Metric(ctx context.Context, alert AlertV2Metric) (AlertV2Metric, error)
	GetAlertV2MetricByID(ctx context.Context, alertID int) (AlertV2Metric, error)
	DeleteAlertV2Metric(ctx context.Context, alertID int) error
}

type AlertV2ChangeInterface interface {
	Base
	CreateAlertV2Change(ctx context.Context, alert AlertV2Change) (AlertV2Change, error)
	UpdateAlertV2Change(ctx context.Context, alert AlertV2Change) (AlertV2Change, error)
	GetAlertV2ChangeByID(ctx context.Context, alertID int) (AlertV2Change, error)
	DeleteAlertV2Change(ctx context.Context, alertID int) error
}

type AlertV2FormBasedPrometheusInterface interface {
	Base
	CreateAlertV2FormBasedPrometheus(ctx context.Context, alert AlertV2FormBasedPrometheus) (AlertV2FormBasedPrometheus, error)
	UpdateAlertV2FormBasedPrometheus(ctx context.Context, alert AlertV2FormBasedPrometheus) (AlertV2FormBasedPrometheus, error)
	GetAlertV2FormBasedPrometheusByID(ctx context.Context, alertID int) (AlertV2FormBasedPrometheus, error)
	DeleteAlertV2FormBasedPrometheus(ctx context.Context, alertID int) error
}

type AlertV2GroupOutlierInterface interface {
	Base
	CreateAlertV2GroupOutlier(ctx context.Context, alert AlertV2GroupOutlier) (AlertV2GroupOutlier, error)
	UpdateAlertV2GroupOutlier(ctx context.Context, alert AlertV2GroupOutlier) (AlertV2GroupOutlier, error)
	GetAlertV2GroupOutlierByID(ctx context.Context, alertID int) (AlertV2GroupOutlier, error)
	DeleteAlertV2GroupOutlier(ctx context.Context, alertID int) error
}

type AlertV2DowntimeInterface interface {
	Base
	CreateAlertV2Downtime(ctx context.Context, alert AlertV2Downtime) (AlertV2Downtime, error)
	UpdateAlertV2Downtime(ctx context.Context, alert AlertV2Downtime) (AlertV2Downtime, error)
	GetAlertV2DowntimeByID(ctx context.Context, alertID int) (AlertV2Downtime, error)
	DeleteAlertV2Downtime(ctx context.Context, alertID int) error
}

func (c *Client) CreateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (createdAlert AlertV2Prometheus, err error) {
	err = c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	payload, err := Marshal(alertV2PrometheusWrapper{Alert: alert})
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	wrapper, err := createAlertV2AndUnmarshal[alertV2PrometheusWrapper](ctx, c, payload)
	if err != nil {
		return AlertV2Prometheus{}, err
	}
	return wrapper.Alert, nil
}

func (c *Client) UpdateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	payload, err := Marshal(alertV2PrometheusWrapper{Alert: alert})
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	wrapper, err := updateAlertV2AndUnmarshal[alertV2PrometheusWrapper](ctx, c, alert.ID, payload)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) GetAlertV2PrometheusByID(ctx context.Context, alertID int) (AlertV2Prometheus, error) {
	wrapper, err := getAlertV2[alertV2PrometheusWrapper](ctx, c, alertID)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) DeleteAlertV2Prometheus(ctx context.Context, alertID int) error {
	return c.deleteAlertV2(ctx, alertID)
}

func (c *Client) CreateAlertV2Event(ctx context.Context, alert AlertV2Event) (AlertV2Event, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Event{}, err
	}

	err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Event{}, err
	}

	payload, err := Marshal(alertV2EventWrapper{Alert: alert})
	if err != nil {
		return AlertV2Event{}, err
	}

	wrapper, err := createAlertV2AndUnmarshal[alertV2EventWrapper](ctx, c, payload)
	if err != nil {
		return AlertV2Event{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) UpdateAlertV2Event(ctx context.Context, alert AlertV2Event) (AlertV2Event, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Event{}, err
	}

	err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Event{}, err
	}

	payload, err := Marshal(alertV2EventWrapper{Alert: alert})
	if err != nil {
		return AlertV2Event{}, err
	}

	wrapper, err := updateAlertV2AndUnmarshal[alertV2EventWrapper](ctx, c, alert.ID, payload)
	if err != nil {
		return AlertV2Event{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) GetAlertV2EventByID(ctx context.Context, alertID int) (AlertV2Event, error) {
	wrapper, err := getAlertV2[alertV2EventWrapper](ctx, c, alertID)
	if err != nil {
		return AlertV2Event{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) DeleteAlertV2Event(ctx context.Context, alertID int) error {
	return c.deleteAlertV2(ctx, alertID)
}

func (c *Client) CreateAlertV2Metric(ctx context.Context, alert AlertV2Metric) (AlertV2Metric, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Metric{}, err
	}

	err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Metric{}, err
	}

	payload, err := Marshal(alertV2MetricWrapper{Alert: alert})
	if err != nil {
		return AlertV2Metric{}, err
	}

	wrapper, err := createAlertV2AndUnmarshal[alertV2MetricWrapper](ctx, c, payload)
	if err != nil {
		return AlertV2Metric{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) UpdateAlertV2Metric(ctx context.Context, alert AlertV2Metric) (AlertV2Metric, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Metric{}, err
	}

	err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Metric{}, err
	}

	payload, err := Marshal(alertV2MetricWrapper{Alert: alert})
	if err != nil {
		return AlertV2Metric{}, err
	}

	wrapper, err := updateAlertV2AndUnmarshal[alertV2MetricWrapper](ctx, c, alert.ID, payload)
	if err != nil {
		return AlertV2Metric{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) GetAlertV2MetricByID(ctx context.Context, alertID int) (AlertV2Metric, error) {
	wrapper, err := getAlertV2[alertV2MetricWrapper](ctx, c, alertID)
	if err != nil {
		return AlertV2Metric{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) DeleteAlertV2Metric(ctx context.Context, alertID int) error {
	return c.deleteAlertV2(ctx, alertID)
}

func (c *Client) CreateAlertV2Downtime(ctx context.Context, alert AlertV2Downtime) (AlertV2Downtime, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	payload, err := Marshal(alertV2DowntimeWrapper{Alert: alert})
	if err != nil {
		return AlertV2Downtime{}, err
	}

	wrapper, err := createAlertV2AndUnmarshal[alertV2DowntimeWrapper](ctx, c, payload)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	return wrapper.Alert, err
}

func (c *Client) UpdateAlertV2Downtime(ctx context.Context, alert AlertV2Downtime) (AlertV2Downtime, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	payload, err := Marshal(alertV2DowntimeWrapper{Alert: alert})
	if err != nil {
		return AlertV2Downtime{}, err
	}

	wrapper, err := updateAlertV2AndUnmarshal[alertV2DowntimeWrapper](ctx, c, alert.ID, payload)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	return wrapper.Alert, err
}

func (c *Client) GetAlertV2DowntimeByID(ctx context.Context, alertID int) (AlertV2Downtime, error) {
	wrapper, err := getAlertV2[alertV2DowntimeWrapper](ctx, c, alertID)
	if err != nil {
		return AlertV2Downtime{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) DeleteAlertV2Downtime(ctx context.Context, alertID int) error {
	return c.deleteAlertV2(ctx, alertID)
}

func (c *Client) CreateAlertV2Change(ctx context.Context, alert AlertV2Change) (AlertV2Change, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Change{}, err
	}

	err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Change{}, err
	}

	payload, err := Marshal(alertV2ChangeWrapper{Alert: alert})
	if err != nil {
		return AlertV2Change{}, err
	}

	wrapper, err := createAlertV2AndUnmarshal[alertV2ChangeWrapper](ctx, c, payload)
	if err != nil {
		return AlertV2Change{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) UpdateAlertV2Change(ctx context.Context, alert AlertV2Change) (AlertV2Change, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Change{}, err
	}

	err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2Change{}, err
	}

	payload, err := Marshal(alertV2ChangeWrapper{Alert: alert})
	if err != nil {
		return AlertV2Change{}, err
	}

	wrapper, err := updateAlertV2AndUnmarshal[alertV2ChangeWrapper](ctx, c, alert.ID, payload)
	if err != nil {
		return AlertV2Change{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) GetAlertV2ChangeByID(ctx context.Context, alertID int) (AlertV2Change, error) {
	wrapper, err := getAlertV2[alertV2ChangeWrapper](ctx, c, alertID)
	if err != nil {
		return AlertV2Change{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) DeleteAlertV2Change(ctx context.Context, alertID int) error {
	return c.deleteAlertV2(ctx, alertID)
}

func (c *Client) CreateAlertV2FormBasedPrometheus(ctx context.Context, alert AlertV2FormBasedPrometheus) (AlertV2FormBasedPrometheus, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	payload, err := Marshal(alertV2FormBasedPrometheusWrapper{Alert: alert})
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	wrapper, err := createAlertV2AndUnmarshal[alertV2FormBasedPrometheusWrapper](ctx, c, payload)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) UpdateAlertV2FormBasedPrometheus(ctx context.Context, alert AlertV2FormBasedPrometheus) (AlertV2FormBasedPrometheus, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	payload, err := Marshal(alertV2FormBasedPrometheusWrapper{Alert: alert})
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	wrapper, err := updateAlertV2AndUnmarshal[alertV2FormBasedPrometheusWrapper](ctx, c, alert.ID, payload)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) GetAlertV2FormBasedPrometheusByID(ctx context.Context, alertID int) (AlertV2FormBasedPrometheus, error) {
	wrapper, err := getAlertV2[alertV2FormBasedPrometheusWrapper](ctx, c, alertID)
	if err != nil {
		return AlertV2FormBasedPrometheus{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) DeleteAlertV2FormBasedPrometheus(ctx context.Context, alertID int) error {
	return c.deleteAlertV2(ctx, alertID)
}

func (c *Client) CreateAlertV2GroupOutlier(ctx context.Context, alert AlertV2GroupOutlier) (AlertV2GroupOutlier, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	payload, err := Marshal(alertV2GroupOutlierWrapper{Alert: alert})
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	wrapper, err := createAlertV2AndUnmarshal[alertV2GroupOutlierWrapper](ctx, c, payload)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) UpdateAlertV2GroupOutlier(ctx context.Context, alert AlertV2GroupOutlier) (AlertV2GroupOutlier, error) {
	err := c.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	err = c.translateScopeSegmentLabels(ctx, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	payload, err := Marshal(alertV2GroupOutlierWrapper{Alert: alert})
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	wrapper, err := updateAlertV2AndUnmarshal[alertV2GroupOutlierWrapper](ctx, c, alert.ID, payload)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) GetAlertV2GroupOutlierByID(ctx context.Context, alertID int) (AlertV2GroupOutlier, error) {
	wrapper, err := getAlertV2[alertV2GroupOutlierWrapper](ctx, c, alertID)
	if err != nil {
		return AlertV2GroupOutlier{}, err
	}

	return wrapper.Alert, nil
}

func (c *Client) DeleteAlertV2GroupOutlier(ctx context.Context, alertID int) error {
	return c.deleteAlertV2(ctx, alertID)
}

func createAlertV2AndUnmarshal[T any](ctx context.Context, c *Client, alertJSON io.Reader) (value T, err error) {
	var zero T

	response, err := c.requester.Request(ctx, http.MethodPost, c.alertsV2URL(), alertJSON)
	if err != nil {
		return zero, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return zero, c.ErrorFromResponse(response)
	}

	return Unmarshal[T](response.Body)
}

func updateAlertV2AndUnmarshal[T any](ctx context.Context, c *Client, alertID int, alertJSON io.Reader) (value T, err error) {
	var zero T

	response, err := c.requester.Request(ctx, http.MethodPut, c.alertV2URL(alertID), alertJSON)
	if err != nil {
		return zero, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return zero, c.ErrorFromResponse(response)
	}

	return Unmarshal[T](response.Body)
}

func (c *Client) deleteAlertV2(ctx context.Context, alertID int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.alertV2URL(alertID), nil)
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

func getAlertV2[T any](ctx context.Context, c *Client, alertID int) (value T, err error) {
	var zero T
	response, err := c.requester.Request(ctx, http.MethodGet, c.alertV2URL(alertID), nil)
	if err != nil {
		return zero, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode == http.StatusNotFound {
		return zero, ErrAlertV2NotFound
	}
	if response.StatusCode != http.StatusOK {
		return zero, c.ErrorFromResponse(response)
	}

	return Unmarshal[T](response.Body)
}

func (c *Client) addNotificationChannelType(ctx context.Context, notificationChannelConfigList []NotificationChannelConfigV2) error {
	// on put/posts the api wants the type of the channel even if it can be inferred
	for i, n := range notificationChannelConfigList {
		nc, err := c.GetNotificationChannelByID(ctx, n.ChannelID)
		if err != nil {
			return fmt.Errorf("error getting info for notification channel %d: %w", n.ChannelID, err)
		}
		notificationChannelConfigList[i].Type = nc.Type
	}
	return nil
}

func (c *Client) translateScopeSegmentLabels(ctx context.Context, scopedSegmentedConfig *ScopedSegmentedConfig) error {
	// the operand of the scope must be in dot notation
	if scopedSegmentedConfig.Scope != nil {
		for i, e := range scopedSegmentedConfig.Scope.Expressions {
			labelDescriptorV3, err := c.getLabelDescriptor(ctx, e.Operand)
			if err != nil {
				return fmt.Errorf("error getting descriptor for label %s: %w", e.Operand, err)
			}
			scopedSegmentedConfig.Scope.Expressions[i].Operand = labelDescriptorV3.ID
		}
	}

	// the label descriptor id must be in dot notation
	for i, d := range scopedSegmentedConfig.SegmentBy {
		labelDescriptorV3, err := c.getLabelDescriptor(ctx, d.ID)
		if err != nil {
			return fmt.Errorf("error getting descriptor for label %s: %w", d.ID, err)
		}
		scopedSegmentedConfig.SegmentBy[i].ID = labelDescriptorV3.ID
	}

	return nil
}

func (c *Client) getLabelDescriptor(ctx context.Context, label string) (LabelDescriptorV3, error) {
	var alertDescriptor LabelDescriptorV3

	labelCache.Lock()
	defer labelCache.Unlock()

	if len(labelCache.labels) == 0 {
		log.Printf("[DEBUG] GetLabel for %s: fetching all labels", label)
		labelDescriptors, err := c.getLabels(ctx)
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
	return c.buildLabelDescriptor(ctx, label)
}

// buildLabelDescriptor gets the descriptor of a label in public notation from the v3/labels/descriptors api
// this is not a general solution to get the descriptor for a public notation label since custom labels will not be properly translated
func (c *Client) buildLabelDescriptor(ctx context.Context, label string) (descriptor LabelDescriptorV3, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.labelsDescriptorsV3URL(label), nil)
	if err != nil {
		return LabelDescriptorV3{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		err = c.ErrorFromResponse(response)
		return LabelDescriptorV3{}, err
	}

	descriptorWrapper, err := Unmarshal[labelDescriptorV3](response.Body)
	if err != nil {
		return LabelDescriptorV3{}, err
	}

	return descriptorWrapper.LabelDescriptor, nil
}

func (c *Client) getLabels(ctx context.Context) (labels []LabelDescriptorV3, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.labelsV3URL(), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		err = c.ErrorFromResponse(response)
		return nil, err
	}

	wrapper, err := Unmarshal[labelsDescriptorV3](response.Body)
	if err != nil {
		return nil, err
	}
	return wrapper.AllLabels, nil
}

func (c *Client) alertsV2URL() string {
	return fmt.Sprintf(alertsV2Path, c.config.url)
}

func (c *Client) alertV2URL(alertID int) string {
	return fmt.Sprintf(alertV2Path, c.config.url, alertID)
}

func (c *Client) labelsV3URL() string {
	return fmt.Sprintf(labelsV3Path, c.config.url)
}

func (c *Client) labelsDescriptorsV3URL(label string) string {
	return fmt.Sprintf(labelsV3DescriptorsPath, c.config.url, label)
}
