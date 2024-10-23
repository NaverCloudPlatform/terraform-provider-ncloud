---
subcategory: "MongoDB"
---

# Resource: ncloud_mongodb_users

Provides a MongoDB User list resources.

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

resource "ncloud_mongodb" "mongodb" {
	vpc_no = data.ncloud_vpc.test_vpc.vpc_no
	subnet_no = data.ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
	server_name_prefix = "ex-svr"
	cluster_type_code = "STAND_ALONE"
	user_name = "testuser"
	user_password = "t123456789!"
}

resource "ncloud_mongodb_users" "mongodb_users" {
	mongodb_instance_no = ncloud_mongodb.mongodb.id
	mongodb_user_list = [
		{
			name = "testuser1",
			password = "t123456789!",
			database_name = "testdb1",
			authority = "READ"
		},
		{
			name = "testuser2",
			password = "t123456789!",
			database_name = "testdb2",
			authority = "READ_WRITE"
		}
	]
}
```

## Argument Reference
The following arguments are supported:

* `mongodb_instance_no` - (Required) The ID of the associated MongoDB Instance.
* `mongodb_user_list` - The list of users to add.
  * `name` - (Required) MongoDB User ID. Allows only alphabets, numbers and underbar (_). Must start with an alphabetic character. Min: 4, Max: 16
  * `password` - (Required) MongoDB User Password. At least one English alphabet, number and special character must be included. Certain special characters ( ` & + \ " ' / space ) cannot be used. Min: 8 , Max: 20
  * `database_name` - (Required) MongoDB Database Name to add MongoDB User. Allows only alphabets, numbers and underbar (_). Must start with an alphabetic character. Min: 4 , Max: 30
  * `authority` - (Required) MongoDB User Authority. You can select `READ|READ_WRITE`.

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `id` - MongoDB User list number. (MongoDB Instance number)

## Import

### `terraform import` command

* MongoDB User can be imported using the `id`. For example:

```console
$ terraform import ncloud_mongodb_users.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import MongoDB User using the `id`. For example:

```terraform
import {
    to = ncloud_mongodb_users.rsc_name
    id = "12345"
}
```