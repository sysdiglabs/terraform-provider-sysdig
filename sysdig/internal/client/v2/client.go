package v2

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/draios/terraform-provider-sysdig/buildinfo"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/jmespath/go-jmespath"
	"github.com/spf13/cast"
)

const (
	GetMePath                  = "/api/users/me"
	UserAgentHeader            = "User-Agent"
	AuthorizationHeader        = "Authorization"
	ContentTypeHeader          = "Content-Type"
	SysdigProviderHeader       = "Sysdig-Provider"
	SysdigProviderHeaderValue  = "Terraform"
	SysdigUserAgentHeaderValue = "SysdigTerraform"
	ContentTypeJSON            = "application/json"
	ContentTypeFormURLEncoded  = "x-www-form-urlencoded"
	SysdigProductHeader        = "X-Sysdig-Product"
)

var errMissingCurrentTeam = errors.New("missing user's current team")

type Base interface {
	CurrentTeamID(ctx context.Context) (int, error)
}

type Common interface {
	UserInterface
	TeamInterface
	NotificationChannelInterface
	IdentityContextInterface
	AgentAccessKeyInterface
}

type MonitorCommon interface {
	AlertInterface
	AlertV2Interface
	DashboardInterface
	SilenceRuleInterface
	InhibitionRuleInterface
}

type SecureCommon interface {
	PosturePolicyInterface
	PostureZoneInterface
	PostureControlInterface
	PostureAcceptRiskInterface
	PostureVulnerabilityAcceptRiskInterface
	ZoneInterface
}

type Requester interface {
	CurrentTeamID(ctx context.Context) (int, error)
	Request(ctx context.Context, method string, url string, payload io.Reader) (*http.Response, error)
}

type Client struct {
	config    *config
	requester Requester
}

func (c *Client) ErrorFromResponse(response *http.Response) error {
	var data any
	err := json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return errors.New(response.Status)
	}

	search, err := jmespath.Search("[message, error, details[], errors[].[reason, message]][][] | join(', ', @)", data)
	if err != nil {
		return errors.New(response.Status)
	}

	if searchArray, ok := search.([]any); ok {
		return errors.New(strings.Join(cast.ToStringSlice(searchArray), ", "))
	}

	searchString := cast.ToString(search)
	if searchString != "" {
		return errors.New(searchString)
	}

	return errors.New(response.Status)
}

func Unmarshal[T any](data io.Reader) (T, error) {
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
	request.Header.Set(UserAgentHeader, fmt.Sprintf("%s/%s", SysdigUserAgentHeaderValue, buildinfo.Version))

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

func getMe(ctx context.Context, cfg *config, httpClient *http.Client, headers map[string]string) (user *User, err error) {
	r, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s%s", cfg.url, GetMePath),
		nil,
	)
	if err != nil {
		return nil, err
	}

	r = r.WithContext(ctx)
	for k, v := range headers {
		r.Header.Set(k, v)
	}

	resp, err := request(httpClient, cfg, r)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := resp.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	wrapper, err := Unmarshal[userWrapper](resp.Body)
	if err != nil {
		return nil, err
	}

	return &wrapper.User, nil
}

func (c *Client) CurrentTeamID(ctx context.Context) (int, error) {
	return c.requester.CurrentTeamID(ctx)
}

func newHTTPClient(cfg *config) *http.Client {
	httpClient := retryablehttp.NewClient()
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: cfg.insecure}
	httpClient.HTTPClient = &http.Client{Transport: transport}
	return httpClient.StandardClient()
}
