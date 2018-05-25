package main

import (
	"github.com/hashicorp/terraform/plugin"
	"terraform-provider-ncloud/ncloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: ncloud.Provider})
}
