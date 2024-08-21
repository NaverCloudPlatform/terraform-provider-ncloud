package objectstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/service/objectstorage"
)

func TestAccResourceNcloudObjectStorage_object_acl_basic(t *testing.T) {
	var aclOutput s3.GetObjectAclOutput
	objectID := "https://kr.object.ncloudstorage.com/tfstate-backend/hello.md"
	acl := "public-read-write"
	resourceName := "ncloud_objectstorage_object_acl.testing_acl"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectACLConfig(objectID, acl),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckObjectACLExists(resourceName, &aclOutput, GetTestProvider(true)),
					resource.TestCheckResourceAttr(resourceName, "object_id", objectID),
					resource.TestCheckResourceAttr(resourceName, "rule", acl),
				),
			},
		},
	})
}

func testAccCheckObjectACLExists(n string, object *s3.GetObjectAclOutput, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		objectID := resource.Primary.Attributes["object_id"]

		bucketName, key := objectstorage.ObjectIDParser(objectID)

		config := provider.Meta().(*conn.ProviderConfig)
		resp, err := config.Client.ObjectStorage.GetObjectAcl(context.Background(), &s3.GetObjectAclInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return err
		}

		if resp != nil {
			object = resp
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
