package hadoop_test

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"

	"testing"
)

func TestAccDataSourceNcloudHadoopBucket_basic(t *testing.T) {
	dataName := "data.ncloud_hadoop_bucket.bucket"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBucketConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataName, "bucket_list.0", regexp.MustCompile(`^[a-z0-9]+[a-z0-9-.]+[a-z0-9]$`)),
				),
			},
		},
	})
}

func testAccDataSourceBucketConfig() string {
	return `
data "ncloud_hadoop_bucket" "bucket" {

}
`
}
