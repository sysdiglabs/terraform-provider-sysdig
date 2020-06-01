package monitor

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

type SysdigMonitorClient interface {
	CreateAlert(Alert) (Alert, error)
	DeleteAlert(int) error
	UpdateAlert(Alert) (Alert, error)
	GetAlertById(int) (Alert, error)

	CreateTeam(Team) (Team, error)
	GetTeamById(int) (Team, error)
	UpdateTeam(Team) (Team, error)
	DeleteTeam(int) error
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

	out, _ := httputil.DumpRequestOut(request, true)
	log.Printf("[DEBUG] %s", string(out))

	response, err := c.httpClient.Do(request)

	out, _ = httputil.DumpResponse(response, true)
	log.Printf("[DEBUG] %s", string(out))

	return response, err
}
