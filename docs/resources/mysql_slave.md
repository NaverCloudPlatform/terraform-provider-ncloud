---
subcategory: "MySQL"
---

# Resource: ncloud_mysql_slave

Provides a MySQL slave resource.

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

resource "ncloud_mysql" "mysql" {
	subnet_no = data.ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
	server_name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!a"
	host_ip = "192.168.0.1"
	database_name = "test_db"
}

resource "ncloud_mysql_slave" "mysql_slave" {
	mysql_instance_no = ncloud_mysql.mysql.id
}
```

## Argument Reference

The following arguments are supported:

* `mysql_instance_no` - (Required) the ID of the associated Mysql Instance.
* `subnet_no` - (Optional, Required if `is_multi_zone` of MySQL Instance is true) The ID of the associate Subnet.

## Attribute Reference

In addition to all arguments above, the following attributes are exported

* `id` - MySQL Slave Server Instance Number.

# Import

### `terraform import` command

* MySQL Slave can be imported using the `id`. For example:
```console
$ terraform import ncloud_mysql_slave.rsc_name 12345
```

### `import` block

* In terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import MySQL Slave using the `id`. For example:

```terraform
import {
    to = ncloud_mysql_slave.rsc_name
    id = "12345"
}
```