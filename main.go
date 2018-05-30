package main

import (
	"github.com/hashicorp/terraform/plugin"
	"oss.navercorp.com/ncloud-paas/terraform-provider-ncloud/ncloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: ncloud.Provider})
}
