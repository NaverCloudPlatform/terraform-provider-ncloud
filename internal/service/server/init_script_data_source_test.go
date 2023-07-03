package server_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudInitScriptBasic(t *testing.T) {
	resourceName := "ncloud_init_script.foo"
	dataName := "data.ncloud_init_script.by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudInitScriptConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "os_type", resourceName, "os_type"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudInitScriptFilter(t *testing.T) {
	resourceName := "ncloud_init_script.foo"
	dataName := "data.ncloud_init_script.by_filter"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudInitScriptConfigFilter,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID("data.ncloud_init_script.by_filter"),
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "os_type", resourceName, "os_type"),
				),
			},
		},
	})
}

var testAccDataSourceNcloudInitScriptConfig = `
resource "ncloud_init_script" "foo" {
	content = "#!/usr/bin/env\nls -al"
}

data "ncloud_init_script" "by_id" {
	id = ncloud_init_script.foo.id
}
`

var testAccDataSourceNcloudInitScriptConfigFilter = `
resource "ncloud_init_script" "foo" {
	content = "#!/usr/bin/env\nls -al"
}

data "ncloud_init_script" "by_filter" {
	filter {
		name   = "init_script_no"
		values = [ncloud_init_script.foo.id]
	}
}
`
