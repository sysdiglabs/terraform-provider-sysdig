package monitor

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor/model"
)

type SysdigMonitorClient interface {
	CreateAlertV2Prometheus(context.Context, AlertV2Prometheus) (AlertV2Prometheus, error)
	DeleteAlertV2Prometheus(context.Context, int) error
	UpdateAlertV2Prometheus(context.Context, AlertV2Prometheus) (AlertV2Prometheus, error)
	GetAlertV2PrometheusById(context.Context, int) (AlertV2Prometheus, error)

	CreateAlertV2Event(context.Context, AlertV2Event) (AlertV2Event, error)
	DeleteAlertV2Event(context.Context, int) error
	UpdateAlertV2Event(context.Context, AlertV2Event) (AlertV2Event, error)
	GetAlertV2EventById(context.Context, int) (AlertV2Event, error)

	CreateAlertV2Metric(context.Context, AlertV2Metric) (AlertV2Metric, error)
	DeleteAlertV2Metric(context.Context, int) error
	UpdateAlertV2Metric(context.Context, AlertV2Metric) (AlertV2Metric, error)
	GetAlertV2MetricById(context.Context, int) (AlertV2Metric, error)

	CreateAlertV2Downtime(context.Context, AlertV2Downtime) (AlertV2Downtime, error)
	DeleteAlertV2Downtime(context.Context, int) error
	UpdateAlertV2Downtime(context.Context, AlertV2Downtime) (AlertV2Downtime, error)
	GetAlertV2DowntimeById(context.Context, int) (AlertV2Downtime, error)

	GetLabelDescriptor(ctx context.Context, label string) (LabelDescriptorV3, error)

	GetNotificationChannelById(context.Context, int) (NotificationChannel, error)

	GetDashboardByID(context.Context, int) (*model.Dashboard, error)
	CreateDashboard(context.Context, *model.Dashboard) (*model.Dashboard, error)
	UpdateDashboard(context.Context, *model.Dashboard) (*model.Dashboard, error)
	DeleteDashboard(context.Context, int) error

	GetCloudAccountById(context.Context, int) (*CloudAccount, error)
	CreateCloudAccount(context.Context, *CloudAccount) (*CloudAccount, error)
	UpdateCloudAccount(context.Context, int, *CloudAccount) (*CloudAccount, error)
	DeleteCloudAccountById(context.Context, int) error
}

func WithExtraHeaders(client SysdigMonitorClient, extraHeaders map[string]string) SysdigMonitorClient {
	rawClient := client.(*sysdigMonitorClient)
	rawClient.extraHeaders = extraHeaders
	return client
}

func NewSysdigMonitorClient(apiToken string, url string, insecure bool) SysdigMonitorClient {
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
			Proxy:           http.ProxyFromEnvironment,
		},
	}

	return &sysdigMonitorClient{
		SysdigMonitorAPIToken: apiToken,
		URL:                   url,
		httpClient:            httpClient.StandardClient(),
	}
}

type sysdigMonitorClient struct {
	SysdigMonitorAPIToken string
	URL                   string
	httpClient            *http.Client
	extraHeaders          map[string]string
}

func (client *sysdigMonitorClient) doSysdigMonitorRequest(ctx context.Context, method string, url string, payload io.Reader) (*http.Response, error) {
	request, _ := http.NewRequest(method, url, payload)
	request = request.WithContext(ctx)
	request.Header.Set("Authorization", "Bearer "+client.SysdigMonitorAPIToken)
	request.Header.Set("Content-Type", "application/json")
	if client.extraHeaders != nil {
		for key, value := range client.extraHeaders {
			request.Header.Set(key, value)
		}
	}

	out, _ := httputil.DumpRequestOut(request, true)
	log.Printf("[DEBUG] %s", string(out))

	response, err := client.httpClient.Do(request)
	if err != nil {
		log.Println(err.Error())
		return response, err
	}

	out, _ = httputil.DumpResponse(response, true)
	log.Printf("[DEBUG] %s", string(out))

	return response, err
}
