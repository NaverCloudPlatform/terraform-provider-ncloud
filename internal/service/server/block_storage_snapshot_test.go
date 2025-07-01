package server_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/server"
)

func TestAccResourceNcloudBlockStorageSnapshot_vpc_basic(t *testing.T) {
	name := fmt.Sprintf("tf-snap-%s", acctest.RandString(5))
	resourceName := "ncloud_block_storage_snapshot.snapshot"
	hypervisorType := "KVM"
	serverSpec := "s2-g3"
	zone := "KR-2"
	volumeType := "CB1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVpcBlockStorageSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageSnapshotVpcConfig(name, hypervisorType, zone, serverSpec, volumeType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name+"-tf"),
					resource.TestCheckResourceAttr(resourceName, "hypervisor_type", hypervisorType),
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

func testAccCheckVpcBlockStorageSnapshotDestroy(s *terraform.State) error {
	config := TestAccProvider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_block_storage_snapshot" {
			continue
		}
		snapshot, err := server.GetVpcBlockStorageSnapshotDetail(config, rs.Primary.ID)

		if err != nil {
			return err
		}
		if snapshot != nil {
			return fmt.Errorf("unterminated snapshot : %s", *snapshot.BlockStorageSnapshotInstanceNo)
		}
	}

	return nil
}

func testAccBlockStorageSnapshotVpcConfig(name string, hypervisorType string, zone string, serverSpec string, volumeType string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_vpc" "test" {
	name               = "%[1]s"
	ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "test" {
	vpc_no             = ncloud_vpc.test.vpc_no
	name               = "%[1]s"
	subnet             = "10.5.0.0/24"
	zone               = "%[3]s"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

data "ncloud_server_image_numbers" "server_images" {
	hypervisor_type = "%[2]s"
	filter {
        name = "name"
        values = ["ubuntu-22.04-base"]
    }
}

resource "ncloud_server" "server" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_number = data.ncloud_server_image_numbers.server_images.image_number_list.0.server_image_number
	server_spec_code = "%[4]s"
	login_key_name = ncloud_login_key.loginkey.key_name
}

resource "ncloud_block_storage" "storage" {
	server_instance_no = ncloud_server.server.id
	name = "%[1]s-tf"
	size = "10"
	hypervisor_type = "%[2]s"
	volume_type = "%[5]s"
	zone = "%[3]s"
}

resource "ncloud_block_storage_snapshot" "snapshot" {
    block_storage_instance_no = ncloud_block_storage.storage.id
	name = "%[1]s-tf"
    description = "Terraform test snapshot"
}
`, name, hypervisorType, zone, serverSpec, volumeType)
}
