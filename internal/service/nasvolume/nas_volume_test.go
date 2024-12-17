package nasvolume_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/nasvolume"
)

func TestAccResourceNcloudNasVolume_classic_basic(t *testing.T) {
	testAccResourceNcloudNasVolumeBasic(t, false)
}

func TestAccResourceNcloudNasVolume_vpc_basic(t *testing.T) {
	testAccResourceNcloudNasVolumeBasic(t, true)
}

func testAccResourceNcloudNasVolumeBasic(t *testing.T, isVpc bool) {
	var volumeInstance nasvolume.NasVolume
	postfix := GetTestPrefix()
	resourceName := "ncloud_nas_volume.test"
	provider := GetTestProvider(isVpc)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: GetTestProviderFactories(isVpc),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNasVolumeDestroy(state, provider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNasVolumeConfig(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasVolumeExists(resourceName, &volumeInstance, provider),
					resource.TestCheckResourceAttr(resourceName, "volume_name_postfix", postfix),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(fmt.Sprintf(`^n\d+_%s$`, postfix))),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "500"),
					resource.TestCheckResourceAttr(resourceName, "volume_total_size", "500"),
					resource.TestCheckResourceAttr(resourceName, "snapshot_volume_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "volume_allotment_protocol_type", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "is_event_configuration", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_snapshot_configuration", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"volume_name_postfix"},
			},
		},
	})
}

func TestAccResourceNcloudNasVolume_classic_resize(t *testing.T) {
	testAccResourceNcloudNasVolumeResize(t, false)
}

func TestAccResourceNcloudNasVolume_vpc_resize(t *testing.T) {
	testAccResourceNcloudNasVolumeResize(t, true)
}

func testAccResourceNcloudNasVolumeResize(t *testing.T, isVpc bool) {
	var before nasvolume.NasVolume
	var after nasvolume.NasVolume
	postfix := GetTestPrefix()
	resourceName := "ncloud_nas_volume.test"
	provider := GetTestProvider(isVpc)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: GetTestProviderFactories(isVpc),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNasVolumeDestroy(state, provider)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNasVolumeConfig(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasVolumeExists(resourceName, &before, provider),
				),
			},
			{
				Config: testAccNasVolumeResizeConfig(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasVolumeExists(resourceName, &after, provider),
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

func TestAccResourceNcloudNasVolume_classic_changeAccessControl(t *testing.T) {
	var before nasvolume.NasVolume
	var after nasvolume.NasVolume
	postfix := GetTestPrefix()
	resourceName := "ncloud_nas_volume.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ClassicProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNasVolumeDestroy(state, GetTestProvider(false))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNasVolumeConfig(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasVolumeExists(resourceName, &before, GetTestProvider(false)),
				),
			},
			{
				Config: testAccNasVolumeChangeAccessControlClassic(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasVolumeExists(resourceName, &after, GetTestProvider(false)),
					testAccCheckNasVolumeNotRecreated(t, &before, &after),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"volume_name_postfix"},
			},
		},
	})
}

func TestAccResourceNcloudNasVolume_vpc_changeAccessControl(t *testing.T) {
	t.Skip()
	{
		// Skip: deprecated server_image_product_code
	}

	var before nasvolume.NasVolume
	var after nasvolume.NasVolume
	postfix := GetTestPrefix()
	resourceName := "ncloud_nas_volume.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckNasVolumeDestroy(state, GetTestProvider(true))
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNasVolumeConfig(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasVolumeExists(resourceName, &before, GetTestProvider(true)),
				),
			},
			{
				Config: testAccNasVolumeChangeAccessControlVpc(postfix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNasVolumeExists(resourceName, &after, GetTestProvider(true)),
					testAccCheckNasVolumeNotRecreated(t, &before, &after),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"volume_name_postfix"},
			},
		},
	})
}

func testAccCheckNasVolumeExists(n string, i *nasvolume.NasVolume, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		nasVolumeInstance, err := nasvolume.GetNasVolume(config, rs.Primary.ID)
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

func testAccCheckNasVolumeDestroy(s *terraform.State, provider *schema.Provider) error {
	config := provider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_nas_volume" {
			continue
		}
		volumeInstance, err := nasvolume.GetNasVolume(config, rs.Primary.ID)
		if volumeInstance == nil {
			return nil
		}
		if err != nil {
			return err
		}
		if volumeInstance != nil && *volumeInstance.Status != "CREAT" {
			return fmt.Errorf("found not deleted nas volume: %s", *volumeInstance.VolumeName)
		}
	}

	return nil
}

func testAccCheckNasVolumeNotRecreated(t *testing.T,
	before, after *nasvolume.NasVolume) resource.TestCheckFunc {
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

func testAccNasVolumeChangeAccessControlClassic(volumeNamePostfix string) string {
	return fmt.Sprintf(`
resource "ncloud_login_key" "loginkey" {
	key_name = "%[1]s-key"
}

resource "ncloud_server" "server-foo" {
	name = "%[1]s-foo"
	server_image_product_code = "SPSW0LINUX000046"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}
resource "ncloud_server" "server-bar" {
	name = "%[1]s-bar"
	server_image_product_code = "SPSW0LINUX000046"
	server_product_code = "SPSVRSTAND000004"
	login_key_name = "${ncloud_login_key.loginkey.key_name}"
}

resource "ncloud_nas_volume" "test" {
	volume_name_postfix = "%[1]s"
	volume_size = "600"
	volume_allotment_protocol_type = "NFS"
	custom_ip_list = [ncloud_server.server-bar.private_ip]
	server_instance_no_list = [ncloud_server.server-foo.id]
}`, volumeNamePostfix)
}

func testAccNasVolumeChangeAccessControlVpc(volumeNamePostfix string) string {
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
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_server" "server-foo" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	server_product_code = "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"
	login_key_name = ncloud_login_key.loginkey.key_name
}

resource "ncloud_server" "server-bar" {
	subnet_no = ncloud_subnet.test.id
	name = "%[1]s"
	server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
	server_product_code = "SVR.VSVR.STAND.C002.M008.NET.HDD.B050.G002"
	login_key_name = ncloud_login_key.loginkey.key_name
}

resource "ncloud_nas_volume" "test" {
	volume_name_postfix = "%[1]s"
	volume_size = "600"
	volume_allotment_protocol_type = "NFS"
	server_instance_no_list = [ncloud_server.server-foo.id,ncloud_server.server-bar.id]
}`, volumeNamePostfix)
}
