package v2

import (
	"bytes"
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

const (
	SysdigTeamIDHeader        = "SysdigTeamID"
	AuthorizationHeader       = "Authorization"
	ContentTypeHeader         = "Content-Type"
	ContentTypeJSON           = "application/json"
	ContentTypeFormURLEncoded = "x-www-form-urlencoded"
)

type Common interface {
	TeamInterface
}

type Requester interface {
	Request(ctx context.Context, method string, url string, payload io.Reader) (*http.Response, error)
}

type Client struct {
	config    *config
	requester Requester
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

func Unmarshal[T any](data io.ReadCloser) (T, error) {
	var result T

	body, err := io.ReadAll(data)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	return result, err
}

func Marshal[T any](data T) (io.Reader, error) {
	payload, err := json.Marshal(data)
	return bytes.NewBuffer(payload), err
}

func request(httpClient *http.Client, cfg *config, request *http.Request) (*http.Response, error) {
	if cfg.extraHeaders != nil {
		for key, value := range cfg.extraHeaders {
			request.Header.Set(key, value)
		}
	}

	out, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s", string(out))
	response, err := httpClient.Do(request)
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

func newHTTPClient(cfg *config) *http.Client {
	httpClient := retryablehttp.NewClient()
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: cfg.insecure}
	httpClient.HTTPClient = &http.Client{Transport: transport}
	return httpClient.StandardClient()
}
