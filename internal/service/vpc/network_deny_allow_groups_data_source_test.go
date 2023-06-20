package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudNetworkACLDenyAllowGroups_basic(t *testing.T) {
	name := fmt.Sprintf("tf-ds-nacl-allow-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl_deny_allow_group.this"
	dataName := "data.ncloud_network_acl_deny_allow_groups.by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: GetTestAccProviders(true),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudNetworkACLDenyAllowGroupsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "network_acl_deny_allow_groups.0.id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataName, "network_acl_deny_allow_groups.0.network_acl_deny_allow_group_no", resourceName, "network_acl_deny_allow_group_no"),
					resource.TestCheckResourceAttrPair(dataName, "network_acl_deny_allow_groups.0.name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "network_acl_deny_allow_groups.0.description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "network_acl_deny_allow_groups.0.vpc_no", resourceName, "vpc_no"),
					resource.TestCheckResourceAttrPair(dataName, "network_acl_deny_allow_groups.0.ip_list", resourceName, "ip_list"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudNetworkACLDenyAllowGroupsConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "this" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.4.0.0/16"
}

resource "ncloud_network_acl_deny_allow_group" "this" {
  vpc_no = ncloud_vpc.this.id
  name        = "%[1]s"
  ip_list     = ["10.0.0.1", "10.0.0.2"]
}

data "ncloud_network_acl_deny_allow_groups" "by_id" {
  network_acl_deny_allow_group_no_list = [ncloud_network_acl_deny_allow_group.this.id]
}
`, name)
}
