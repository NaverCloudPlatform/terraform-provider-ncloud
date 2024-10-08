---
subcategory: "MySQL"
---

# Data Source: ncloud_mysql_databases

Get a list of MySQL databases.

~> **NOTE:** This only supports VPC environments.

## Example Usage

```terraform
data "ncloud_mysql_databases" "all" {
    id = 12345
    filter {
        name = "name"
        values = ["testdb1"]
    }
    
    output_file = "databases.json"
}
output "database_list" {
    values = [for db in data.ncloud_mysql_databases.all.mysql_database_list : db.name]
}
```

Outputs:
```terraform
database_list = [
    "testdb1"
]
```

## Argument Reference

The following arguments are required:

* `mysql_instance_no` - (Required) The ID of the associated Mysql Instance.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `vlaues` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the argument above: 

* `id` - Mysql Databases number(Mysql Instance number).
* `mysql_database_list` - The list of databases to add .
  * `name` - MySQL Database Name.