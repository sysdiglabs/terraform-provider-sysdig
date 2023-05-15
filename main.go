package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: sysdig.Provider})
}
