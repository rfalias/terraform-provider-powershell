package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/rfalias/terraform-provider-powershell/pypwsh"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pypwsh.Provider,
	})
}
