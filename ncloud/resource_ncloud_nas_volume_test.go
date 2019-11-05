package ncloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudNasVolumeBasic(t *testing.T) {
	var volumeInstance server.NasVolumeInstance
	prefix := getTestPrefix()
	testVolumeName := prefix + "_vol"

	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			// volume_name_postfix : tf8214_vol => volume_name: n000300_tf8214_vol
			if !strings.Contains(*volumeInstance.VolumeName, testVolumeName) {
				return fmt.Errorf("not found: %s", testVolumeName)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNasVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasVolumeConfig(testVolumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasVolumeExists("ncloud_nas_volume.test", &volumeInstance),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_nas_volume.test",
						"volume_name_postfix",
						testVolumeName),
					resource.TestCheckResourceAttr(
						"ncloud_nas_volume.test",
						"volume_size",
						"500"),
				),
			},
			{
				ResourceName:            "ncloud_nas_volume.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"volume_name_postfix"},
			},
		},
	})
}

func TestAccResourceNcloudNasVolumeResize(t *testing.T) {
	var before server.NasVolumeInstance
	var after server.NasVolumeInstance
	prefix := getTestPrefix()
	testVolumeName := prefix + "_vol"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNasVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasVolumeConfig(testVolumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasVolumeExists("ncloud_nas_volume.test", &before),
				),
			},
			{
				Config: testAccNasVolumeResizeConfig(testVolumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasVolumeExists("ncloud_nas_volume.test", &after),
					testAccCheckNasVolumeNotRecreated(t, &before, &after),
				),
			},
			{
				ResourceName:            "ncloud_nas_volume.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"volume_name_postfix"},
			},
		},
	})
}

func TestAccResourceNcloudNasVolumeChangeAccessControl(t *testing.T) {
	var before server.NasVolumeInstance
	var after server.NasVolumeInstance
	prefix := getTestPrefix()
	testVolumeName := prefix + "_vol"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNasVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasVolumeConfig(testVolumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasVolumeExists("ncloud_nas_volume.test", &before),
				),
			},
			{
				Config: testAccNasVolumeChangeAccessControl(testVolumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasVolumeExists("ncloud_nas_volume.test", &after),
					testAccCheckNasVolumeNotRecreated(t, &before, &after),
				),
			},
			{
				ResourceName:            "ncloud_nas_volume.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_ip_list", "instance_custom_ip_list", "volume_name_postfix"},
			},
		},
	})
}

func testAccCheckNasVolumeExists(n string, i *server.NasVolumeInstance) resource.TestCheckFunc {
	return testAccCheckNasVolumeExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckNasVolumeExistsWithProvider(n string, i *server.NasVolumeInstance, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		client := provider.Meta().(*NcloudAPIClient)
		nasVolumeInstance, err := getNasVolumeInstance(client, rs.Primary.ID)
		if err != nil {
			return nil
		}

		if nasVolumeInstance != nil {
			*i = *nasVolumeInstance
			return nil
		}

		return fmt.Errorf("nas volume instance not found")
	}
}

func testAccCheckNasVolumeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*NcloudAPIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_nas_volume" {
			continue
		}
		volumeInstance, err := getNasVolumeInstance(client, rs.Primary.ID)
		if volumeInstance == nil {
			return nil
		}
		if err != nil {
			return err
		}
		if volumeInstance != nil && *volumeInstance.NasVolumeInstanceStatus.Code != "CREAT" {
			return fmt.Errorf("found not deleted nas volume: %s", *volumeInstance.VolumeName)
		}
	}

	return nil
}

func testAccCheckNasVolumeNotRecreated(t *testing.T,
	before, after *server.NasVolumeInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *before.NasVolumeInstanceNo != *after.NasVolumeInstanceNo {
			t.Fatalf("Ncloud NasVolumeInstanceNo have changed. Before %s. After %s", *before.NasVolumeInstanceNo, *after.NasVolumeInstanceNo)
		}
		return nil
	}
}

func testAccNasVolumeConfig(volumeNamePostfix string) string {
	return fmt.Sprintf(`
resource "ncloud_nas_volume" "test" {
	volume_name_postfix = "%s"
	volume_size = "500"
	volume_allotment_protocol_type = "NFS"
}`, volumeNamePostfix)
}

func testAccNasVolumeResizeConfig(volumeNamePostfix string) string {
	return fmt.Sprintf(`
resource "ncloud_nas_volume" "test" {
	volume_name_postfix = "%s"
	volume_size = "600"
	volume_allotment_protocol_type = "NFS"
}`, volumeNamePostfix)
}

func testAccNasVolumeChangeAccessControl(volumeNamePostfix string) string {
	return fmt.Sprintf(`
resource "ncloud_nas_volume" "test" {
	volume_name_postfix = "%s"
	volume_size = "600"
	volume_allotment_protocol_type = "NFS"
	custom_ip_list = ["10.10.10.1", "10.10.10.2"]
}`, volumeNamePostfix)
}
