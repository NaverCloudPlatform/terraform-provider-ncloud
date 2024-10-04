---
subcategory: "MySQL"
---

# Resource: ncloud_mysql_databases

Provides a MySQL Database list resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage

```terraform
resource "ncloud_vpc" "test_vpc" {
	ipv4_cidr_block  = "10.5.0.0/16"
}
resource "ncloud_subnet" "test_subnet" {
	vpc_no             = ncloud_vpc.test_vpc.vpc_no
	subnet             = "10.5.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test_vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
}
resource "ncloud_mysql" "mysql" {
	subnet_no = ncloud_subnet.test_subnet.id
	service_name = "tf-mysql"
	server_name_prefix = "testprefix"
	user_name = "testusername"
	user_password = "t123456789!a"
	host_ip = "192.168.0.1"
	database_name = "test_db"
}
resource "ncloud_mysql_databases" "mysql_db" {
	mysql_instance_no = ncloud_mysql.mysql.id
	mysql_database_list = [
		{
			name = "testdb1"
		},
		{
			name = "testdb2"
		}
	]
}
```

## Argument Reference
The following arguments are supported:

* `mysql_instance_no` - (Required) The ID of the associated Mysql Instance.
* `mysql_database_list` - The list of databases to add .
  * `name` - (Required) Database name to create. Only English alphabets, numbers and special characters ( \ _ , - ) are allowed and must start with an English alphabet. Min: 1, Max: 30

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `id` - Mysql Database List number.(Mysql Instance number)

## Import

### `terraform import` command

* MySQL Database can be imported using the `id`. For example:

```console
$ terraform import ncloud_mysql_databases.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import MySQL Database using the `id`. For example:

```terraform
import {
    to = ncloud_mysql_databases.rsc_name
    id = "12345"
}
```