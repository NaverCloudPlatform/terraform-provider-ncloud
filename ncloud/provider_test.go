package ncloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"ncloud": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("NCLOUD_PROFILE"); v == "" {
		if v := os.Getenv("NCLOUD_ACCESS_KEY"); v == "" {
			t.Fatal("NCLOUD_ACCESS_KEY must be set for acceptance tests")
		}
		if v := os.Getenv("NCLOUD_SECRET_KEY"); v == "" {
			t.Fatal("NCLOUD_SECRET_KEY must be set for acceptance tests")
		}
	}

	region := testAccGetRegion()
	log.Printf("[INFO] Test: Using %s as test region", region)
	os.Setenv("NCLOUD_DEFAULT_REGION", region)

	err := testAccProvider.Configure(terraform.NewResourceConfig(nil))
	if err != nil {
		t.Fatal(err)
	}
}

func testAccGetRegion() string {
	v := os.Getenv("NCLOUD_DEFAULT_REGION")
	if v == "" {
		return "KR"
	}
	return v
}
