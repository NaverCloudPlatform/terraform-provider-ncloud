package apigw_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/ncloudsdk"
)

func TestAccResourceNcloudApigw_product_basic(t *testing.T) {
	productName := fmt.Sprintf("tf-product-%s", acctest.RandString(5))

	resourceName := "ncloud_apigw_product.testing_product"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccproductConfig(productName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckproductExists(resourceName, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "product_name", productName),
					// check all the other attributes
				),
			},
		},
	})
}

func testAccCheckproductExists(n string, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		c := ncloudsdk.NewClient("https://apigateway.apigw.ntruss.com/api/v1", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))

		response, err := c.GETProductsProductid(context.Background(), &ncloudsdk.GETProductsProductidRequestQuery{
			// change value with "resource.Primary.ID"
			Productid: &resource.Primary.ID,
		})
		if response == nil {
			return err
		}
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckProductDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_apigw_product.testing_product" {
			continue
		}

		c := ncloudsdk.NewClient("https://apigateway.apigw.ntruss.com/api/v1", os.Getenv("NCLOUD_ACCESS_KEY"), os.Getenv("NCLOUD_SECRET_KEY"))
		_, err := c.GETProductsProductid(context.Background(), &ncloudsdk.GETProductsProductidRequestQuery{
			// change value with "rs.Primary.ID"
			Productid: &rs.Primary.ID,
		})
		if err != nil {
			return nil
		}
	}

	return nil
}

func testAccproductConfig(productName string) string {
	return fmt.Sprintf(`
	resource "ncloud_apigw_product" "testing_product" {
		product_name = "%[1]s"
		subscription_code = "tf-swm90"

	}`, productName)
}
