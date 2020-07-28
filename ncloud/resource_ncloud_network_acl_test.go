package ncloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudNetworkACL_basic(t *testing.T) {
	var networkACL vpc.NetworkAcl
	resourceName := "ncloud_network_acl.nacl"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "network_acl_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", "tf-testacc-network-acl"),
					resource.TestCheckResourceAttr(resourceName, "status", "RUN"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceNcloudNetworkACLConfig() string {
	return `
resource "ncloud_vpc" "vpc" {
	name            = "tf-testacc-network-acl"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_network_acl" "nacl" {
	vpc_no      = ncloud_vpc.vpc.vpc_no
	name        = "tf-testacc-network-acl"
	description = "for test acc"
}
`
}

func testAccCheckNetworkACLExists(n string, networkACL *vpc.NetworkAcl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network acl id is set: %s", n)
		}

		client := testAccProvider.Meta().(*NcloudAPIClient)
		instance, err := getNetworkACLInstance(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		*networkACL = *instance

		return nil
	}
}
