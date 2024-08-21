package objectstorage_test

import (
	"fmt"
	"testing"

	randacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudObjectStorage_bucket_basic(t *testing.T) {
	dataName := "data.ncloud_objectstorage_bucket.by_name"
	resourceName := "ncloud_objectstorage_bucket.testing_bucket"
	testBucketName := fmt.Sprintf("tf-bucket-%s", randacctest.RandString(4))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBucketConfig(testBucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataName, "bucket_name", resourceName, "bucket_name"),
				),
			},
		},
	})

}

func testAccDataSourceBucketConfig(testBucketName string) string {
	return fmt.Sprintf(`
	resource "ncloud_objectstorage_bucket" "testing_bucket" {
		bucket_name				= "%[1]s"
	}

	data "ncloud_objectstorage_bucket" "by_name" {
		bucket_name				= ncloud_objectstorage_bucket.testing_bucket.bucket_name
	}
	`, testBucketName)
}
