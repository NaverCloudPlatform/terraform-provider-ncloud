package ncloud

import (
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccNcloudInstance_basic(t *testing.T) {
	var instance sdk.ServerInstance
	testServerName := getTestServerName()

	testCheck := func() func(*terraform.State) error {
		return func(*terraform.State) error {
			if instance.ServerName != testServerName {
				return fmt.Errorf("not found: %s", testServerName)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "ncloud_instance.instance",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig(testServerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"ncloud_instance.instance", &instance),
					testCheck(),
					resource.TestCheckResourceAttr(
						"ncloud_instance.instance",
						"server_image_product_code",
						"SPSW0LINUX000032"),
					resource.TestCheckResourceAttr(
						"ncloud_instance.instance",
						"server_product_code",
						"SPSVRSTAND000004"),
				),
			},
		},
	})
}

func TestAccNcloudInstance_changeServerInstanceSpec(t *testing.T) {
	var before sdk.ServerInstance
	var after sdk.ServerInstance
	testServerName := getTestServerName()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "ncloud_instance.instance",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig(testServerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"ncloud_instance.instance", &before),
				),
			},
			{
				Config: testAccInstanceChangeSpecConfig(testServerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"ncloud_instance.instance", &after),
					testAccCheckInstanceNotRecreated(
						t, &before, &after),
				),
			},
		},
	})
}

func getTestServerName() string {
	rInt := acctest.RandIntRange(1, 9999)
	testServerName := fmt.Sprintf("terraform-test-%d", rInt)
	return testServerName
}

func testAccCheckInstanceExists(n string, i *sdk.ServerInstance) resource.TestCheckFunc {
	return testAccCheckInstanceExistsWithProvider(n, i, func() *schema.Provider { return testAccProvider })
}

func testAccCheckInstanceExistsWithProvider(n string, i *sdk.ServerInstance, providerF func() *schema.Provider) resource.TestCheckFunc {
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
		instance, err := getServerInstance(conn, rs.Primary.ID)
		if err != nil {
			return nil
		}

		if instance != nil {
			*i = *instance
			return nil
		}

		return fmt.Errorf("instance not found")
	}
}

func testAccCheckInstanceNotRecreated(t *testing.T,
	before, after *sdk.ServerInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ServerInstanceNo != after.ServerInstanceNo {
			t.Fatalf("Ncloud Instance IDs have changed. Before %s. After %s", before.ServerInstanceNo, after.ServerInstanceNo)
		}
		return nil
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	return testAccCheckInstanceDestroyWithProvider(s, testAccProvider)
}

func testAccCheckInstanceDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*NcloudSdk).conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_instance" {
			continue
		}

		instance, err := getServerInstance(conn, rs.Primary.ID)
		if err == nil {
			if instance != nil && instance.ServerInstanceStatusName != "stopped" {
				return fmt.Errorf("found unstopped instance: %s", instance.ServerInstanceNo)
			}
		}

		if instance == nil {
			continue
		}

		return err
	}

	return nil
}

func testAccInstanceConfig(testServerName string) string {
	return fmt.Sprintf(`
resource "ncloud_instance" "instance" {
"server_name" = "%s"
"server_image_product_code" = "SPSW0LINUX000032"
"server_product_code" = "SPSVRSTAND000004"
}
`, testServerName)
}

func testAccInstanceChangeSpecConfig(testServerName string) string {
	return fmt.Sprintf(`
resource "ncloud_instance" "instance" {
"server_name" = "%s"
"server_image_product_code" = "SPSW0LINUX000032"
"server_product_code" = "SPSVRSTAND000056"
}
`, testServerName)
}
