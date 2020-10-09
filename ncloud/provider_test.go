package ncloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccClassicProviders map[string]terraform.ResourceProvider

var testAccProvider *schema.Provider
var testAccClassicProvider *schema.Provider

var credsEnvVars = []string{
	"NCLOUD_ACCESS_KEY",
	"NCLOUD_SECRET_KEY",
}

var regionEnvVar = "NCLOUD_REGION"

func init() {
	testAccProvider = getTestAccProvider(true).(*schema.Provider)
	testAccClassicProvider = getTestAccProvider(false).(*schema.Provider)

	testAccProviders = map[string]terraform.ResourceProvider{
		"ncloud": testAccProvider,
	}

	testAccClassicProviders = map[string]terraform.ResourceProvider{
		"ncloud": testAccClassicProvider,
	}
}

func getTestProvider(isVpc bool) *schema.Provider {
	if isVpc {
		return testAccProvider
	}

	return testAccClassicProvider
}

func getTestAccProviders(isVpc bool) map[string]terraform.ResourceProvider {
	if isVpc {
		return testAccProviders
	}

	return testAccClassicProviders
}

func getTestAccProvider(isVpc bool) terraform.ResourceProvider {
	testProvider := &schema.Provider{
		Schema:         schemaMap(),
		DataSourcesMap: DataSourcesMap(),
		ResourcesMap:   ResourcesMap(),
		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			d.Set("region", testAccGetRegion())
			d.Set("support_vpc", fmt.Sprintf("%v", isVpc))
			return providerConfigure(d)
		},
	}

	return testProvider
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
	if v := multiEnvSearch(credsEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(credsEnvVars, ", "))
	}

	region := testAccGetRegion()
	log.Printf("[INFO] Test: Using %s as test region", region)
}

func testAccGetRegion() string {
	v := os.Getenv(regionEnvVar)
	if v == "" {
		return "KR"
	}
	return v
}

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func getTestPrefix() string {
	rand := acctest.RandString(5)
	return fmt.Sprintf("tf%s", rand)
}

func composeConfig(config ...string) string {
	var str strings.Builder

	for _, conf := range config {
		str.WriteString(conf)
	}

	return str.String()
}
