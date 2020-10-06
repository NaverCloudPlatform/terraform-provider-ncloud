package ncloud

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceNcloudPlacementGroup_basic(t *testing.T) {
	var PlacementGroup vserver.PlacementGroup
	resourceName := "ncloud_placement_group.test"
	name := fmt.Sprintf("tf-pl-group-basic-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccResourceNcloudPlacementGroupConfig(name),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementGroupExists(resourceName, &PlacementGroup),
					resource.TestMatchResourceAttr(resourceName, "placement_group_no", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "placement_group_type", "AA"),
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

func TestAccResourceNcloudPlacementGroup_disappears(t *testing.T) {
	var PlacementGroup vserver.PlacementGroup
	resourceName := "ncloud_placement_group.test"
	name := fmt.Sprintf("tf-pl-group-disappear-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccResourceNcloudPlacementGroupConfig(name),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementGroupExists(resourceName, &PlacementGroup),
					testAccCheckPlacementGroupDisappears(&PlacementGroup),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNcloudPlacementGroup_updateName(t *testing.T) {
	var PlacementGroup vserver.PlacementGroup
	resourceName := "ncloud_placement_group.test"
	name := fmt.Sprintf("tf-pl-group-update-%s", acctest.RandString(5))
	updateName := fmt.Sprintf("tf-pl-group-update-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:   testAccResourceNcloudPlacementGroupConfig(name),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementGroupExists(resourceName, &PlacementGroup),
				),
			},
			{
				Config:   testAccResourceNcloudPlacementGroupConfig(updateName),
				SkipFunc: testOnlyVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementGroupExists(resourceName, &PlacementGroup),
				),
				ExpectError: regexp.MustCompile("Change 'name' is not support, Please set `name` as a old value"),
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

func testAccResourceNcloudPlacementGroupConfig(name string) string {
	return fmt.Sprintf(`
resource "ncloud_placement_group" "test" {
	name                  = "%[1]s"
}
`, name)
}

func testAccCheckPlacementGroupExists(n string, PlacementGroup *vserver.PlacementGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Placement group id is set: %s", n)
		}

		config := testAccProvider.Meta().(*ProviderConfig)
		instance, err := getPlacementGroupInstance(config, rs.Primary.ID)
		if err != nil {
			return err
		}

		*PlacementGroup = *instance

		return nil
	}
}

func testAccCheckPlacementGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_placement_group" {
			continue
		}

		instance, err := getPlacementGroupInstance(config, rs.Primary.ID)

		if err != nil {
			return err
		}

		if instance != nil {
			return errors.New("Placement group still exists")
		}
	}

	return nil
}

func testAccCheckPlacementGroupDisappears(instance *vserver.PlacementGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*ProviderConfig)

		reqParams := &vserver.DeletePlacementGroupRequest{
			RegionCode:       &config.RegionCode,
			PlacementGroupNo: instance.PlacementGroupNo,
		}

		_, err := config.Client.vserver.V2Api.DeletePlacementGroup(reqParams)

		return err
	}
}
