package objectstorage_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func TestAccResourceNcloudObjectStorage_object_basic(t *testing.T) {
	bucket := "tfstate-backend"
	key := fmt.Sprintf("%s.md", acctest.RandString(5))
	resourceName := "ncloud_objectstorage_object.testing_object"
	content := "content for file upload testing"

	tmpFile := createTempFile(t, content, key)
	source := tmpFile.Name()
	defer os.Remove(source)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig(bucket, key, source),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectExists(resourceName, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^https:\/\/.*\.object\.ncloudstorage\.com\/[^\/]+\/[^\/]+\.*$`)),
					resource.TestCheckResourceAttr(resourceName, "bucket", bucket),
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
			Bucket: aws.String(resource.Primary.Attributes["bucket"]),
			Key:    aws.String(resource.Primary.Attributes["key"]),
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
			Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			Key:    aws.String(rs.Primary.Attributes["key"]),
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

func testAccObjectConfig(bucket, key, source string) string {
	return fmt.Sprintf(`
	resource "ncloud_objectstorage_object" "testing_object" {
		bucket				= "%[1]s"
		key 				= "%[2]s"
		source				= "%[3]s"	
	}`, bucket, key, source)
}

func createTempFile(t *testing.T, content, key string) *os.File {
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
