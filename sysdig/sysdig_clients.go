package sysdig

import (
	"errors"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/common"
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor"
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
)

type SysdigClients interface {
	GetSecureEndpoint() (string, error)
	GetSecureApiToken() (string, error)

	sysdigMonitorClient() (monitor.SysdigMonitorClient, error)
	sysdigSecureClient() (secure.SysdigSecureClient, error)
	sysdigCommonClient() (common.SysdigCommonClient, error)

	// v2
	sysdigMonitorClientV2() (v2.Monitor, error)
	sysdigSecureClientV2() (v2.Secure, error)
}

type sysdigClients struct {
	d             *schema.ResourceData
	mu            sync.Mutex
	onceMonitor   sync.Once
	onceSecure    sync.Once
	onceCommon    sync.Once
	monitorClient monitor.SysdigMonitorClient
	secureClient  secure.SysdigSecureClient
	commonClient  common.SysdigCommonClient

	// v2
	monitorClientV2 v2.Monitor
	secureClientV2  v2.Secure
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

func (c *sysdigClients) sysdigMonitorClient() (m monitor.SysdigMonitorClient, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	monitorAPIToken := c.d.Get("sysdig_monitor_api_token").(string)
	if monitorAPIToken == "" {
		err = errors.New("sysdig monitor token not provided")
		return
	}

	c.onceMonitor.Do(func() {
		c.monitorClient = monitor.NewSysdigMonitorClient(
			monitorAPIToken,
			c.d.Get("sysdig_monitor_url").(string),
			c.d.Get("sysdig_monitor_insecure_tls").(bool),
		)

		if headers, ok := c.d.GetOk("extra_headers"); ok {
			extraHeaders := headers.(map[string]interface{})
			extraHeadersTransformed := map[string]string{}
			for key := range extraHeaders {
				extraHeadersTransformed[key] = extraHeaders[key].(string)
			}
			c.monitorClient = monitor.WithExtraHeaders(c.monitorClient, extraHeadersTransformed)
		}
	})

	return c.monitorClient, nil
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
	iamURL     string
	instanceID string
	apiKey     string
	teamID     string
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
		return nil, errors.New("missing sysdig monitor URL")
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

	if apiURL, ok = data.GetOk("ibm_monitor_url"); !ok {
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

	return &ibmVariables{
		globalVariables: &globalVariables{
			apiURL:       apiURL.(string),
			insecure:     data.Get("ibm_monitor_insecure_tls").(bool),
			extraHeaders: getExtraHeaders(data),
		},
		iamURL:     iamURL.(string),
		instanceID: instanceID.(string),
		apiKey:     apiKey.(string),
		teamID:     data.Get("ibm_monitor_team_id").(string),
	}, nil
}

func getIBMSecureVariables(data *schema.ResourceData) (*ibmVariables, error) {
	var ok bool
	var apiURL, iamURL, instanceID, apiKey interface{}

	if apiURL, ok = data.GetOk("ibm_secure_url"); !ok {
		return nil, errors.New("missing secure IBM URL")
	}

	if iamURL, ok = data.GetOk("ibm_secure_iam_url"); !ok {
		return nil, errors.New("missing secure IBM IAM URL")
	}

	if instanceID, ok = data.GetOk("ibm_secure_instance_id"); !ok {
		return nil, errors.New("missing secure IBM instance ID")
	}

	if apiKey, ok = data.GetOk("ibm_secure_api_key"); !ok {
		return nil, errors.New("missing secure IBM API key")
	}

	return &ibmVariables{
		globalVariables: &globalVariables{
			apiURL:       apiURL.(string),
			insecure:     data.Get("ibm_secure_insecure_tls").(bool),
			extraHeaders: getExtraHeaders(data),
		},
		iamURL:     iamURL.(string),
		instanceID: instanceID.(string),
		apiKey:     apiKey.(string),
		teamID:     data.Get("ibm_secure_team_id").(string),
	}, nil
}

func newSysdigMonitor(data *schema.ResourceData) (v2.Monitor, error) {
	vars, err := getSysdigMonitorVariables(data)
	if err != nil {
		return nil, err
	}
	return v2.NewSysdigMonitor(
		v2.WithToken(vars.token),
		v2.WithURL(vars.apiURL),
		v2.WithInsecure(vars.insecure),
		v2.WithExtraHeaders(vars.extraHeaders),
	), nil
}

func newSysdigSecure(data *schema.ResourceData) (v2.Secure, error) {
	vars, err := getSysdigSecureVariables(data)
	if err != nil {
		return nil, err
	}
	return v2.NewSysdigSecure(
		v2.WithToken(vars.token),
		v2.WithURL(vars.apiURL),
		v2.WithInsecure(vars.insecure),
		v2.WithExtraHeaders(vars.extraHeaders),
	), nil
}

func newIBMMonitor(data *schema.ResourceData) (v2.Monitor, error) {
	vars, err := getIBMMonitorVariables(data)
	if err != nil {
		return nil, err
	}

	return v2.NewIBMMonitor(
		v2.WithURL(vars.apiURL),
		v2.WithIBMIamURL(vars.iamURL),
		v2.WithIBMInstanceID(vars.instanceID),
		v2.WithIBMAPIKey(vars.apiKey),
		v2.WithInsecure(vars.insecure),
	), nil
}

func newIBMSecure(data *schema.ResourceData) (v2.Secure, error) {
	vars, err := getIBMSecureVariables(data)
	if err != nil {
		return nil, err
	}

	return v2.NewIBMSecure(
		v2.WithURL(vars.apiURL),
		v2.WithIBMIamURL(vars.iamURL),
		v2.WithIBMInstanceID(vars.instanceID),
		v2.WithIBMAPIKey(vars.apiKey),
		v2.WithInsecure(vars.insecure),
	), nil
}

func (c *sysdigClients) sysdigMonitorClientV2() (v2.Monitor, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var err error

	// return if already initialized
	if c.monitorClientV2 != nil {
		return c.monitorClientV2, nil
	}

	// try to initialize IBM client
	c.monitorClientV2, err = newIBMMonitor(c.d)
	if err == nil {
		return c.monitorClientV2, nil
	}

	// initialize sysdig client
	c.monitorClientV2, err = newSysdigMonitor(c.d)
	return c.monitorClientV2, err
}

func (c *sysdigClients) sysdigSecureClient() (s secure.SysdigSecureClient, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	secureAPIToken := c.d.Get("sysdig_secure_api_token").(string)
	if secureAPIToken == "" {
		err = errors.New("sysdig secure token not provided")
		return
	}

	c.onceSecure.Do(func() {
		c.secureClient = secure.NewSysdigSecureClient(
			secureAPIToken,
			c.d.Get("sysdig_secure_url").(string),
			c.d.Get("sysdig_secure_insecure_tls").(bool),
		)

		if headers, ok := c.d.GetOk("extra_headers"); ok {
			extraHeaders := headers.(map[string]interface{})
			extraHeadersTransformed := map[string]string{}
			for key := range extraHeaders {
				extraHeadersTransformed[key] = extraHeaders[key].(string)
			}
			c.secureClient = secure.WithExtraHeaders(c.secureClient, extraHeadersTransformed)
		}
	})

	return c.secureClient, nil
}

func (c *sysdigClients) sysdigSecureClientV2() (v2.Secure, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var err error

	// return if already initialized
	if c.secureClientV2 != nil {
		return c.secureClientV2, nil
	}

	// try to initialize IBM client
	c.secureClientV2, err = newIBMSecure(c.d)
	if err == nil {
		return c.secureClientV2, nil
	}

	// initialize sysdig client
	c.secureClientV2, err = newSysdigSecure(c.d)
	return c.secureClientV2, err
}

func (c *sysdigClients) sysdigCommonClient() (co common.SysdigCommonClient, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	monitorAPIToken := c.d.Get("sysdig_monitor_api_token").(string)
	secureAPIToken := c.d.Get("sysdig_secure_api_token").(string)

	if monitorAPIToken == "" && secureAPIToken == "" {
		err = errors.New("sysdig monitor and sysdig secure tokens not provided")
		return
	}

	commonAPIToken := monitorAPIToken
	commonURL := c.d.Get("sysdig_monitor_url").(string)
	commonInsecure := c.d.Get("sysdig_monitor_insecure_tls").(bool)
	if monitorAPIToken == "" {
		commonAPIToken = secureAPIToken
		commonURL = c.d.Get("sysdig_secure_url").(string)
		commonInsecure = c.d.Get("sysdig_secure_insecure_tls").(bool)
	}

	c.onceCommon.Do(func() {
		c.commonClient = common.NewSysdigCommonClient(
			commonAPIToken,
			commonURL,
			commonInsecure,
		)

		if headers, ok := c.d.GetOk("extra_headers"); ok {
			extraHeaders := headers.(map[string]interface{})
			extraHeadersTransformed := map[string]string{}
			for key := range extraHeaders {
				extraHeadersTransformed[key] = extraHeaders[key].(string)
			}
			c.commonClient = common.WithExtraHeaders(c.commonClient, extraHeadersTransformed)
		}
	})

	return c.commonClient, nil
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
