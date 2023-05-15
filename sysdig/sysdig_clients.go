package sysdig

import (
	"errors"
	"fmt"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SysdigClients interface {
	GetClientType() ClientType
	GetSecureEndpoint() (string, error)
	GetSecureApiToken() (string, error)

	// v2
	sysdigMonitorClientV2() (v2.SysdigMonitor, error)
	sysdigSecureClientV2() (v2.SysdigSecure, error)
	ibmMonitorClient() (v2.IBMMonitor, error)
	commonClientV2() (v2.Common, error)
	sysdigCommonClientV2() (v2.SysdigCommon, error)
}

//go:generate stringer -type ClientType
type ClientType int

const (
	SysdigMonitor ClientType = iota
	SysdigSecure
	IBMMonitor
)

type sysdigClients struct {
	d        *schema.ResourceData
	mu       sync.Mutex
	commonMu sync.Mutex

	// v2
	monitorClientV2  v2.SysdigMonitor
	secureClientV2   v2.SysdigSecure
	monitorIBMClient v2.IBMMonitor
	commonV2         v2.Common
	sysdigCommonV2   v2.SysdigCommon
}

type globalVariables struct {
	apiURL       string
	insecure     bool
	extraHeaders map[string]string
}

type sysdigVariables struct {
	*globalVariables
	token string
}

type ibmVariables struct {
	*globalVariables
	iamURL         string
	instanceID     string
	apiKey         string
	sysdigTeamName string
	sysdigTeamID   *int
}

func getSysdigMonitorVariables(data *schema.ResourceData) (*sysdigVariables, error) {
	var ok bool
	var apiURL, token interface{}

	if apiURL, ok = data.GetOk("sysdig_monitor_url"); !ok {
		return nil, errors.New("missing sysdig monitor URL")
	}

	if token, ok = data.GetOk("sysdig_monitor_api_token"); !ok {
		return nil, errors.New("missing sysdig monitor token")
	}

	return &sysdigVariables{
		globalVariables: &globalVariables{
			apiURL:       apiURL.(string),
			insecure:     data.Get("sysdig_monitor_insecure_tls").(bool),
			extraHeaders: getExtraHeaders(data),
		},
		token: token.(string),
	}, nil
}

func getSysdigSecureVariables(data *schema.ResourceData) (*sysdigVariables, error) {
	var ok bool
	var apiURL, token interface{}

	if apiURL, ok = data.GetOk("sysdig_secure_url"); !ok {
		return nil, errors.New("missing sysdig secure URL")
	}

	if token, ok = data.GetOk("sysdig_secure_api_token"); !ok {
		return nil, errors.New("missing sysdig monitor token")
	}

	return &sysdigVariables{
		globalVariables: &globalVariables{
			apiURL:       apiURL.(string),
			insecure:     data.Get("sysdig_secure_insecure_tls").(bool),
			extraHeaders: getExtraHeaders(data),
		},
		token: token.(string),
	}, nil
}

func getIBMMonitorVariables(data *schema.ResourceData) (*ibmVariables, error) {
	var ok bool
	var apiURL, iamURL, instanceID, apiKey interface{}
	var teamID *int

	if apiURL, ok = data.GetOk("sysdig_monitor_url"); !ok {
		return nil, errors.New("missing monitor IBM URL")
	}

	if iamURL, ok = data.GetOk("ibm_monitor_iam_url"); !ok {
		return nil, errors.New("missing monitor IBM IAM URL")
	}

	if instanceID, ok = data.GetOk("ibm_monitor_instance_id"); !ok {
		return nil, errors.New("missing monitor IBM instance ID")
	}

	if apiKey, ok = data.GetOk("ibm_monitor_api_key"); !ok {
		return nil, errors.New("missing monitor IBM API key")
	}

	if id, ok := data.GetOk("sysdig_monitor_team_id"); ok {
		tmp := id.(int)
		teamID = &tmp
	}

	return &ibmVariables{
		globalVariables: &globalVariables{
			apiURL:       apiURL.(string),
			insecure:     data.Get("sysdig_monitor_insecure_tls").(bool),
			extraHeaders: getExtraHeaders(data),
		},
		iamURL:         iamURL.(string),
		instanceID:     instanceID.(string),
		apiKey:         apiKey.(string),
		sysdigTeamID:   teamID,
		sysdigTeamName: data.Get("sysdig_monitor_team_name").(string),
	}, nil
}

func (c *sysdigClients) GetSecureEndpoint() (string, error) {
	endpoint := c.d.Get("sysdig_secure_url").(string)
	if endpoint == "" {
		return "", errors.New("GetSecureEndpoint, sysdig_secure_url not provided")
	}
	return endpoint, nil
}

func (c *sysdigClients) GetSecureApiToken() (string, error) {
	secureAPIToken := c.d.Get("sysdig_secure_api_token").(string)
	if secureAPIToken == "" {
		return "", errors.New("GetSecureApiToken, sysdig secure token not provided")
	}
	return secureAPIToken, nil
}

func (c *sysdigClients) sysdigMonitorClientV2() (v2.SysdigMonitor, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.monitorClientV2 != nil {
		return c.monitorClientV2, nil
	}

	vars, err := getSysdigMonitorVariables(c.d)
	if err != nil {
		return nil, err
	}

	c.monitorClientV2 = v2.NewSysdigMonitor(
		v2.WithToken(vars.token),
		v2.WithURL(vars.apiURL),
		v2.WithInsecure(vars.insecure),
		v2.WithExtraHeaders(vars.extraHeaders),
	)

	return c.monitorClientV2, nil
}

func (c *sysdigClients) sysdigSecureClientV2() (v2.SysdigSecure, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.secureClientV2 != nil {
		return c.secureClientV2, nil
	}

	vars, err := getSysdigSecureVariables(c.d)
	if err != nil {
		return nil, err
	}

	c.secureClientV2 = v2.NewSysdigSecure(
		v2.WithToken(vars.token),
		v2.WithURL(vars.apiURL),
		v2.WithInsecure(vars.insecure),
		v2.WithExtraHeaders(vars.extraHeaders),
	)

	return c.secureClientV2, nil
}

func (c *sysdigClients) ibmMonitorClient() (v2.IBMMonitor, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.monitorIBMClient != nil {
		return c.monitorIBMClient, nil
	}

	vars, err := getIBMMonitorVariables(c.d)
	if err != nil {
		return nil, err
	}

	c.monitorIBMClient = v2.NewIBMMonitor(
		v2.WithURL(vars.apiURL),
		v2.WithIBMIamURL(vars.iamURL),
		v2.WithIBMInstanceID(vars.instanceID),
		v2.WithIBMAPIKey(vars.apiKey),
		v2.WithInsecure(vars.insecure),
		v2.WithSysdigTeamID(vars.sysdigTeamID),
		v2.WithSysdigTeamName(vars.sysdigTeamName),
	)

	return c.monitorIBMClient, nil
}

func (c *sysdigClients) sysdigCommonClientV2() (v2.SysdigCommon, error) {
	c.commonMu.Lock()
	defer c.commonMu.Unlock()

	var err error
	clientType := c.GetClientType()

	switch clientType {
	case SysdigMonitor:
		c.sysdigCommonV2, err = c.sysdigMonitorClientV2()
	case SysdigSecure:
		c.sysdigCommonV2, err = c.sysdigSecureClientV2()
	default:
		return nil, fmt.Errorf("failed to create common sysdig client, %s is not supported", clientType)
	}

	return c.sysdigCommonV2, err
}

func (c *sysdigClients) commonClientV2() (v2.Common, error) {
	c.commonMu.Lock()
	defer c.commonMu.Unlock()

	var err error
	clientType := c.GetClientType()

	switch clientType {
	case SysdigMonitor:
		c.commonV2, err = c.sysdigMonitorClientV2()
	case SysdigSecure:
		c.commonV2, err = c.sysdigSecureClientV2()
	case IBMMonitor:
		c.commonV2, err = c.ibmMonitorClient()
	default:
		return nil, fmt.Errorf("failed to create common client, %s is not supported", clientType)
	}

	return c.commonV2, err
}

func (c *sysdigClients) GetClientType() ClientType {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := getIBMMonitorVariables(c.d); err == nil {
		return IBMMonitor
	}

	if _, err := getSysdigMonitorVariables(c.d); err == nil {
		return SysdigMonitor
	}

	return SysdigSecure
}

func getExtraHeaders(d *schema.ResourceData) map[string]string {
	if headers, ok := d.GetOk("extra_headers"); ok {
		extraHeaders := headers.(map[string]interface{})
		extraHeadersTransformed := map[string]string{}
		for key := range extraHeaders {
			extraHeadersTransformed[key] = extraHeaders[key].(string)
		}
		return extraHeadersTransformed
	}
	return nil
}
