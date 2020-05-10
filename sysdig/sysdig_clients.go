package sysdig

import (
	"errors"
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"sync"
)

type SysdigClients interface {
	sysdigMonitorClient() (monitor.SysdigMonitorClient, error)
	sysdigSecureClient() (secure.SysdigSecureClient, error)
}

type sysdigClients struct {
	d                    *schema.ResourceData
	onceMonitor          sync.Once
	onceSecure           sync.Once
	_sysdigMonitorClient monitor.SysdigMonitorClient
	_sysdigSecureClient  secure.SysdigSecureClient
}

func (c *sysdigClients) sysdigMonitorClient() (m monitor.SysdigMonitorClient, err error) {
	monitorAPIToken := c.d.Get("sysdig_monitor_api_token").(string)
	if monitorAPIToken == "" {
		err = errors.New("sysdig monitor token not provided")
		return
	}

	c.onceMonitor.Do(func() {
		c._sysdigMonitorClient = monitor.NewSysdigMonitorClient(
			monitorAPIToken,
			c.d.Get("sysdig_monitor_url").(string),
		)
	})

	return c._sysdigMonitorClient, nil
}

func (c *sysdigClients) sysdigSecureClient() (s secure.SysdigSecureClient, err error) {
	secureAPIToken := c.d.Get("sysdig_secure_api_token").(string)
	if secureAPIToken == "" {
		err = errors.New("sysdig secure token not provided")
		return
	}

	c.onceSecure.Do(func() {
		c._sysdigSecureClient = secure.NewSysdigSecureClient(
			secureAPIToken,
			c.d.Get("sysdig_secure_url").(string),
		)
	})

	return c._sysdigSecureClient, nil
}
