package ncloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
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

func TestAccValidateOneResult(t *testing.T) {
	if err := validateOneResult(0); err == nil {
		t.Fatalf("0 result must throw 'no results' error")
	}
	if err := validateOneResult(1); err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := validateOneResult(2); err == nil {
		t.Fatalf("2 results must throw 'more than one found results'")
	}
}
