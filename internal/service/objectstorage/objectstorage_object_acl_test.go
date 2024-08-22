package objectstorage_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	awsTypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/objectstorage"
)

func TestAccResourceNcloudObjectStorage_object_acl_basic(t *testing.T) {
	objectID := objectstorage.ObjectIDGenerator(strings.ToLower("KR"), "tfstate-backend", "hello.md")
	aclOptions := []string{string(awsTypes.ObjectCannedACLPrivate),
		string(awsTypes.ObjectCannedACLPublicRead),
		string(awsTypes.ObjectCannedACLPublicReadWrite),
		string(awsTypes.ObjectCannedACLAuthenticatedRead)}
	acl := aclOptions[acctest.RandIntRange(0, len(aclOptions)-1)]
	resourceName := "ncloud_objectstorage_object_acl.testing_acl"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectACLConfig(objectID, acl),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectACLExists(resourceName, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "object_id", objectID),
					resource.TestCheckResourceAttr(resourceName, "rule", acl),
				),
			},
		},
	})
}

func testAccCheckObjectACLExists(n string, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		objectID := resource.Primary.Attributes["object_id"]

		_, _, bucketName, key := objectstorage.ObjectIDParser(objectID)

		config := provider.Meta().(*conn.ProviderConfig)
		resp, err := config.Client.ObjectStorage.GetObjectAcl(context.Background(), &s3.GetObjectAclInput{
			Bucket: ncloud.String(bucketName),
			Key:    ncloud.String(key),
		})
		if err != nil {
			return err
		}

		if resp != nil {
			return nil
		}

		return fmt.Errorf("Object ACL not found")
	}
}

func testAccObjectACLConfig(objectID, acl string) string {
	return fmt.Sprintf(`
	resource "ncloud_objectstorage_object_acl" "testing_acl" {
		object_id			= "%[1]s"
		rule				= "%[2]s" 
	}`, objectID, acl)
}
