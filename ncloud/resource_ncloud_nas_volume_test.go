package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"strings"
	"testing"
)

func TestAccNcloudNasVolume_basic(t *testing.T) {
	var volumeInstance sdk.NasVolumeInstance
	prefix := getTestPrefix()
	testVolumeName := prefix + "_vol"
	log.Printf("[DEBUG] testVolumeName: %s", testVolumeName)

	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			// volume_name_postfix : tf8214_vol => volume_name: n000300_tf8214_vol
			if !strings.Contains(volumeInstance.VolumeName, testVolumeName) {
				return fmt.Errorf("not found: %s", testVolumeName)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "ncloud_nas_volume.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNasVolumeDestroy,
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
						"volume_size_gb",
						"500"),
				),
			},
		},
	})
}

func testAccCheckNasVolumeExists(n string, i *sdk.NasVolumeInstance) resource.TestCheckFunc {
	return testAccCheckNasVolumeExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckNasVolumeExistsWithProvider(n string, i *sdk.NasVolumeInstance, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		provider := providerF()
		conn := provider.Meta().(*NcloudSdk).conn
		nasVolumeInstance, err := getNasVolumeInstance(conn, rs.Primary.ID)
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
	return testAccCheckNasVolumeDestroyWithProvider(s, testAccProvider)
}

func testAccCheckNasVolumeDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*NcloudSdk).conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_nas_volume" {
			continue
		}
		volumeInstance, err := getNasVolumeInstance(conn, rs.Primary.ID)
		log.Printf("[DEBUG] testAccCheckNasVolumeDestroyWithProvider volumeInstance: %#v", volumeInstance)
		if volumeInstance == nil {
			return nil
		}
		if err != nil {
			return err
		}
		if volumeInstance != nil && volumeInstance.NasVolumeInstanceStatus.Code != "CREAT" {
			return fmt.Errorf("found not deleted nas volume: %s", volumeInstance.VolumeName)
		}
	}

	return nil
}

func testAccNasVolumeConfig(volumeNamePostfix string) string {
	return fmt.Sprintf(`
resource "ncloud_nas_volume" "test" {
	"volume_name_postfix" = "%s"
	"volume_size_gb" = "500"
	"volume_allotment_protocol_type_code" = "NFS"
}`, volumeNamePostfix)
}
