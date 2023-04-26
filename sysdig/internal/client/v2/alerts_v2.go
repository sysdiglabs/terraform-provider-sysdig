package v2

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const (
	GetLabelsDescriptorsV3Path = "%s/api/v3/labels/descriptors/%s"
	GetLabelsV3Path            = "%s/api/v3/labels"
	CreateAlertV2Path          = "%s/api/v2/alerts"
	UpdateAlertV2Path          = "%s/api/v2/alerts/%d"
	GetAlertV2Path             = "%s/api/v2/alerts/%d"
	DeleteAlertV2Path          = "%s/api/v2/alerts/%d"
)

var labelCache struct {
	sync.Mutex

	labels []LabelDescriptorV3
}

type AlertsV2Interface interface {
	AlertsV2PrometheusInterface
}

type AlertsV2PrometheusInterface interface {
	CreateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error)
	UpdateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error)
	GetAlertV2PrometheusByID(ctx context.Context, alertID int) (AlertV2Prometheus, error)
	DeleteAlertV2Prometheus(ctx context.Context, alertID int) (err error)
}

func (client *Client) CreateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	payload, err := Marshal(AlertV2PrometheusWrapper{Alert: alert})
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	body, err := client.createAlertV2(ctx, payload)
	if err != nil {
		return AlertV2Prometheus{}, nil
	}

	wrapper, err := Unmarshal[AlertV2PrometheusWrapper](body)
	if err != nil {
		return AlertV2Prometheus{}, nil
	}

	return wrapper.Alert, nil
}

func (client *Client) UpdateAlertV2Prometheus(ctx context.Context, alert AlertV2Prometheus) (AlertV2Prometheus, error) {
	err := client.addNotificationChannelType(ctx, alert.NotificationChannelConfigList)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	payload, err := Marshal(AlertV2PrometheusWrapper{Alert: alert})
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	body, err := client.updateAlertV2(ctx, alert.ID, payload)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	wrapper, err := Unmarshal[AlertV2PrometheusWrapper](body)
	if err != nil {
		return AlertV2Prometheus{}, nil
	}

	return wrapper.Alert, nil
}

func (client *Client) GetAlertV2PrometheusByID(ctx context.Context, alertID int) (AlertV2Prometheus, error) {
	body, err := client.getAlertV2ByID(ctx, alertID)
	if err != nil {
		return AlertV2Prometheus{}, err
	}

	wrapper, err := Unmarshal[AlertV2PrometheusWrapper](body)
	if err != nil {
		return AlertV2Prometheus{}, nil
	}

	return wrapper.Alert, nil
}

func (client *Client) DeleteAlertV2Prometheus(ctx context.Context, alertID int) (err error) {
	return client.deleteAlertV2(ctx, alertID)
}

func (client *Client) createAlertV2(ctx context.Context, payload io.Reader) (io.ReadCloser, error) {
	response, err := client.requester.Request(ctx, http.MethodPost, client.CreateAlertV2URL(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return response.Body, nil
}

func (client *Client) updateAlertV2(ctx context.Context, alertID int, payload io.Reader) (io.ReadCloser, error) {
	response, err := client.requester.Request(ctx, http.MethodPut, client.UpdateAlertV2URL(alertID), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return response.Body, nil
}

func (client *Client) getAlertV2ByID(ctx context.Context, alertID int) (io.ReadCloser, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetAlertV2URL(alertID), nil)
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
	response, err := client.requester.Request(ctx, http.MethodDelete, client.DeleteAlertV2URL(alertID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
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

// getLabelDescriptor gets the descriptor from a label in public notation
func (client *Client) getLabelDescriptor(ctx context.Context, label string) (LabelDescriptorV3, error) {
	var alertDescriptor LabelDescriptorV3

	labelCache.Lock()
	defer labelCache.Unlock()

	if len(labelCache.labels) == 0 {
		log.Printf("[DEBUG] getLabelDescriptor for %s: fetching all labels", label)
		labelDescriptors, err := client.getLabels(ctx)
		if err != nil {
			return alertDescriptor, err
		}
		labelCache.labels = labelDescriptors
	}

	for _, l := range labelCache.labels {
		if l.PublicID == label {
			return l, nil
		}
	}

	// if the label did not exist, build the descriptor from /v3/labels/descriptor
	log.Printf("[DEBUG] getLabelDescriptor for %s: not found in existing customer labels", label)
	return client.buildLabelDescriptor(ctx, label)
}

// buildLabelDescriptor gets the descriptor of a label in public notation from the v3/labels/descriptors api
// this is not a general solution to get the descriptor for a public notation label since custom labels will not be properly translated
// e.g. the public notation cloud_provider_tag_k8s_io_role_master will not be translated to the correct cloudProvider.tag.k8s.io/role/master id
func (client *Client) buildLabelDescriptor(ctx context.Context, label string) (LabelDescriptorV3, error) {
	// always returns 200, even if the label does not exist for the customer
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetLabelsDescriptorsV3URL(label), nil)
	if err != nil {
		return LabelDescriptorV3{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return LabelDescriptorV3{}, client.ErrorFromResponse(response)
	}

	wrapper, err := Unmarshal[labelsDescriptorsV3Wrapper](response.Body)
	if err != nil {
		return LabelDescriptorV3{}, err
	}

	return wrapper.LabelDescriptorV3, nil
}

func (client *Client) getLabels(ctx context.Context) ([]LabelDescriptorV3, error) {
	limit := 6000
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetLabelsV3URL(&limit), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	wrapper, err := Unmarshal[labelsWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return wrapper.AllLabels, err
}

func (client *Client) GetLabelsDescriptorsV3URL(label string) string {
	return fmt.Sprintf(GetLabelsDescriptorsV3Path, client.config.url, label)
}

func (client *Client) GetLabelsV3URL(limit *int) string {
	u := fmt.Sprintf(GetLabelsV3Path, client.config.url)
	if limit != nil {
		u = fmt.Sprintf("%s?limit=%d", u, limit)
	}
	return u
}

func (client *Client) CreateAlertV2URL() string {
	return fmt.Sprintf(CreateAlertV2Path, client.config.url)
}

func (client *Client) UpdateAlertV2URL(alertID int) string {
	return fmt.Sprintf(UpdateAlertV2Path, client.config.url, alertID)
}

func (client *Client) GetAlertV2URL(alertID int) string {
	return fmt.Sprintf(GetAlertV2Path, client.config.url, alertID)
}

func (client *Client) DeleteAlertV2URL(alertID int) string {
	return fmt.Sprintf(DeleteAlertV2Path, client.config.url, alertID)
}
