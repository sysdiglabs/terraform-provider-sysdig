package common

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/jmespath/go-jmespath"
	"github.com/spf13/cast"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Client struct {
	APIToken     string
	URL          string
	httpClient   *http.Client
	extraHeaders map[string]string
}

func NewClient(token string, url string, insecure bool) *Client {
	httpClient := retryablehttp.NewClient()
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}
	httpClient.HTTPClient = &http.Client{Transport: transport}

	return &Client{
		APIToken:   token,
		URL:        url,
		httpClient: httpClient.StandardClient(),
	}
}

func (client *Client) DoSysdigRequest(ctx context.Context, method string, url string, payload io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	request.Header.Set("Authorization", "Bearer "+client.APIToken)
	request.Header.Set("Content-Type", "application/json")
	if client.extraHeaders != nil {
		for key, value := range client.extraHeaders {
			request.Header.Set(key, value)
		}
	}

	out, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s", string(out))
	response, err := client.httpClient.Do(request)
	if err != nil {
		log.Println(err.Error())
		return response, err
	}

	out, err = httputil.DumpResponse(response, true)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s", string(out))
	return response, err
}

func (client *Client) SetExtraHeaders(extraHeaders map[string]string) {
	client.extraHeaders = extraHeaders
}

func ErrorFromResponse(response *http.Response) error {
	var data interface{}
	err := json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return errors.New(response.Status)
	}

	search, err := jmespath.Search("[message, errors[].[reason, message]][][] | join(', ', @)", data)
	if err != nil {
		return errors.New(response.Status)
	}

	if searchArray, ok := search.([]interface{}); ok {
		return errors.New(strings.Join(cast.ToStringSlice(searchArray), ", "))
	}

	return errors.New(cast.ToString(search))
}
