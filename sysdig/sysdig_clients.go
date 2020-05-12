package sysdig

import (
	"errors"
	"github.com/draios/terraform-provider-sysdig/sysdig/common"
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"sync"
)

type SysdigClients interface {
	sysdigMonitorClient() (monitor.SysdigMonitorClient, error)
	sysdigSecureClient() (secure.SysdigSecureClient, error)
	sysdigCommonClient() (common.SysdigCommonClient, error)
}

type sysdigClients struct {
	d             *schema.ResourceData
	onceMonitor   sync.Once
	onceSecure    sync.Once
	onceCommon    sync.Once
	monitorClient monitor.SysdigMonitorClient
	secureClient  secure.SysdigSecureClient
	commonClient  common.SysdigCommonClient
}

func (c *sysdigClients) sysdigMonitorClient() (m monitor.SysdigMonitorClient, err error) {
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
	})

	return c.monitorClient, nil
}

func (c *sysdigClients) sysdigSecureClient() (s secure.SysdigSecureClient, err error) {
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
	})

	return c.secureClient, nil
}

func (c *sysdigClients) sysdigCommonClient() (co common.SysdigCommonClient, err error) {
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
	})

	return c.commonClient, nil

}
