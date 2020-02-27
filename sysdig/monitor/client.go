package monitor

import (
	"io"
	"net/http"
)

type SysdigMonitorClient interface {
	CreateAlert(Alert) (Alert, error)
	DeleteAlert(int) error
	UpdateAlert(Alert) (Alert, error)
	GetAlertById(int) (Alert, error)
}

func NewSysdigMonitorClient(apiToken string, url string) SysdigMonitorClient {
	return &sysdigMonitorClient{
		SysdigMonitorAPIToken: apiToken,
		URL:                   url,
		httpClient:            http.DefaultClient,
	}
}

type sysdigMonitorClient struct {
	SysdigMonitorAPIToken string
	URL                   string
	httpClient            *http.Client
}

func (c *sysdigMonitorClient) doSysdigMonitorRequest(method string, url string, payload io.Reader) (*http.Response, error) {
	request, _ := http.NewRequest(method, url, payload)
	request.Header.Set("Authorization", "Bearer "+c.SysdigMonitorAPIToken)
	request.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(request)
}
