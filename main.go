package main

import (
	"log/slog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func main() {
	sysdigClient := sysdig.NewSysdigClients()
	defer func() {
		err := sysdigClient.Close()
		if err != nil {
			slog.Default().Error("error closing the provider", "error", err)
		}
	}()

	provider := &sysdig.SysdigProvider{SysdigClient: sysdigClient}
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: provider.Provider})
}
