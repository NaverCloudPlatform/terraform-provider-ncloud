package objectstorage_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
)

func TestAccDataSourceNcloudObjectStorage_object_basic(t *testing.T) {
	bucket := fmt.Sprintf("tf-bucket-%s", acctest.RandString(5))
	key := fmt.Sprintf("%s.md", acctest.RandString(5))
	dataName := "data.ncloud_objectstorage_object.by_id"
	resourceName := "ncloud_objectstorage_object.testing_object"
	content := "content for file upload testing"

	tmpFile := CreateTempFile(t, content, key)
	source := tmpFile.Name()
	defer os.Remove(source)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceObjectConfig(bucket, key, source),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^[a-z0-9-_]+\/[a-zA-Z0-9_.-]+$`)),
					resource.TestCheckResourceAttrPair(dataName, "object_id", resourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceObjectConfig(bucket, key, source string) string {
	return fmt.Sprintf(`
	resource "ncloud_objectstorage_bucket" "testing_bucket" {
		bucket_name				= "%[1]s"
	}

	resource "ncloud_objectstorage_object" "testing_object" {
		bucket					= ncloud_objectstorage_bucket.testing_bucket.bucket_name
		key						= "%[2]s"
		source					= "%[3]s"
	}

	data "ncloud_objectstorage_object" "by_id" {
		object_id				= ncloud_objectstorage_object.testing_object.id
	}
	`, bucket, key, source)
}
