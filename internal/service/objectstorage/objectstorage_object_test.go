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

func TestAccResourceNcloudObjectStorage_object_basic(t *testing.T) {
	bucketName := fmt.Sprintf("tf-bucket-%s", acctest.RandString(5))
	key := fmt.Sprintf("%s.md", acctest.RandString(5))
	resourceName := "ncloud_objectstorage_object.testing_object"
	content := "content for file upload testing"

	tmpFile := CreateTempFile(t, content, key)
	source := tmpFile.Name()
	defer os.Remove(source)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig(bucketName, key, source),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectExists(resourceName, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^objectstorage_object::[a-z]{2}::[a-z0-9-_]+::[a-zA-Z0-9_.-]+$`)),
					resource.TestCheckResourceAttr(resourceName, "bucket", bucketName),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "source", source),
				),
			},
		},
	})
}

func testAccCheckObjectExists(n string, provider *schema.Provider) resource.TestCheckFunc {
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

func testAccCheckObjectDestroy(s *terraform.State) error {
	config := GetTestProvider(true).Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_objectstorage_object" {
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

func testAccObjectConfig(bucketName, key, source string) string {
	return fmt.Sprintf(`
	resource "ncloud_objectstorage_bucket" "testing_bucket" {
		bucket_name			= "%[1]s"
	}

	resource "ncloud_objectstorage_object" "testing_object" {
		bucket				= ncloud_objectstorage_bucket.testing_bucket.bucket_name
		key 				= "%[2]s"
		source				= "%[3]s"	
	}`, bucketName, key, source)
}

func CreateTempFile(t *testing.T, content, key string) *os.File {
	tmpFile, err := os.CreateTemp("", key)
	if err != nil {
		t.Error("Error Occur: ", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Error("Error Occur: ", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Error("Error Occur: ", err)
	}

	return tmpFile
}
