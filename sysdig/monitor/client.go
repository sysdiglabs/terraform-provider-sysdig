package monitor

import (
	"crypto/tls"
	"io"
	"net/http"
)

type SysdigMonitorClient interface {
	CreateAlert(Alert) (Alert, error)
	DeleteAlert(int) error
	UpdateAlert(Alert) (Alert, error)
	GetAlertById(int) (Alert, error)
}

func NewSysdigMonitorClient(apiToken string, url string, insecure bool) SysdigMonitorClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}

	return &sysdigMonitorClient{
		SysdigMonitorAPIToken: apiToken,
		URL:                   url,
		httpClient:            httpClient,
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
