package v2

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/jmespath/go-jmespath"
	"github.com/spf13/cast"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Common interface {
	TeamInterface
}

type Monitor interface {
	Common
}

type Secure interface {
	Common
}

type RequestFunc func(ctx context.Context, method string, url string, payload io.Reader) (*http.Response, error)

type Client struct {
	config     *config
	httpClient *http.Client
	DoRequest  RequestFunc
}

type config struct {
	url          string
	token        string
	insecure     bool
	extraHeaders map[string]string
}

type ClientOption func(c *config)

func WithURL(url string) ClientOption {
	return func(c *config) {
		c.url = url
	}
}

func WithToken(token string) ClientOption {
	return func(c *config) {
		c.token = token
	}
}

func WithInsecure(insecure bool) ClientOption {
	return func(c *config) {
		c.insecure = insecure
	}
}

func WithExtraHeaders(headers map[string]string) ClientOption {
	return func(c *config) {
		c.extraHeaders = headers
	}
}

func NewMonitor(opts ...ClientOption) Monitor {
	return newClient(opts...)
}

func NewSecure(opts ...ClientOption) Secure {
	return newClient(opts...)
}

func newClient(opts ...ClientOption) *Client {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	httpClient := retryablehttp.NewClient()
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: cfg.insecure}
	httpClient.HTTPClient = &http.Client{Transport: transport}

	c := &Client{
		config:     cfg,
		httpClient: httpClient.StandardClient(),
	}

	// default to sysdig request
	c.DoRequest = c.doSysdigRequest

	return c
}

func (client *Client) request(request *http.Request) (*http.Response, error) {
	if client.config.extraHeaders != nil {
		for key, value := range client.config.extraHeaders {
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

func (client *Client) doSysdigRequest(ctx context.Context, method string, url string, payload io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.token))
	request.Header.Set("Content-Type", "application/json")

	return client.request(request)
}

func (client *Client) ErrorFromResponse(response *http.Response) error {
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

func Unmarshal[T any](body []byte) (T, error) {
	var result T
	err := json.Unmarshal(body, &result)
	return result, err
}

func Marshal[T any](data T) (io.Reader, error) {
	payload, err := json.Marshal(data)
	return bytes.NewBuffer(payload), err
}
