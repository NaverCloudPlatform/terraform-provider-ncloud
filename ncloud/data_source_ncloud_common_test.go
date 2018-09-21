package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const skipNoResultsTest = true

func testAccCheckDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("source ID not set")
		}
		return nil
	}
}
