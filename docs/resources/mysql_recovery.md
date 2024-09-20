---
subcategory: "MySQL"
---

# Resource: ncloud_mysql_recovery

Provides a MySQL instance resource.
 
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

resource "ncloud_mysql_recovery" "mysql_recovery" {
	mysql_instance_no = ncloud_mysql.mysql.id
	recovery_server_name = "test-recovery"
	file_name = "20210722"
}
```

## Argument Reference

The following arguments are supported:

* `mysql_instance_no` - (Required) the ID of the associated Mysql Instance.
* `subnet_no` - (Optional, Required if `is_multi_zone` of MySQL Instance is true) The ID of the associate Subnet.
* `recovery_server_name` (Required) Recovery server name to create. In order to prevent overlapping host names, random text is added. Can comprise only lower-case English alphabets, numbers and dash ( - ). The first letter must be an English alphabet and the last letter must be an English alphabet or a number. Min: 3, Max: 25
* `file_name` - (Optional, One of `file_name` and `recovery_time` is required) The name of backup file. If you enter `file_name`, ignore the entry of `recovery_time`. 
* `recovery_time` - (Optional, One of `file_name` and `recovery_time` is required) The time of recovery. If you enter `recovery_time`, ignore the entry of `file_name`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported

* `id` - MySQL Recovery Server Instance Number.

# Import

### `terraform import` command

* MySQL Recovery can be imported using the `id`. For example:
```console
$ terraform import ncloud_mysql_recovery.rsc_name 12345
```

### `import` block

* In terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import MySQL Recovery using the `id`. For example:

```terraform
import {
    to = ncloud_mysql_recovery.rsc_name
    id = "12345"
}
```