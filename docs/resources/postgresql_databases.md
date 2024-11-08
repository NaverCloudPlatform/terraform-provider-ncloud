---
subcategory: "PostgreSQL"
---

# Resource: ncloud_postgresql_databases

Provides a PostgreSQL Database list resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage

```terraform
resource "ncloud_vpc" "test" {
    ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "test" {
  vpc_no         = ncloud_vpc.test.vpc_no
  subnet         = cidrsubnet(ncloud_vpc.test.ipv4_cidr_block, 8, 1)
  zone           = "KR-2"
  network_acl_no = ncloud_vpc.test.default_network_acl_no
  subnet_type    = "PUBLIC"
}

resource "ncloud_postgresql" "postgresql" {
  subnet_no = ncloud_subnet.test.id
  service_name = "tf-postgresql"
  server_name_prefix = "name-prefix"
  user_name = "username"
  user_password = "password1!"
  client_cidr = "0.0.0.0/0"
  database_name = "db_name"
}

resource "ncloud_postgresql_users" "postgresql_users" {
	postgresql_instance_no = ncloud_postgresql.postgresql.postgresql_instance_no
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

resource "ncloud_postgresql_databases" "postgresql_databases" {
	postgresql_instance_no = ncloud_postgresql.postgresql.postgresql_instance_no
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

* `postgresql_instance_no` - (Required) The ID of the associated Postgresql Instance.
* `postgresql_database_list` - The list of databases to add.
    * `name` - (Required) Database name to create. Only English alphabets, numbers and special characters ( \ _ , - ) are allowed and must start with an English alphabet. Min: 1, Max: 30
    * `owner` - (Required) User ID to manage the database.

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `id` - Postgresql Database List number.(Postgresql Instance number)

## Import

### `terraform import` command

* PostgreSQL Database can be imported using the `id`. For example:

```console
$ terraform import ncloud_postgresql_databases.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import PostgreSQL Database using the `id`. For example:

```terraform
import {
    to = ncloud_postgresql_databases.rsc_name
    id = "12345"
}
```