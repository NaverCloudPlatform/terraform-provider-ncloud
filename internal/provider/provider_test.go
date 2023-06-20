package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccClassicProviders map[string]*schema.Provider

var testAccProvider *schema.Provider
var testAccClassicProvider *schema.Provider

var credsEnvVars = []string{
	"NCLOUD_ACCESS_KEY",
	"NCLOUD_SECRET_KEY",
}

var regionEnvVar = "NCLOUD_REGION"

func init() {
	testAccProvider = getTestAccProvider(true)
	testAccClassicProvider = getTestAccProvider(false)

	testAccProviders = map[string]*schema.Provider{
		"ncloud": testAccProvider,
	}

	testAccClassicProviders = map[string]*schema.Provider{
		"ncloud": testAccClassicProvider,
	}
}

func getTestAccProvider(isVpc bool) *schema.Provider {
	testProvider := &schema.Provider{
		Schema:         SchemaMap(),
		DataSourcesMap: DataSourcesMap(),
		ResourcesMap:   ResourcesMap(),
		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			d.Set("region", testAccGetRegion())
			d.Set("support_vpc", isVpc)
			return ProviderConfigure(d)
		},
	}

	return testProvider
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccGetRegion() string {
	v := os.Getenv(regionEnvVar)
	if v == "" {
		return "KR"
	}
	return v
}
