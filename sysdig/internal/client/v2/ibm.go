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
	ibmInstanceIDHeader   = "IBMInstanceID"
	ibmIAMPath            = "/identity/token"
	ibmGrantTypeFormValue = "grant_type"
	ibmAPIKeyFormValue    = "apikey"
	ibmAPIKeyGrantType    = "urn:ibm:params:oauth:grant-type:apikey"
	sysdigTeamIDHeader    = "SysdigTeamID"
	getTeamByNamePath     = "/api/v2/teams/light/name/"
	ibmProductHeader      = "SysdigProduct"
)

type IBMCommon interface {
	Common
}

type IBMMonitor interface {
	IBMCommon
	MonitorCommon
}

type IBMSecure interface {
	IBMCommon
	SecureCommon
}

type (
	IBMAccessToken string
	UnixTimestamp  int64
)

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

func (ir *IBMRequest) getIBMIAMToken() (token IBMAccessToken, err error) {
	ir.tokenLock.Lock()
	defer ir.tokenLock.Unlock()

	if UnixTimestamp(time.Now().Unix()) < ir.tokenExpiration {
		return ir.token, nil
	}

	data := url.Values{}
	data.Set(ibmGrantTypeFormValue, ibmAPIKeyGrantType)
	data.Set(ibmAPIKeyFormValue, ir.config.ibmAPIKey)
	identityURL := fmt.Sprintf("%s%s", ir.config.ibmIamURL, ibmIAMPath)

	r, err := http.NewRequest(http.MethodPost, identityURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	r.Header.Set(ContentTypeHeader, ContentTypeFormURLEncoded)

	response, err := request(ir.httpClient, ir.config, r)
	if err != nil {
		return "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	iamToken, err := Unmarshal[IAMTokenResponse](response.Body)
	if err != nil {
		return "", err
	}

	ir.token = IBMAccessToken(iamToken.AccessToken)
	ir.tokenExpiration = UnixTimestamp(iamToken.Expiration)

	return ir.token, nil
}

func (ir *IBMRequest) getTeamIDByName(ctx context.Context, name string, token IBMAccessToken) (teamID int, err error) {
	r, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s%s%s", ir.config.url, getTeamByNamePath, name),
		nil,
	)
	if err != nil {
		return -1, err
	}

	r = r.WithContext(ctx)
	r.Header.Set(ibmInstanceIDHeader, ir.config.ibmInstanceID)
	r.Header.Set(AuthorizationHeader, fmt.Sprintf("Bearer %s", token))
	r.Header.Set(SysdigProductHeader, ir.config.product)
	r.Header.Set(ibmProductHeader, ir.config.product)

	response, err := request(ir.httpClient, ir.config, r)
	if err != nil {
		return -1, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	wrapper, err := Unmarshal[teamWrapper](response.Body)
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
		ibmInstanceIDHeader: ir.config.ibmInstanceID,
		AuthorizationHeader: fmt.Sprintf("Bearer %s", token),
		SysdigProductHeader: ir.config.product,
		ibmProductHeader:    ir.config.product,
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
	r.Header.Set(ibmInstanceIDHeader, ir.config.ibmInstanceID)
	r.Header.Set(AuthorizationHeader, fmt.Sprintf("Bearer %s", token))
	r.Header.Set(sysdigTeamIDHeader, strconv.Itoa(teamID))
	r.Header.Set(ContentTypeHeader, ContentTypeJSON)
	r.Header.Set(SysdigProviderHeader, SysdigProviderHeaderValue)
	r.Header.Set(SysdigProductHeader, ir.config.product)
	r.Header.Set(ibmProductHeader, ir.config.product)

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

func NewIBMSecure(opts ...ClientOption) IBMSecure {
	return newIBMClient(opts...)
}
