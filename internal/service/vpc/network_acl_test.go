package vpc_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	vpcservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/vpc"
)

func TestAccResourceNcloudNetworkACL_basic(t *testing.T) {
	var networkACL vpc.NetworkAcl
	name := fmt.Sprintf("test-network-acl-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl.nacl"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "network_acl_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
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

func TestAccResourceNcloudNetworkACL_disappears(t *testing.T) {
	var networkACL vpc.NetworkAcl
	name := fmt.Sprintf("test-network-acl-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl.nacl"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
					testAccCheckNetworkACLDisappears(&networkACL),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudNetworkACL_onlyRequiredParam(t *testing.T) {
	var networkACL vpc.NetworkAcl
	name := fmt.Sprintf("test-network-acl-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl.nacl"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLConfigOnlyRequired(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
					resource.TestMatchResourceAttr(resourceName, "vpc_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "network_acl_no", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^[a-z0-9]+$`)),
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

func TestAccResourceNcloudNetworkACL_updateName(t *testing.T) {
	var networkACL vpc.NetworkAcl
	name := fmt.Sprintf("test-network-acl-basic-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl.nacl"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLConfigOnlyRequired(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
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

func TestAccResourceNcloudNetworkACL_description(t *testing.T) {
	var networkACL vpc.NetworkAcl
	name := fmt.Sprintf("test-network-acl-desc-%s", acctest.RandString(5))
	resourceName := "ncloud_network_acl.nacl"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNcloudNetworkACLConfigDescription(name, "foo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
					resource.TestCheckResourceAttr(resourceName, "description", "foo"),
				),
			},
			{
				Config: testAccResourceNcloudNetworkACLConfigDescription(name, "bar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLExists(resourceName, &networkACL),
					resource.TestCheckResourceAttr(resourceName, "description", "bar"),
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

func testAccResourceNcloudNetworkACLConfig(name string) string {
	return testAccResourceNcloudNetworkACLConfigDescription(name, "for test acc")
}

func testAccResourceNcloudNetworkACLConfigDescription(name, description string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%[1]s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_network_acl" "nacl" {
	vpc_no      = ncloud_vpc.vpc.vpc_no
	name        = "%[1]s"
	description = "%[2]s"
}
`, name, description)
}

func testAccResourceNcloudNetworkACLConfigOnlyRequired(name string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "vpc" {
	name            = "%s"
	ipv4_cidr_block = "10.3.0.0/16"
}

resource "ncloud_network_acl" "nacl" {
	vpc_no      = ncloud_vpc.vpc.vpc_no
}
`, name)
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

		config := GetTestProvider(true).Meta().(*conn.ProviderConfig)
		instance, err := vpcservice.GetNetworkACLInstance(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*networkACL = *instance

		return nil
	}
}

func testAccCheckNetworkACLDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_network_acl" {
			continue
		}

		instance, err := vpcservice.GetNetworkACLInstance(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("Network ACL still exists")
		}
	}

	return nil
}

func testAccCheckNetworkACLDisappears(instance *vpc.NetworkAcl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

		reqParams := &vpc.DeleteNetworkAclRequest{
			RegionCode:   &config.RegionCode,
			NetworkAclNo: instance.NetworkAclNo,
		}

		_, err := config.Client.Vpc.V2Api.DeleteNetworkAcl(reqParams)

		if err := vpcservice.WaitForNcloudNetworkACLDeletion(config, *instance.NetworkAclNo); err != nil {
			return err
		}

		return err
	}
}
