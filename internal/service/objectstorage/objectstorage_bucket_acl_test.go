package objectstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/objectstorage"
)

func TestAccResourceNcloudObjectStorage_bucket_acl_basic(t *testing.T) {
	var aclOutput s3.GetBucketAclOutput
	bucketID := "https://kr.object.ncloudstorage.com/tfstate-backend"
	aclOptions := []string{string(awsTypes.BucketCannedACLPrivate),
		string(awsTypes.BucketCannedACLPublicRead),
		string(awsTypes.BucketCannedACLPublicReadWrite),
		string(awsTypes.BucketCannedACLAuthenticatedRead)}
	acl := aclOptions[acctest.RandIntRange(0, len(aclOptions)-1)]
	resourceName := "ncloud_objectstorage_bucket_acl.testing_acl"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketACLConfig(bucketID, acl),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBucketACLExists(resourceName, &aclOutput, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "bucket_id", bucketID),
					resource.TestCheckResourceAttr(resourceName, "rule", acl),
				),
			},
		},
	})
}

func testAccCheckBucketACLExists(n string, object *s3.GetBucketAclOutput, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		bucketID := resource.Primary.Attributes["bucket_id"]

		bucketName := objectstorage.BucketIDParser(bucketID)

		config := provider.Meta().(*conn.ProviderConfig)
		resp, err := config.Client.ObjectStorage.GetBucketAcl(context.Background(), &s3.GetBucketAclInput{
			Bucket: ncloud.String(bucketName),
		})

		if err != nil {
			return err
		}

		if resp != nil {
			object = resp
			return nil
		}

		return fmt.Errorf("Bucket ACL not found")
	}
}

func testAccBucketACLConfig(bucketID, acl string) string {
	return fmt.Sprintf(`
		resource "ncloud_objectstorage_bucket_acl" "testing_acl" {
			bucket_id				= "%[1]s"
			rule					= "%[2]s"
		}
	`, bucketID, acl)
}
