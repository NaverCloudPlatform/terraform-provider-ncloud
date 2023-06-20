package vpc

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
)

func TestAccResourceNcloudVpc_basic(t *testing.T) {
	var vpc vpc.Vpc
	rInt := rand.Intn(16)
	cidr := fmt.Sprintf("10.%d.0.0/16", rInt)
	name := fmt.Sprintf("test-vpc-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_vpc.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    GetTestAccProviders(true),
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudVpcConfig(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "ipv4_cidr_block", cidr),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "default_network_acl_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "default_access_control_group_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "default_public_route_table_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "default_private_route_table_no", regexp.MustCompile(`^\d+$`)),
				),
			},
		},
	})
}

func TestAccResourceNcloudVpc_disappears(t *testing.T) {
	var vpc vpc.Vpc
	rInt := rand.Intn(16)
	cidr := fmt.Sprintf("10.%d.0.0/16", rInt)
	name := fmt.Sprintf("test-vpc-disapr-%s", acctest.RandString(5))
	resourceName := "ncloud_vpc.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    GetTestAccProviders(true),
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudVpcConfig(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(resourceName, &vpc),
					testAccCheckVpcDisappears(&vpc),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudVpc_updateName(t *testing.T) {
	var vpc vpc.Vpc
	rInt := rand.Intn(16)
	cidr := fmt.Sprintf("10.%d.0.0/16", rInt)
	name := fmt.Sprintf("test-vpc-name-%s", acctest.RandString(5))
	resourceName := "ncloud_vpc.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    GetTestAccProviders(true),
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudVpcConfig(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(resourceName, &vpc),
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

		config := GetTestProvider(true).Meta().(*provider.ProviderConfig)
		vpcInstance, err := GetVpcInstance(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*vpc = *vpcInstance

		return nil
	}
}

func testAccCheckVpcDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*provider.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_vpc" {
			continue
		}

		instance, err := GetVpcInstance(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("VPC still exists")
		}
	}

	return nil
}

func testAccCheckVpcDisappears(instance *vpc.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := GetTestProvider(true).Meta().(*provider.ProviderConfig)

		reqParams := &vpc.DeleteVpcRequest{
			RegionCode: &config.RegionCode,
			VpcNo:      instance.VpcNo,
		}

		_, err := config.Client.Vpc.V2Api.DeleteVpc(reqParams)

		if err := waitForNcloudVpcDeletion(config, *instance.VpcNo); err != nil {
			return err
		}

		return err
	}
}
