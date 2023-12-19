package hadoop_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"

	"testing"
)

func TestAccDataSourceNcloudHadoopBucketbasic(t *testing.T) {
	dataName := "data.ncloud_hadoop_bucket.bucket"
	firstBucketName := "akj1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBucketConfig(),
				Check: resource.ComposeTestCheckFunc(
					TestAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "bucket_list.0", firstBucketName),
				),
			},
		},
	})
}

func testAccDataSourceBucketConfig() string {
	return fmt.Sprintf(`
data "ncloud_hadoop_bucket" "bucket" {

}
`)
}
