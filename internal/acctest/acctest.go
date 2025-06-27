package acctest

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/terraform-providers/terraform-provider-ncloud/internal/provider"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/provider/fwprovider"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/verify"
)

const (
	SkipNoResultsTest = true
	// Provider name for single configuration testing
	ProviderName = "ncloud"
)

// ProtoV6ProviderFactories is a static map containing only the main provider instance
// for testing
var (
	ProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error) = protoV6ProviderFactoriesInit(context.Background(), ProviderName)
)

var TestAccProvider *schema.Provider

// testAccProviderConfigure ensures Provider is only configured once
var testAccProviderConfigure sync.Once

var credsEnvVars = []string{
	"NCLOUD_ACCESS_KEY",
	"NCLOUD_SECRET_KEY",
}

var regionEnvVar = "NCLOUD_REGION"

func init() {
	TestAccProvider = getTestAccProvider()
}

func getTestAccProvider() *schema.Provider {
	p := provider.New(context.Background())
	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		d.Set("region", testAccGetRegion())
		d.Set("support_vpc", true)
		return provider.ProviderConfigure(ctx, d)
	}
	return p
}

func TestAccPreCheck(t *testing.T) {
	testAccProviderConfigure.Do(func() {
		if v := multiEnvSearch(credsEnvVars); v == "" {
			t.Fatalf("One of %s must be set for acceptance tests", strings.Join(credsEnvVars, ", "))
		}

		region := testAccGetRegion()
		log.Printf("[INFO] Test: Using %s as test region", region)

		diags := TestAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
		if diags.HasError() {
			t.Fatalf("configuring provider: %v", diags)
		}
	})
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

func GetTestPrefix() string {
	rand := acctest.RandString(5)
	return fmt.Sprintf("tf%s", rand)
}

func ComposeConfig(config ...string) string {
	var str strings.Builder

	for _, conf := range config {
		str.WriteString(conf)
	}

	return str.String()
}

func TestAccCheckResourceDisappears(provider *schema.Provider, resource *schema.Resource, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceState.Primary.ID == "" {
			return fmt.Errorf("resource ID missing: %s", resourceName)
		}

		return resource.Delete(resource.Data(resourceState.Primary), provider.Meta())
	}
}

func TestAccCheckDataSourceID(n string) resource.TestCheckFunc {
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
	if err := ValidateOneResult(0); err == nil {
		t.Fatalf("0 result must throw 'no results' error")
	}
	if err := ValidateOneResult(1); err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := ValidateOneResult(2); err == nil {
		t.Fatalf("2 results must throw 'more than one found results'")
	}
}

func GetTestClusterName() string {
	rInt := acctest.RandIntRange(1, 9999)
	testClusterName := fmt.Sprintf("tf-%d-cluster", rInt)
	return testClusterName
}

func protoV6ProviderFactoriesInit(ctx context.Context, providerNames ...string) map[string]func() (tfprotov6.ProviderServer, error) {
	factories := make(map[string]func() (tfprotov6.ProviderServer, error), len(providerNames))

	for _, name := range providerNames {
		factories[name] = func() (tfprotov6.ProviderServer, error) {
			providerServerFactory, _, err := protoV6TestProviderServerFactory(ctx)

			if err != nil {
				return nil, err
			}

			return providerServerFactory(), nil
		}
	}

	return factories
}

func protoV6TestProviderServerFactory(ctx context.Context) (func() tfprotov6.ProviderServer, *schema.Provider, error) {
	primary := provider.New(ctx)
	primary.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		d.Set("region", testAccGetRegion())
		d.Set("support_vpc", true)
		return provider.ProviderConfigure(ctx, d)
	}

	upgradedSdkProvider, err := tf5to6server.UpgradeServer(
		ctx,
		primary.GRPCProvider,
	)
	if err != nil {
		return nil, nil, err
	}

	servers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedSdkProvider
		},
		providerserver.NewProtocol6(fwprovider.New(primary)),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		return nil, nil, err
	}

	return muxServer.ProviderServer, primary, nil
}
