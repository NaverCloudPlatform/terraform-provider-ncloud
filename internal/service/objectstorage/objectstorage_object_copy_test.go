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
	sourceName := fmt.Sprintf("%s.md", acctest.RandString(5))
	resourceName := "ncloud_objectstorage_object_copy.testing_copy"
	content := "content for file upload testing"
	key := "test/key/" + sourceName

	tmpFile := CreateTempFile(t, content, sourceName)
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
					testAccCheckObjectCopyExists(resourceName, TestAccProvider),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^[a-z0-9-_.-]+(\/[a-z0-9-_.-]+)+$`)),
					resource.TestCheckResourceAttr(resourceName, "bucket", bucketName+"-to"),
					resource.TestCheckResourceAttr(resourceName, "key", key),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source"},
			},
		},
	})
}

func TestAccResourceNcloudObjectStorage_object_copy_update_source(t *testing.T) {
	bucketName := fmt.Sprintf("tf-bucket-%s", acctest.RandString(5))
	resourceName := "ncloud_objectstorage_object_copy.testing_copy"
	sourceName := fmt.Sprintf("%s.md", acctest.RandString(5))
	content := "content for file upload testing"
	preObjectkey := "test/key/" + sourceName

	tmpFile := CreateTempFile(t, content, sourceName)
	source := tmpFile.Name()
	defer os.Remove(source)

	preObjectResourceName := "ncloud_objectstorage_object.testing_object_pre"
	postObjectResourceName := "ncloud_objectstorage_object.testing_object_post"
	postObjectKey := "test/post-key/" + sourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectCopySourcePreUpdateConfig(bucketName, preObjectkey, source, postObjectKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectCopyExists(resourceName, TestAccProvider),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^[a-z0-9-_.-]+(\/[a-z0-9-_.-]+)+$`)),
					resource.TestCheckResourceAttr(resourceName, "bucket", bucketName+"-to"),
					resource.TestCheckResourceAttr(resourceName, "key", preObjectkey),
					resource.TestCheckResourceAttrPair(resourceName, "source", preObjectResourceName, "id"),
				),
			},
			{
				Config: testAccObjectCopySourcePostUpdateConfig(bucketName, preObjectkey, source, postObjectKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectCopyExists(resourceName, TestAccProvider),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^[a-z0-9-_.-]+(\/[a-z0-9-_.-]+)+$`)),
					resource.TestCheckResourceAttr(resourceName, "bucket", bucketName+"-to"),
					resource.TestCheckResourceAttr(resourceName, "key", postObjectKey),
					resource.TestCheckResourceAttrPair(resourceName, "source", postObjectResourceName, "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source"},
			},
		},
	})
}

func TestAccResourceNcloudObjectStorage_object_copy_update_content_type(t *testing.T) {
	bucketName := fmt.Sprintf("tf-bucket-%s", acctest.RandString(5))
	resourceName := "ncloud_objectstorage_object_copy.testing_copy"
	sourceName := fmt.Sprintf("%s.md", acctest.RandString(5))
	content := "content for file upload testing"
	key := "test/key/" + sourceName

	tmpFile := CreateTempFile(t, content, sourceName)
	source := tmpFile.Name()
	defer os.Remove(source)

	newContentType := "application/json"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectCopyConfig(bucketName, key, source),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectExists(resourceName, TestAccProvider),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^[a-z0-9-_.-]+(\/[a-z0-9-_.-]+)+$`)),
					resource.TestCheckResourceAttr(resourceName, "bucket", bucketName+"-to"),
					resource.TestCheckResourceAttr(resourceName, "key", key),
				),
			},
			{
				Config: testAccObjectCopyContentTypeConfig(bucketName, key, source, newContentType),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectExists(resourceName, TestAccProvider),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^[a-z0-9-_.-]+(\/[a-z0-9-_.-]+)+$`)),
					resource.TestCheckResourceAttr(resourceName, "bucket", bucketName+"-to"),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "content_type", newContentType),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source"},
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
	config := TestAccProvider.Meta().(*conn.ProviderConfig)

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

func testAccObjectCopyContentTypeConfig(bucketName, key, source, contentType string) string {
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
		content_type		= "%[4]s"
	}	
	`, bucketName, key, source, contentType)
}

func testAccObjectCopySourcePreUpdateConfig(bucketName, key, source, postObjectKey string) string {
	return fmt.Sprintf(`
	resource "ncloud_objectstorage_bucket" "testing_bucket_from" {
		bucket_name			= "%[1]s-from"
	}

	resource "ncloud_objectstorage_bucket" "testing_bucket_to" {
		bucket_name			= "%[1]s-to"
	}

	resource "ncloud_objectstorage_object" "testing_object_pre" {
		bucket				= ncloud_objectstorage_bucket.testing_bucket_from.bucket_name
		key 				= "%[2]s"
		source				= "%[3]s"	
	}
	
	resource "ncloud_objectstorage_object" "testing_object_post" {
		bucket				= ncloud_objectstorage_bucket.testing_bucket_from.bucket_name
		key					= "%[4]s"
		source				= "%[3]s"	
	}
		
	resource "ncloud_objectstorage_object_copy" "testing_copy" {
		bucket 				= ncloud_objectstorage_bucket.testing_bucket_to.bucket_name
		key 				= "%[2]s"
		source 				= ncloud_objectstorage_object.testing_object_pre.id
	}	
	`, bucketName, key, source, postObjectKey)
}

func testAccObjectCopySourcePostUpdateConfig(bucketName, key, source, postObjectKey string) string {
	return fmt.Sprintf(`
	resource "ncloud_objectstorage_bucket" "testing_bucket_from" {
		bucket_name			= "%[1]s-from"
	}

	resource "ncloud_objectstorage_bucket" "testing_bucket_to" {
		bucket_name			= "%[1]s-to"
	}

	resource "ncloud_objectstorage_object" "testing_object_pre" {
		bucket				= ncloud_objectstorage_bucket.testing_bucket_from.bucket_name
		key 				= "%[2]s"
		source				= "%[3]s"	
	}
	
	resource "ncloud_objectstorage_object" "testing_object_post" {
		bucket				= ncloud_objectstorage_bucket.testing_bucket_from.bucket_name
		key					= "%[4]s"
		source				= "%[3]s"	
	}
		
	resource "ncloud_objectstorage_object_copy" "testing_copy" {
		bucket 				= ncloud_objectstorage_bucket.testing_bucket_to.bucket_name
		key 				= "%[4]s"
		source 				= ncloud_objectstorage_object.testing_object_post.id
	}	
	`, bucketName, key, source, postObjectKey)
}
