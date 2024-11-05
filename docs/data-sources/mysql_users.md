---
subcategory: "MySQL"
---

# Data Source: ncloud_mysql_users

Get a list of MySQL users.

~> **NOTE:** This only supports VPC environments.

## Example Usage

```terraform
data "ncloud_mysql_users" "all" {
    id = 12345
    filter {
        name = "name"
        values = ["test-user1"]
    }
    
    output_file = "users.json"
}

output "user_list" {
    values = {
        for user in data.ncloud_mysql_users.all.mysql_user_list:
            user.name => user.host_ip
    }
}
```

Outputs:
```terraform
user_list = {
    "test-user1": "192.168.0.1"
}
```

## Argument Reference

The following arguments are required:

* `id` - (Required) Mysql Users number. Either `id` or `mysql_instance_no` must be provided.
* `mysql_instance_no` - (Required) Mysql Instance No, either `id` or `mysql_instance_no` must be provided.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `vlaues` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the argument above: 

* `mysql_user_list` - The list of users to add .
  * `name` - MySQL User ID.
  * `host_ip` -  MySQL user host.
  * `authority` - MySQL User Authority.
  * `is_system_table_access` - MySQL User system table accessibility.
