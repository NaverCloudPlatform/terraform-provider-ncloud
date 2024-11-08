---
subcategory: "PostgresSQL"
---

# Data Source: ncloud_postgresql_databases

Get a list of PostgresQL databases.

~> **NOTE:** This only supports VPC environments.

## Example Usage

```terraform
data "ncloud_postgresql_databases" "all" {
    id = 12345
    filter {
        name = "name"
        vlaues = ["test-db1"]
    }

    output_file = "databasess.json"
}

output "database_list" {
    values = {
        for database in data.ncloud_postgresql_databases.all.postgresqul_database_list:
            database.name => database.owner
    }
}
```

Outputs:
```terraform
database_list = {
    "test-db1": "test-user1"
}
```

## Argument Reference

The following arguments are required:

* `postgresql_instance_no` - (Required) Postgresql Instance No.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in adddition to the argument above:

* `postgresql_database_list` - The list of databases to add.
  * `name` - PostgreSQL Databases ID.
  * `owner` - User ID to manage the database.