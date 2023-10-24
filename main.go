package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func main() {
	sysdigClient := sysdig.NewSysdigClients()
	defer sysdigClient.Close()

	provider := &sysdig.SysdigProvider{SysdigClient: sysdigClient}
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: provider.Provider})
}
