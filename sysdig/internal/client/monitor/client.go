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
	CreateAlert(context.Context, Alert) (Alert, error)
	DeleteAlert(context.Context, int) error
	UpdateAlert(context.Context, Alert) (Alert, error)
	GetAlertById(context.Context, int) (Alert, error)

	CreateTeam(context.Context, Team) (Team, error)
	GetTeamById(context.Context, int) (Team, error)
	UpdateTeam(context.Context, Team) (Team, error)
	DeleteTeam(context.Context, int) error

	CreateNotificationChannel(context.Context, NotificationChannel) (NotificationChannel, error)
	GetNotificationChannelById(context.Context, int) (NotificationChannel, error)
	GetNotificationChannelByName(context.Context, string) (NotificationChannel, error)
	DeleteNotificationChannel(context.Context, int) error
	UpdateNotificationChannel(context.Context, NotificationChannel) (NotificationChannel, error)

	GetDashboardByID(context.Context, int) (*model.Dashboard, error)
	CreateDashboard(context.Context, *model.Dashboard) (*model.Dashboard, error)
	UpdateDashboard(context.Context, *model.Dashboard) (*model.Dashboard, error)
	DeleteDashboard(context.Context, int) error
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
