package common

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/hashicorp/go-retryablehttp"
)

type SysdigCommonClient interface {
	CreateUser(context.Context, *User) (*User, error)
	GetUserById(context.Context, int) (*User, error)
	GetUserByEmail(context.Context, string) (*User, error)
	DeleteUser(context.Context, int) error
	UpdateUser(context.Context, *User) (*User, error)
	GetCurrentUser(context.Context) (*User, error)
}

func WithExtraHeaders(client SysdigCommonClient, extraHeaders map[string]string) SysdigCommonClient {
	rawClient := client.(*sysdigCommonClient)
	rawClient.extraHeaders = extraHeaders
	return client
}

func NewSysdigCommonClient(sysdigAPIToken string, url string, insecure bool) SysdigCommonClient {
	client := retryablehttp.NewClient()
    transport := http.DefaultTransport.(*http.Transport).Clone()
    transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}

    client.HTTPClient = &http.Client{
        Transport: transport,
    }

	return &sysdigCommonClient{
		SysdigAPIToken: sysdigAPIToken,
		URL:            url,
		httpClient:     client.StandardClient(),
	}
}

type sysdigCommonClient struct {
	SysdigAPIToken string
	URL            string
	httpClient     *http.Client
	extraHeaders   map[string]string
}

func (client *sysdigCommonClient) doSysdigCommonRequest(ctx context.Context, method string, url string, payload io.Reader) (*http.Response, error) {
	request, _ := http.NewRequest(method, url, payload)
	request = request.WithContext(ctx)
	request.Header.Set("Authorization", "Bearer "+client.SysdigAPIToken)
	request.Header.Set("Content-Type", "application/json")
	if client.extraHeaders != nil {
		for key, value := range client.extraHeaders {
			request.Header.Set(key, value)
		}
	}

	out, _ := httputil.DumpRequestOut(request, true)
	log.Printf("[DEBUG] %s", string(out))
	response, error := client.httpClient.Do(request)

	out, _ = httputil.DumpResponse(response, true)
	log.Printf("[DEBUG] %s", string(out))
	return response, error
}
