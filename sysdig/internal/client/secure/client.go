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
	request.Header.Set("Sysdig-Provider", "Terraform")
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
