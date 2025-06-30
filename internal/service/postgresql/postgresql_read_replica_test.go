package postgresql_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/vpostgresql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	. "github.com/terraform-providers/terraform-provider-ncloud/internal/acctest"
	"github.com/terraform-providers/terraform-provider-ncloud/internal/conn"
	postgresqlservice "github.com/terraform-providers/terraform-provider-ncloud/internal/service/postgresql"
)

func TestAccResourceNcloudPostgresqlReadReplica_vpc_basic(t *testing.T) {
	var postgresqlServerInstance vpostgresql.CloudPostgresqlServerInstance
	testName := fmt.Sprintf("tf-postgresqlrr-%s", acctest.RandString(5))
	resourceName := "ncloud_postgresql_read_replica.postgresql_rr"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPostgresqlReadReplicaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlReadReplicaConfig(testName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlReadReplicaExists(resourceName, &postgresqlServerInstance, TestAccProvider),
					resource.TestCheckResourceAttrSet(resourceName, "postgresql_instance_no"),
				),
			},
		},
	})
}

func testAccPostgresqlReadReplicaConfig(testPostgresqlName string) string {
	return fmt.Sprintf(`
resource "ncloud_vpc" "test_vpc" {
	name             = "%[1]s"
	ipv4_cidr_block  = "10.5.0.0/16"
}

resource "ncloud_subnet" "test_subnet" {
	vpc_no             = ncloud_vpc.test_vpc.vpc_no
	name               = "%[1]s"
	subnet             = "10.5.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test_vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
}

resource "ncloud_postgresql" "postgresql" {
	vpc_no            = ncloud_vpc.test_vpc.vpc_no
	subnet_no         = ncloud_subnet.test_subnet.id
	service_name      = "%[1]s"
	server_name_prefix = "testprefix"
	user_name         = "testusername"
	user_password     = "t123456789!a"
	client_cidr       = "0.0.0.0/0"
	database_name     = "test_db"
}

resource "ncloud_postgresql_read_replica" "postgresql_rr" {
	postgresql_instance_no = ncloud_postgresql.postgresql.id
}
`, testPostgresqlName)
}

func testAccCheckPostgresqlReadReplicaExists(n string, readreplica *vpostgresql.CloudPostgresqlServerInstance, provider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := provider.Meta().(*conn.ProviderConfig)
		postgresqlReadReplica, err := postgresqlservice.GetPostgresqlReadReplicaServer(context.Background(), config, resource.Primary.Attributes["postgresql_instance_no"], resource.Primary.Attributes["id"])
		if err != nil {
			return err
		}

		if postgresqlReadReplica != nil {
			*readreplica = *postgresqlReadReplica[0]
			return nil
		}

		return fmt.Errorf("postgresql read replica not found")
	}
}

func testAccCheckPostgresqlReadReplicaDestroy(s *terraform.State) error {
	config := TestAccProvider.Meta().(*conn.ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ncloud_postgresql_read_replica" {
			continue
		}
		instance, err := postgresqlservice.GetPostgresqlReadReplicaServer(context.Background(), config, rs.Primary.Attributes["postgresql_instance_no"], rs.Primary.Attributes["id"])
		if err != nil && !strings.Contains(err.Error(), "5001017") {
			return err
		}

		if instance != nil {
			return errors.New("postgresql read replica still exists")
		}
	}

	return nil
}
