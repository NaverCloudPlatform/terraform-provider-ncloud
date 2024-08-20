package objectstorage_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
)

func TestAccResourceNcloudObjectStorage_bucket_basic(t *testing.T) {
	var bucket types.Bucket
	bucketName := "tf-state-new8"
	resourceName := "ncloud_objectstorage.testing_bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig(bucketName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBucketExists(resourceName, &bucket, GetTestProvider(true)),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
				),
			},
		},
	})
}

func testAccCheckBucketExists(n string, objectStorage *types.Bucket, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		resp, err := config.Client.ObjectStorage.ListBuckets(context.Background(), &s3.ListBucketsInput{})

		if err != nil {
			return err
		}

		for _, bucket := range resp.Buckets {
			if *bucket.Name == resource.Primary.Attributes["bucket_name"] {
				*objectStorage = bucket
				return nil
			}
		}

		return fmt.Errorf("Bucket not found")

	}
}

func testAccCheckBucketDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_objectstorage" {
			continue
		}

		config := GetTestProvider(true).Meta().(*conn.ProviderConfig)
		resp, err := config.Client.ObjectStorage.ListBuckets(context.Background(), &s3.ListBucketsInput{})

		if err != nil {
			return err
		}

		for _, bucket := range resp.Buckets {
			if *bucket.Name == rs.Primary.Attributes["bucket_name"] {
				return fmt.Errorf("Bucket found")
			}
		}
	}

	return nil
}

func testAccBucketConfig(bucketName string) string {
	return fmt.Sprintf(`
	resource "ncloud_objectstorage" "testing_bucket" {
		bucket_name				= "%[1]s"
	}`, bucketName)
}
