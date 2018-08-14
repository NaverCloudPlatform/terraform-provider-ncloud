package main

import (
	"github.com/NaverCloudPlatform/terraform-provider-ncloud/ncloud"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: ncloud.Provider})
}
