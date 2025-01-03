---
subcategory: "PostgreSQL"
---

# Resource: ncloud_postgresql_users

Provides a PostgreSQL User list resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage

```terraform
resource "ncloud_vpc" "vpc" {
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "subnet" {
  vpc_no         = ncloud_vpc.vpc.vpc_no
  subnet         = cidrsubnet(ncloud_vpc.vpc.ipv4_cidr_block, 8, 1)
  zone           = "KR-2"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PUBLIC"
}

resource "ncloud_postgresql" "postgresql" {
  vpc_no             = ncloud_vpc.vpc.vpc_no
  subnet_no          = ncloud_subnet.subnet.id
  service_name       = "tf-postgresql"
  server_name_prefix = "name-prefix"
  user_name          = "username"
  user_password      = "password1!"
  client_cidr        = "0.0.0.0/0"
  database_name      = "db_name"
}

resource "ncloud_postgresql_users" "postgresql_users" {
	id = ncloud_postgresql.postgresql.id
	postgresql_user_list = [
		{
			name = "test1",
			password = "t123456789!",
			client_cidr = "0.0.0.0/0",
			is_replication_role = "false"
		},
		{
			name = "test2",
			password = "t123456789!",
			client_cidr = "0.0.0.0/0",
			is_replication_role = "false"
		}
	]
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the associated Postgresql Instance.
* `postgresql_user_list` - The list of users to add.
  * `name` - (Required) PostgreSQL User ID. Only English alphabets, numbers and special characters ( \ _ , - ) are allowed and must start with an English alphabet. Cannot include User ID. Min: 4, Max: 16
  * `password` - (Required) PostgreSQL User Password. At least one English alphabet, number and special character must be included. Certain special characters ( ` & + \ " ' / space ) cannot be used. Min: 8, Max: 20
  * `client_cidr` - (Required) Access Control (CIDR) of the client you want to connect to EX) Allow all access: 0.0.0.0/0, Allow specific IP access: 192.168.1.1/32, Allow IP band access: 192.168.1.0/24
  * `replication_role` - (Required) Replication Role or not (true/false).

## Import

### `terraform import` command

* PostgreSQL User can be imported using the `id`:`name`:`name`:... . For example:

```console
$ terraform import ncloud_postgresql_users.rsc_name 12345:name1:name2
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import PostgreSQL User using the `id`:`name`:`name`:... . For example:

```terraform
import {
    to = ncloud_postgresql_users.rsc_name
    id = "12345:name1:name2"
}
```
