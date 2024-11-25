---
subcategory: "PostgreSQL"
---

# Resource: ncloud_postgresql_databases

Provides a PostgreSQL Database list resource.

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
			replication_role = "false"
		},
		{
			name = "test2",
			password = "t123456789!",
			client_cidr = "0.0.0.0/0",
			replication_role = "false"
		}
	]
}

resource "ncloud_postgresql_databases" "postgresql_databases" {
	id = ncloud_postgresql.postgresql.id
	postgresql_database_list = [
		{
			name = "testdb1",
			owner = ncloud_postgresql_users.postgresql_users.postgresql_user_list[0].name
		},
		{
			name = "testdb2",
			owner = ncloud_postgresql_users.postgresql_users.postgresql_user_list[1].name
		}
	]
}
```

## Argument Reference
The following arguments are supported:

* `id` - (Required) The ID of the associated Postgresql Instance.
* `postgresql_database_list` - The list of databases to add.
    * `name` - (Required) Database name to create. Only English alphabets, numbers and special characters ( \ _ , - ) are allowed and must start with an English alphabet. Min: 1, Max: 30
    * `owner` - (Required) User ID to manage the database.

## Import

### `terraform import` command

* PostgreSQL Database can be imported using the `id`:`name`:`name`:... . For example:

```console
$ terraform import ncloud_postgresql_databases.rsc_name 12345:name1:name2
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import PostgreSQL Database using the `id`:`name`:`name`:... . For example:

```terraform
import {
    to = ncloud_postgresql_databases.rsc_name
    id = "12345:name1:name2"
}
```
