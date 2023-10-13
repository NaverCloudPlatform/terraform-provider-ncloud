package server_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudVpcAccessControlGroupBasic(t *testing.T) {
	name := fmt.Sprintf("tf-ds-acg-basic-%s", acctest.RandString(5))
	dataName := "data.ncloud_access_control_group.by_id"
	resourceName := "ncloud_access_control_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcAccessControlGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "access_control_group_no", resourceName, "access_control_group_no"),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "vpc_no", resourceName, "vpc_no"),
					resource.TestCheckResourceAttrPair(dataName, "is_default", resourceName, "is_default"),
					TestAccCheckDataSourceID("data.ncloud_access_control_group.by_filter"),
				),
			},
		},
	})
}

func TestAccDataSourceNcloudClassicAccessControlGroup_basic(t *testing.T) {
	dataName := "data.ncloud_access_control_group.by_name"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudClassicAccessControlConfig,
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "name", "ncloud-default-acg"),
					resource.TestMatchResourceAttr(dataName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(dataName, "access_control_group_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(dataName, "configuration_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(dataName, "description", "Default AccessControlGroup"),
					TestAccCheckDataSourceID("data.ncloud_access_control_group.by_filter"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudVpcAccessControlGroupConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.4.0.0/16"
}

resource "ncloud_access_control_group" "test" {
	name                  = "%[1]s"
	description           = "for acc test"
	vpc_no                = ncloud_vpc.test.id
}

data "ncloud_access_control_group" "by_id" {
	id = ncloud_access_control_group.test.id
}

data "ncloud_access_control_group" "by_filter" {
	filter {
		name   = "access_control_group_no"
		values = [ncloud_access_control_group.test.id]
	}
}
`, name)
}

var testAccDataSourceNcloudClassicAccessControlConfig = `
data "ncloud_access_control_group" "by_name" {
	name = "ncloud-default-acg"
}

data "ncloud_access_control_group" "by_filter" {
	filter {
		name   = "name"
		values = ["ncloud-default-acg"]
	}
}
`
