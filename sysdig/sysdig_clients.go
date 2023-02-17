package sysdig

import (
	"errors"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/common"
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor"
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
	monitorv2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/monitor"
	securev2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/secure"
)

type SysdigClients interface {
	GetSecureEndpoint() (string, error)
	GetSecureApiToken() (string, error)

	sysdigMonitorClient() (monitor.SysdigMonitorClient, error)
	sysdigMonitorV2Client() (monitorv2.Client, error)
	sysdigSecureClient() (secure.SysdigSecureClient, error)
	sysdigSecureV2Client() (securev2.Client, error)
	sysdigCommonClient() (common.SysdigCommonClient, error)
}

type sysdigClients struct {
	d             *schema.ResourceData
	mu            sync.Mutex
	onceMonitor   sync.Once
	onceMonitorV2 sync.Once
	onceSecure    sync.Once
	onceSecureV2  sync.Once
	onceCommon    sync.Once
	monitorClient monitor.SysdigMonitorClient
	secureClient  secure.SysdigSecureClient
	commonClient  common.SysdigCommonClient

	monitorClientV2 monitorv2.Client
	secureClientV2  securev2.Client
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

func (c *sysdigClients) sysdigMonitorV2Client() (m monitorv2.Client, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	monitorAPIToken := c.d.Get("sysdig_monitor_api_token").(string)
	if monitorAPIToken == "" {
		err = errors.New("sysdig monitor token not provided")
		return
	}

	c.onceMonitorV2.Do(func() {
		c.monitorClientV2 = monitorv2.NewClient(
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
			c.monitorClientV2.SetExtraHeaders(extraHeadersTransformed)
		}
	})

	return c.monitorClientV2, nil
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

func (c *sysdigClients) sysdigSecureV2Client() (s securev2.Client, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	secureAPIToken := c.d.Get("sysdig_secure_api_token").(string)
	if secureAPIToken == "" {
		err = errors.New("sysdig secure token not provided")
		return
	}

	c.onceSecureV2.Do(func() {
		c.secureClientV2 = securev2.NewClient(
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
			c.secureClientV2.SetExtraHeaders(extraHeadersTransformed)
		}
	})

	return c.secureClientV2, nil
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
