package objectstorage_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func TestAccResourceNcloudObjectStorage_object_copy_basic(t *testing.T) {
	bucketName := fmt.Sprintf("tf-bucket-%s", acctest.RandString(5))
	key := fmt.Sprintf("%s.md", acctest.RandString(5))
	resourceName := "ncloud_objectstorage_object_copy.testing_copy"
	content := "content for file upload testing"

	tmpFile := CreateTempFile(t, content, key)
	source := tmpFile.Name()
	defer os.Remove(source)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckObjectCopyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectCopyConfig(bucketName, key, source),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectCopyExists(resourceName, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^[a-z0-9-_.-]+(\/[a-z0-9-_.-]+)+$`)),
					resource.TestCheckResourceAttr(resourceName, "bucket", bucketName+"-to"),
					resource.TestCheckResourceAttr(resourceName, "key", key),
				),
			},
		},
	})
}

func testAccCheckObjectCopyExists(n string, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		resp, err := config.Client.ObjectStorage.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: ncloud.String(resource.Primary.Attributes["bucket"]),
			Key:    ncloud.String(resource.Primary.Attributes["key"]),
		})
		if err != nil {
			return err
		}

		if resp != nil {
			return nil
		}

		return fmt.Errorf("Object not found")
	}
}

func testAccCheckObjectCopyDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_objectstorage_object_copy" {
			continue
		}

		resp, err := config.Client.ObjectStorage.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: ncloud.String(rs.Primary.Attributes["bucket"]),
			Key:    ncloud.String(rs.Primary.Attributes["key"]),
		})
		if resp != nil {
			return fmt.Errorf("Object found")
		}

		if err != nil {
			return nil
		}
	}

	return nil
}

func testAccObjectCopyConfig(bucketName, key, source string) string {
	return fmt.Sprintf(`
	resource "ncloud_objectstorage_bucket" "testing_bucket_from" {
		bucket_name			= "%[1]s-from"
	}

	resource "ncloud_objectstorage_bucket" "testing_bucket_to" {
		bucket_name			= "%[1]s-to"
	}

	resource "ncloud_objectstorage_object" "testing_object" {
		bucket				= ncloud_objectstorage_bucket.testing_bucket_from.bucket_name
		key 				= "%[2]s"
		source				= "%[3]s"	
	}
		
	resource "ncloud_objectstorage_object_copy" "testing_copy" {
		bucket 				= ncloud_objectstorage_bucket.testing_bucket_to.bucket_name
		key 				= "%[2]s"
		source 				= ncloud_objectstorage_object.testing_object.id
	}	
	`, bucketName, key, source)
}
