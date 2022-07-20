package ncloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNcloudSourceCommitRepository(t *testing.T) {
	dataName := "data.ncloud_sourcecommit_repository.test-repo"
	resourceName := "ncloud_sourcecommit_repository.test-repo"
	repositoryName := getTestRepositoryName()
	repositoryDesc := fmt.Sprintf("description of %v", repositoryName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNcloudSourceCommitRepositoryConfig(repositoryName, repositoryDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(dataName),
					resource.TestCheckResourceAttrPair(dataName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataName, "filesafer", resourceName, "filesafer"),
					resource.TestCheckResourceAttrPair(dataName, "git_https", resourceName, "git_https"),
					resource.TestCheckResourceAttrPair(dataName, "git_ssh", resourceName, "git_ssh"),
					resource.TestCheckResourceAttrPair(dataName, "creator", resourceName, "creator"),
				),
			},
		},
	})
}

func testAccDataSourceNcloudSourceCommitRepositoryConfig(name string, description string) string {
	return fmt.Sprintf(`
resource "ncloud_sourcecommit_repository" "test-repo" {
	name = "%[1]s"
	description = "%[2]s"
	filesafer = false
}

data "ncloud_sourcecommit_repository" "test-repo" {
	name = ncloud_sourcecommit_repository.test-repo.name
}
`, name, description)
}
