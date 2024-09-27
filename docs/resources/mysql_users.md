---
subcategory: "MySQL"
---

# Resource: ncloud_mysql_users

Provides a MySQL User list resource.

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

resource "ncloud_mysql_users" "mysql_users" {
	mysql_instance_no = ncloud_mysql.mysql.id
	mysql_user_list = [
		{
			name = "test1",
			password = "t123456789!",
			host_ip = "%%",
			authority = "READ"
		},
		{
			name = "test2",
			password = "t123456789!",
			host_ip = "%%",
			authority = "DDL"
		}
	]
}
```

## Argument Reference
The following arguments are supported:

* `mysql_instance_no` - (Required) The ID of the associated Mysql Instance.
* `mysql_user_list` - The list of users to add .
  * `name` - (Required) MySQL User ID. Only English alphabets, numbers and special characters ( \ _ , - ) are allowed and must start with an English alphabet. Min: 4, Max: 16
  * `password` - (Required) MySQL User Password. At least one English alphabet, number and special character must be included. Certain special characters ( ` & + \ " ' / space ) cannot be used. Min: 8, Max: 20
  * `host_ip` - (Required) MySQL user host. ex) Overall connection permitted: %, Connection by specific IPs permitted: 1.1.1.1, IP band connection permitted: 1.1.1.%
  * `authority` - (Required) MySQL User Authority. You can select `READ|CRUD|DDL`.

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `id` - Mysql User List number.(Mysql Instance number)

## Import

### `terraform import` command

* MySQL User can be imported using the `id`. For example:

```console
$ terraform import ncloud_mysql_users.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import MySQL User using the `id`. For example:

```terraform
import {
    to = ncloud_mysql_users.rsc_name
    id = "12345"
}
```