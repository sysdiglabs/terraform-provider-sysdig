package v2

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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
	SysdigTeamIDHeader    = "SysdigTeamID"
	GetTeamByNamePath     = "/api/v2/teams/light/name/"
)

type IBMCommon interface {
	Common
}

type IBMMonitor interface {
	IBMCommon
	MonitorCommon
}

type IBMAccessToken string
type UnixTimestamp int64

type IBMRequest struct {
	config     *config
	httpClient *http.Client

	tokenLock       *sync.Mutex
	tokenExpiration UnixTimestamp
	token           IBMAccessToken

	teamIDLock *sync.Mutex
	teamID     *int
}

type IAMTokenResponse struct {
	AccessToken string `json:"access_token"`
	Expiration  int64  `json:"expiration"`
}

func (ir *IBMRequest) getIBMIAMToken() (IBMAccessToken, error) {
	ir.tokenLock.Lock()
	defer ir.tokenLock.Unlock()

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

func (ir *IBMRequest) getTeamIDByName(ctx context.Context, name string, token IBMAccessToken) (int, error) {
	r, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s%s%s", ir.config.url, GetTeamByNamePath, name),
		nil,
	)
	if err != nil {
		return -1, err
	}

	r = r.WithContext(ctx)
	r.Header.Set(IBMInstanceIDHeader, ir.config.ibmInstanceID)
	r.Header.Set(AuthorizationHeader, fmt.Sprintf("Bearer %s", token))

	resp, err := request(ir.httpClient, ir.config, r)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	wrapper, err := Unmarshal[teamWrapper](resp.Body)
	if err != nil {
		return -1, err
	}

	return wrapper.Team.ID, nil
}

func (ir *IBMRequest) CurrentTeamID(ctx context.Context) (int, error) {
	ir.teamIDLock.Lock()
	defer ir.teamIDLock.Unlock()

	if ir.teamID != nil {
		return *ir.teamID, nil
	}

	token, err := ir.getIBMIAMToken()
	if err != nil {
		return -1, err
	}

	if ir.config.sysdigTeamName != "" {
		teamID, err := ir.getTeamIDByName(ctx, ir.config.sysdigTeamName, token)
		if err != nil {
			return -1, err
		}

		ir.teamID = &teamID
		return *ir.teamID, nil
	}

	// use default current team
	user, err := getMe(ctx, ir.config, ir.httpClient, map[string]string{
		IBMInstanceIDHeader: ir.config.ibmInstanceID,
		AuthorizationHeader: fmt.Sprintf("Bearer %s", token),
	})
	if err != nil {
		return -1, err
	}

	if user.CurrentTeam == nil {
		return -1, errMissingCurrentTeam
	}

	ir.teamID = user.CurrentTeam

	return *ir.teamID, nil
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

	teamID, err := ir.CurrentTeamID(ctx)
	if err != nil {
		return nil, err
	}

	r = r.WithContext(ctx)
	r.Header.Set(IBMInstanceIDHeader, ir.config.ibmInstanceID)
	r.Header.Set(AuthorizationHeader, fmt.Sprintf("Bearer %s", token))
	r.Header.Set(SysdigTeamIDHeader, strconv.Itoa(teamID))
	r.Header.Set(ContentTypeHeader, ContentTypeJSON)
	r.Header.Set(SysdigProviderHeader, SysdigProviderHeaderValue)

	return request(ir.httpClient, ir.config, r)
}

func newIBMClient(opts ...ClientOption) *Client {
	cfg := configure(opts...)
	return &Client{
		config: cfg,
		requester: &IBMRequest{
			tokenLock:  &sync.Mutex{},
			teamIDLock: &sync.Mutex{},
			config:     cfg,
			httpClient: newHTTPClient(cfg),
			teamID:     cfg.sysdigTeamID,
		},
	}
}

func NewIBMMonitor(opts ...ClientOption) IBMMonitor {
	return newIBMClient(opts...)
}
