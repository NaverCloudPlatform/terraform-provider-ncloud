---
subcategory: "MySQL"
---

# Resource: ncloud_mysql_recovery

Provides a MySQL recovery resource.
 
~> **NOTE:** This resource only supports VPC environment. Only one server related resource (ncloud_mysql_slave, ncloud_mysql_recovery) can be created or deleted at a time.

## Example Usage

```terraform
resource "ncloud_vpc" "vpc" {
  name             = "mysql-vpc"
  ipv4_cidr_block  = "10.5.0.0/16"
}

resource "ncloud_subnet" "subnet" {
  vpc_no             = ncloud_vpc.vpc.vpc_no
  name               = "mysql-subnet"
  subnet             = "10.5.0.0/24"
  zone               = "KR-2"
  network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
  subnet_type        = "PUBLIC"
}

resource "ncloud_mysql" "mysql" {
  subnet_no          = data.ncloud_subnet.subnet.id
  service_name       = "tf-mysql"
  server_name_prefix = "testprefix"
  user_name          = "testusername"
  user_password      = "t123456789!a"
  host_ip            = "192.168.0.1"
  database_name      = "test_db"
}

resource "ncloud_mysql_recovery" "mysql_recovery" {
  mysql_instance_no    = ncloud_mysql.mysql.id
  recovery_server_name = "test-recovery"
  file_name            = "20210722"
}
```

## Argument Reference

The following arguments are supported:

* `mysql_instance_no` - (Required) the ID of the associated Mysql Instance.
* `recovery_server_name` (Required) Recovery server name to create. In order to prevent overlapping host names, random text is added. Can comprise only lower-case English alphabets, numbers and dash ( - ). The first letter must be an English alphabet and the last letter must be an English alphabet or a number. Min: 3, Max: 25
* `subnet_no` - (Optional, Required if `is_multi_zone` of MySQL Instance is true) The ID of the associate Subnet. Not available in Neurocloud and `gov` site.
* `file_name` - (Optional, One of `file_name` and `recovery_time` is required) The name of backup file.
* `recovery_time` - (Optional, One of `file_name` and `recovery_time` is required) The time of recovery.

## Attribute Reference

In addition to all arguments above, the following attributes are exported

* `id` - MySQL Recovery Server Instance Number.
* `mysql_server_list` - The list of the MySQL server.
  * `server_instance_no` - Server instance number.
  * `server_name` - Server name.
  * `server_role` - Server role code. ex) M(Master), H(Standby Master)
  * `zone_code` - Zone code.
  * `subnet_no` - The ID of the associated Subnet.
  * `product_code` - Product code.
  * `is_public_subnet` - Public subnet status.
  * `private_domain` - Private domain.
  * `public_domain` - Public domain.
  * `memory_size` - Available memory size.
  * `cpu_count` - CPU count.
  * `data_storage_size` - Storage size.
  * `used_data_storage_size` - Size of data storage in use.
  * `uptime` - Running start time.
  * `create_date` - Server create date.

# Import

### `terraform import` command

* MySQL Recovery can be imported using the `mysql_instance_no`:`id`. For example:
```console
$ terraform import ncloud_mysql_recovery.rsc_name 12345:24678
```

### `import` block

* In terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import MySQL Recovery using the `mysql_instance_no`:`id`. For example:

```terraform
import {
    to = ncloud_mysql_recovery.rsc_name
    id = "12345:24678"
}
```
