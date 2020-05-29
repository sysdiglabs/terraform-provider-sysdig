package common

import (
	"strings"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

type SysdigCommonClient interface {
	CreateUser(User) (User, error)
	GetUserById(int) (User, error)
	DeleteUser(int) error
	UpdateUser(User) (User, error)

	CreateTeam(Team) (Team, error)
	GetTeamById(int) (Team, error)
	DeleteTeam(int) error
	UpdateTeam(Team) (Team, error)
}

func NewSysdigCommonClient(sysdigAPIToken string, url string, insecure bool) SysdigCommonClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}

	return &sysdigCommonClient{
		SysdigAPIToken: sysdigAPIToken,
		URL:            url,
		httpClient:     httpClient,
	}
}

type sysdigCommonClient struct {
	SysdigAPIToken string
	URL            string
	httpClient     *http.Client
}

func (client *sysdigCommonClient) doSysdigCommonRequest(method string, url string, payload io.Reader) (*http.Response, error) {

	var IBMInstanceID string
	var URL string

	if url[0:10] == "instanceid" {
		result := strings.Split(url, "::")
		IBMInstanceID = strings.Split(result[0], "=")[1]
		URL = strings.Split(result[1], "=")[1]
	} else {
		URL = url
	}

	request, _ := http.NewRequest(method, URL, payload)
	request.Header.Set("Content-Type", "application/json")

	if IBMInstanceID != "" {
		request.Header.Set("Authorization", client.SysdigAPIToken)
		request.Header.Set("IBMInstanceID", IBMInstanceID)
	} else {
		request.Header.Set("Authorization", "Bearer "+client.SysdigAPIToken)
	}

	out, _ := httputil.DumpRequestOut(request, true)
	log.Printf("[DEBUG] %s", string(out))
	response, error := client.httpClient.Do(request)

	out, _ = httputil.DumpResponse(response, true)
	log.Printf("[DEBUG] %s", string(out))
	return response, error
}
