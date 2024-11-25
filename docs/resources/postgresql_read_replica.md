---
subcategory: "PostgreSQL"
---

# Resource: ncloud_postgresql_read_replica

Provides a PostgreSQL instance resource.

~> **NOTE:** This resource only supports VPC environment. You can't create and delete more than one resource at the same time.

## Example Usage

```terraform
resource "ncloud_vpc" "vpc" {
	name             = "post-vpc"
	ipv4_cidr_block  = "10.5.0.0/16"
}

resource "ncloud_subnet" "subnet" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	name               = "post-subnet"
	subnet             = "10.5.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
}

resource "ncloud_postgresql" "postgresql" {
	vpc_no            = ncloud_vpc.vpc.vpc_no
	subnet_no         = ncloud_subnet.subnet.id
	service_name      = "tf-postgresql"
	server_name_prefix = "testprefix"
	user_name         = "testusername"
	user_password     = "t123456789!a"
	client_cidr       = "0.0.0.0/0"
	database_name     = "test_db"
}

resource "ncloud_postgresql_read_replica" "postgresql_rr" {
	postgresql_instance_no = ncloud_postgresql.postgresql.id
}
```

## Argument Reference

The following arguments are supported:

* `postgresql_instance_no` - (Required) The ID of the associated PostgreSQL instance.
* `subnet_no` - (Optional, Required if `multi_zone` of PostgreSQL instance is true) The ID of the associate Subnet. Not available in Neurocloud and `gov` site.

## Attribute Reference

In addition to all arguments above, the following attributes are exported

* `id` - PostgreSQL Read Replica server instance number. 
* `postgresql_server_list` - The list of the PostgreSQL server.
  * `server_instance_no` - Server instance number.
  * `server_name` - Server name.
  * `server_role` - Server role code. M(Primary), H(Secondary), S(Read Replica)
  * `product_code` - Product code.
  * `zone_code` - Zone code.
  * `subnet_no` - Number of the associated Subnet.
  * `public_subnet` - Public subnet status. (`true` or `false`)
  * `public_domain` - Public domain.
  * `private_domain` - Private domain.
  * `private_ip` - Private IP.
  * `data_storage_size` - Storage size.
  * `used_data_storage_size` - Size of data storage in use.
  * `cpu_count` - CPU count.
  * `memory_size` - Available memory size.
  * `uptime` - Running start time.
  * `create_date` - Server create date.

# Import

### `terraform import` command

* PostgreSQL Read Replica can be imported using the `postgresql_instance_no`:`id`. For example:
```console
$ terraform import ncloud_postgresql_read_replica.rsc_name 12345:24678
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import PostgreSQL Read Replica using the `postgresql_instance_no`:`id`. For example:

```terraform
import {
    to = nlcoud_postgresql_read_replica.rsc_name
    id = "12345:24678"
}
```
