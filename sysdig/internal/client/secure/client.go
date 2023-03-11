package secure

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/hashicorp/go-retryablehttp"
)

type SysdigSecureClient interface {
	CreatePolicy(context.Context, Policy) (Policy, error)
	DeletePolicy(context.Context, int) error
	UpdatePolicy(context.Context, Policy) (Policy, error)
	GetPolicyById(context.Context, int) (Policy, error)

	CreateRule(context.Context, Rule) (Rule, error)
	GetRuleByID(context.Context, int) (Rule, error)
	UpdateRule(context.Context, Rule) (Rule, error)
	DeleteRule(context.Context, int) error

	CreateNotificationChannel(context.Context, NotificationChannel) (NotificationChannel, error)
	GetNotificationChannelById(context.Context, int) (NotificationChannel, error)
	GetNotificationChannelByName(context.Context, string) (NotificationChannel, error)
	DeleteNotificationChannel(context.Context, int) error
	UpdateNotificationChannel(context.Context, NotificationChannel) (NotificationChannel, error)

	CreateTeam(context.Context, Team) (Team, error)
	GetTeamById(context.Context, int) (Team, error)
	DeleteTeam(context.Context, int) error
	UpdateTeam(context.Context, Team) (Team, error)

	CreateList(context.Context, List) (List, error)
	GetListById(context.Context, int) (List, error)
	DeleteList(context.Context, int) error
	UpdateList(context.Context, List) (List, error)

	CreateMacro(context.Context, Macro) (Macro, error)
	GetMacroById(context.Context, int) (Macro, error)
	DeleteMacro(context.Context, int) error
	UpdateMacro(context.Context, Macro) (Macro, error)

	CreateVulnerabilityExceptionList(context.Context, *VulnerabilityExceptionList) (*VulnerabilityExceptionList, error)
	GetVulnerabilityExceptionListByID(context.Context, string) (*VulnerabilityExceptionList, error)
	DeleteVulnerabilityExceptionList(context.Context, string) error
	UpdateVulnerabilityExceptionList(context.Context, *VulnerabilityExceptionList) (*VulnerabilityExceptionList, error)

	CreateVulnerabilityException(context.Context, string, *VulnerabilityException) (*VulnerabilityException, error)
	GetVulnerabilityExceptionByID(context.Context, string, string) (*VulnerabilityException, error)
	DeleteVulnerabilityException(context.Context, string, string) error
	UpdateVulnerabilityException(context.Context, string, *VulnerabilityException) (*VulnerabilityException, error)

	CreateCloudAccount(context.Context, *CloudAccount) (*CloudAccount, error)
	GetCloudAccountById(context.Context, string) (*CloudAccount, error)
	DeleteCloudAccount(context.Context, string) error
	UpdateCloudAccount(context.Context, string, *CloudAccount) (*CloudAccount, error)
	GetTrustedCloudIdentity(context.Context, string) (string, error)

	CreateBenchmarkTask(context.Context, *BenchmarkTask) (*BenchmarkTask, error)
	GetBenchmarkTask(context.Context, string) (*BenchmarkTask, error)
	DeleteBenchmarkTask(context.Context, string) error
	SetBenchmarkTaskEnabled(context.Context, string, bool) error

	CreateScanningPolicy(context.Context, ScanningPolicy) (ScanningPolicy, error)
	GetScanningPolicyById(context.Context, string) (ScanningPolicy, error)
	DeleteScanningPolicyById(context.Context, string) error
	UpdateScanningPolicyById(context.Context, ScanningPolicy) (ScanningPolicy, error)

	CreateScanningPolicyAssignmentList(context.Context, ScanningPolicyAssignmentList) (ScanningPolicyAssignmentList, error)
	GetScanningPolicyAssignmentList(context.Context) (ScanningPolicyAssignmentList, error)
	DeleteScanningPolicyAssignmentList(context.Context, ScanningPolicyAssignmentList) error
}

func WithExtraHeaders(client SysdigSecureClient, extraHeaders map[string]string) SysdigSecureClient {
	rawClient := client.(*sysdigSecureClient)
	rawClient.extraHeaders = extraHeaders
	return client
}

func NewSysdigSecureClient(sysdigSecureAPIToken string, url string, insecure bool) SysdigSecureClient {
	httpClient := retryablehttp.NewClient()
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}
	httpClient.HTTPClient = &http.Client{Transport: transport}

	return &sysdigSecureClient{
		SysdigSecureAPIToken: sysdigSecureAPIToken,
		URL:                  url,
		httpClient:           httpClient.StandardClient(),
	}
}

type sysdigSecureClient struct {
	SysdigSecureAPIToken string
	URL                  string
	httpClient           *http.Client
	extraHeaders         map[string]string
}

func (client *sysdigSecureClient) doSysdigSecureRequest(ctx context.Context, method string, url string, payload io.Reader) (*http.Response, error) {
	request, _ := http.NewRequest(method, url, payload)
	request = request.WithContext(ctx)
	request.Header.Set("Authorization", "Bearer "+client.SysdigSecureAPIToken)
	request.Header.Set("Content-Type", "application/json")
	if client.extraHeaders != nil {
		for key, value := range client.extraHeaders {
			request.Header.Set(key, value)
		}
	}

	out, _ := httputil.DumpRequestOut(request, true)
	log.Printf("[DEBUG] %s", string(out))
	response, err := client.httpClient.Do(request)
	if err != nil {
		log.Println(err.Error())
		return response, err
	}

	out, _ = httputil.DumpResponse(response, true)
	log.Printf("[DEBUG] %s", string(out))
	return response, err
}
