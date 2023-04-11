package v2

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	IBMInstanceIDHeader   = "IBMInstanceID"
	IBMIAMPath            = "/identity/token"
	IBMGrantTypeFormValue = "grant_type"
	IBMApiKeyFormValue    = "apikey"
	IBMAPIKeyGrantType    = "urn:ibm:params:oauth:grant-type:apikey"
)

type IBMCommon interface {
	Common
}

type IBMMonitor interface {
	IBMCommon
}

type IBMAccessToken string
type UnixTimestamp int64

type IBMRequest struct {
	mu              sync.Mutex
	config          *config
	httpClient      *http.Client
	tokenExpiration UnixTimestamp
	token           IBMAccessToken
}

type IAMTokenResponse struct {
	AccessToken string `json:"access_token"`
	Expiration  int64  `json:"expiration"`
}

func (ir *IBMRequest) getIBMIAMToken() (IBMAccessToken, error) {
	ir.mu.Lock()
	defer ir.mu.Unlock()

	if UnixTimestamp(time.Now().Unix()) < ir.tokenExpiration {
		return ir.token, nil
	}

	data := url.Values{}
	data.Set(IBMGrantTypeFormValue, IBMAPIKeyGrantType)
	data.Set(IBMApiKeyFormValue, ir.config.ibmAPIKey)
	identityURL := fmt.Sprintf("%s%s", ir.config.ibmIamURL, IBMIAMPath)

	r, err := http.NewRequest(http.MethodPost, identityURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	r.Header.Set(ContentTypeHeader, ContentTypeFormURLEncoded)

	resp, err := request(ir.httpClient, ir.config, r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	iamToken, err := Unmarshal[IAMTokenResponse](resp.Body)
	if err != nil {
		return "", err
	}

	ir.token = IBMAccessToken(iamToken.AccessToken)
	ir.tokenExpiration = UnixTimestamp(iamToken.Expiration)

	return ir.token, nil
}

func (ir *IBMRequest) Request(ctx context.Context, method string, url string, payload io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	token, err := ir.getIBMIAMToken()
	if err != nil {
		return nil, err
	}

	r = r.WithContext(ctx)
	r.Header.Set(IBMInstanceIDHeader, ir.config.ibmInstanceID)
	r.Header.Set(AuthorizationHeader, fmt.Sprintf("Bearer %s", token))
	r.Header.Set(ContentTypeHeader, ContentTypeJSON)

	return request(ir.httpClient, ir.config, r)
}

func newIBMClient(opts ...ClientOption) *Client {
	cfg := configure(opts...)
	return &Client{
		config: cfg,
		requester: &IBMRequest{
			mu:         sync.Mutex{},
			config:     cfg,
			httpClient: newHTTPClient(cfg),
		},
	}
}

func NewIBMMonitor(opts ...ClientOption) IBMMonitor {
	return newIBMClient(opts...)
}
