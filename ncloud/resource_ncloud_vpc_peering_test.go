package ncloud

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudVpcPeering_basic(t *testing.T) {
	var vpcPeeringInstance vpc.VpcPeeringInstance
	resourceName := "ncloud_vpc_peering.foo"
	name := fmt.Sprintf("test-peering-basic-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccResourceNcloudVpcPeeringConfig(name),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringExists(resourceName, &vpcPeeringInstance),
					resource.TestMatchResourceAttr(resourceName, "source_vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "target_vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "has_reverse_vpc_peering", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_between_accounts", "false"),
				),
			},
			{
				SkipFunc:          testOnlyVpc,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNcloudVpcPeering_Peering(t *testing.T) {
	var vpcPeeringInstance vpc.VpcPeeringInstance
	resourceNameMain := "ncloud_vpc_peering.foo"
	resourceNamePeer := "ncloud_vpc_peering.bar"
	name := fmt.Sprintf("test-peering-basic-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccResourceNcloudVpcPeeringConfigAdd(name),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringExists(resourceNamePeer, &vpcPeeringInstance),
					resource.TestMatchResourceAttr(resourceNamePeer, "source_vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceNamePeer, "target_vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceNamePeer, "has_reverse_vpc_peering", "true"),
					resource.TestCheckResourceAttr(resourceNamePeer, "is_between_accounts", "false"),
					resource.TestCheckResourceAttr(resourceNameMain, "has_reverse_vpc_peering", "true"),
				),
			},
		},
	})
}

func TestAccResourceNcloudVpcPeering_disappears(t *testing.T) {
	var vpcPeeringInstance vpc.VpcPeeringInstance
	resourceName := "ncloud_vpc_peering.foo"
	name := fmt.Sprintf("test-peering-disap-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccResourceNcloudVpcPeeringConfig(name),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringExists(resourceName, &vpcPeeringInstance),
					testAccCheckVpcPeeringDisappears(&vpcPeeringInstance),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudVpcPeering_description(t *testing.T) {
	var vpcPeeringInstance vpc.VpcPeeringInstance
	resourceName := "ncloud_vpc_peering.foo"
	name := fmt.Sprintf("test-peering-desc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccResourceNcloudVpcPeeringConfigDescription(name, "foo"),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringExists(resourceName, &vpcPeeringInstance),
					resource.TestCheckResourceAttr(resourceName, "description", "foo"),
				),
			},
			{
				Config:   testAccResourceNcloudVpcPeeringConfigDescription(name, "bar"),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringExists(resourceName, &vpcPeeringInstance),
					resource.TestCheckResourceAttr(resourceName, "description", "bar"),
				),
			},
			{
				SkipFunc:          testOnlyVpc,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceNcloudVpcPeeringConfig(name string) string {
	return testAccResourceNcloudVpcPeeringConfigDescription(name, "for acc test")
}

func testAccResourceNcloudVpcPeeringConfigDescription(name, description string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "main" {
	name               = "%[1]s-a"
	ipv4_cidr_block    = "10.4.0.0/16"
}

resource "ncloud_vpc" "peer" {
	name               = "%[1]s-b"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_vpc_peering" "foo" {
	name          = "%[1]s-foo"
	source_vpc_no = ncloud_vpc.main.id
	target_vpc_no = ncloud_vpc.peer.id
	description   = "%[2]s"
}
`, name, description)
}

func testAccResourceNcloudVpcPeeringConfigAdd(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "main" {
	name               = "%[1]s-a"
	ipv4_cidr_block    = "10.4.0.0/16"
}

resource "ncloud_vpc" "peer" {
	name               = "%[1]s-b"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_vpc_peering" "foo" {
	name          = "%[1]s-foo"
	source_vpc_no = ncloud_vpc.main.id
	target_vpc_no = ncloud_vpc.peer.id
}

resource "ncloud_vpc_peering" "bar" {
	name          = "%[1]s-bar"
	source_vpc_no = ncloud_vpc.peer.id
	target_vpc_no = ncloud_vpc.main.id
}
`, name)
}

func testAccCheckVpcPeeringExists(n string, vpcPeering *vpc.VpcPeeringInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC peering id is set: %s", n)
		}

		config := testAccProvider.Meta().(*ProviderConfig)
		instance, err := getVpcPeeringInstance(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*vpcPeering = *instance

		return nil
	}
}

func testAccCheckVpcPeeringDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_vpc_peering" {
			continue
		}

		instance, err := getVpcPeeringInstance(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("VPC Peering still exists")
		}
	}

	return nil
}

func testAccCheckVpcPeeringDisappears(instance *vpc.VpcPeeringInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)

		reqParams := &vpc.DeleteVpcPeeringInstanceRequest{
			RegionCode:           &config.RegionCode,
			VpcPeeringInstanceNo: instance.VpcPeeringInstanceNo,
		}

		_, err := config.Client.vpc.V2Api.DeleteVpcPeeringInstance(reqParams)

		if err := waitForNcloudVpcPeeringDeletion(config, *instance.VpcPeeringInstanceNo); err != nil {
			return err
		}

		return err
	}
}
