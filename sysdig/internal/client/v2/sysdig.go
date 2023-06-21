package v2

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type SysdigRequest struct {
	config     *config
	httpClient *http.Client

	teamIDLock *sync.Mutex
	teamID     *int
}

type SysdigCommon interface {
	Common
	GroupMappingInterface
	GroupMappingConfigInterface
}

type SysdigMonitor interface {
	SysdigCommon
	MonitorCommon
	DashboardInterface
	CloudAccountMonitorInterface
	AlertV2Interface
}

type SysdigSecure interface {
	SysdigCommon
	SecureCommon
	PolicyInterface
	RuleInterface
	ListInterface
	MacroInterface
	ScanningPolicyInterface
	ScanningPolicyAssignmentInterface
	VulnerabilityExceptionListInterface
	VulnerabilityExceptionInterface
	CloudAccountSecureInterface
}

func (sr *SysdigRequest) Request(ctx context.Context, method string, url string, payload io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	r = r.WithContext(ctx)
	r.Header.Set(AuthorizationHeader, fmt.Sprintf("Bearer %s", sr.config.token))
	r.Header.Set(ContentTypeHeader, ContentTypeJSON)
	r.Header.Set(SysdigProviderHeader, SysdigProviderHeaderValue)

	return request(sr.httpClient, sr.config, r)
}

func NewSysdigMonitor(opts ...ClientOption) SysdigMonitor {
	return newSysdigClient(opts...)
}

func NewSysdigSecure(opts ...ClientOption) SysdigSecure {
	return newSysdigClient(opts...)
}

func newSysdigClient(opts ...ClientOption) *Client {
	cfg := configure(opts...)
	return &Client{
		config: cfg,
		requester: &SysdigRequest{
			teamIDLock: &sync.Mutex{},
			config:     cfg,
			httpClient: newHTTPClient(cfg),
		},
	}
}

func (sr *SysdigRequest) CurrentTeamID(ctx context.Context) (int, error) {
	sr.teamIDLock.Lock()
	defer sr.teamIDLock.Unlock()

	if sr.teamID != nil {
		return *sr.teamID, nil
	}

	user, err := getMe(ctx, sr.config, sr.httpClient, map[string]string{
		AuthorizationHeader: fmt.Sprintf("Bearer %s", sr.config.token),
	})
	if err != nil {
		return -1, err
	}

	if user.CurrentTeam == nil {
		return -1, errMissingCurrentTeam
	}

	sr.teamID = user.CurrentTeam

	return *sr.teamID, nil
}
