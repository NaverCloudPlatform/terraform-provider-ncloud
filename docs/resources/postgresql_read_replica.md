---
subcategory: "PostgreSQL"
---

# Resource: ncloud_postgresql_read_replica

Provides a PostgreSQL instance resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage

```terraform
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
	postgresql_instance_no = ncloud_postgresql.postgresql.postgresql_instance_no
}
```

## Argument Reference

The following arguments are supported:

* `postgresql_instance_no` - (Required) The ID of the associated Postgresql Instance.
* `subnet_no` - (Optional, Required if `is_multi_zone` of MySQL Instance is true) The ID of the associate Subnet.

## Attribute Reference

In addition to all arguments above, the following attributes are exported

* `id` - PostgreSQL Read Replica Server Instance Number. 

# Import

### `terraform import` command

* PostgreSQL Read Replica can be imported using the `id`. For example:
```console
$ terraform import ncloud_postgresql_read_replica.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import PostgreSQL Read Replica using the `id`. For example:

```terraform
import {
    to = nlcoud_postgresql_read_replica.rsc_name
    id = "12345"
}
```