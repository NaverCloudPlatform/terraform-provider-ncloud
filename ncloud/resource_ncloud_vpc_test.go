package ncloud

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudVpc_basic(t *testing.T) {
	var vpc vpc.Vpc
	rInt := rand.Intn(16)
	cidr := fmt.Sprintf("10.%d.0.0/16", rInt)
	name := fmt.Sprintf("testacc-vpc-basic-%d", rInt)
	resourceName := "ncloud_vpc.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcConfig(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "ipv4_cidr_block", cidr),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "RUN"),
				),
			},
		},
	})
}

func testAccResourceNcloudVpcConfig(name, cidr string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test" {
	name               = "%s"
  ipv4_cidr_block    = "%s"
}
`, name, cidr)
}

func testAccCheckVpcExists(n string, vpc *vpc.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		client := testAccProvider.Meta().(*NcloudAPIClient)
		vpcInstance, err := getVpcInstance(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		*vpc = *vpcInstance

		return nil
	}
}
