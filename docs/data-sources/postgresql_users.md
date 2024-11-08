---
subcategory: "PostgresSQL"
---

# Data Source: ncloud_postgresql_users

Get a list of PostgresQL users.

~> **NOTE:** This only supports VPC environments.

## Example Usage

```terraform
data "ncloud_postgresql_users" "all" {
    id = 12345
    filter {
        name = "name"
        values = ["test-user1"]
    }

    output_file = "users.json"
}

output "user_list" {
    values = {
        for user in data.ncloud_postgresql_users.all.postgresql_user_list:
            user.name => user.cidr
    }
}
```


Outputs:
```terraform
user_list = {
    "test-user1": "0.0.0.0/0"
}
```

## Argument Reference

The following arguments are required:

* `postgresql_instance_no` - (Required) Postgresql Instance No.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in adddition to the argument above:

* `postgresql_user_list` - The list of users to add.
  * `name` - PostgreSQL User ID.
  * `password` - PostgreSQL User Password.
  * `client_cidr` - Access Control (CIDR) of the client you want to connect 
  * `is_replication_role` - Replication Role or not