---
subcategory: "MongoDB"
---

# Data Source: ncloud_mongodb_users

Get a list of MongoDB users.

~> **NOTE:** This only supports VPC environments.

## Example Usage

```terraform 
data "ncloud_mongodb_uers" "all" {
    id = 12345
    filter {
        name = "name"
        values = ["test-user1"]
    }

    output_file = "users.json"
}

output "user_list" {
    values = {
        for user in data.ncloud_mongodb_users.all.mongodb_user_list:
            user.name => user.database_name
    }
}
```

Outputs:
```terraform
user_list = {
    "test-user1": "test-db1"
}
```

## Argument Reference

The following arguments are required:

* `id` - (Required) MongoDB Users number. Either `id` or `mongodb_instance_no` must be provided.
* `mongodb_instance_no` - (Required) MongoDB Instance No, either `id` or `mongodb_instance_no` must be provided.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the argument above:

* `mongodb_user_list` - The list of users to add.
  * `name` - MongoDB User ID.
  * `password` - MongoDB User Password.
  * `database_name` - MongoDB Database Name that MongoDB User belongs to.
  * `authority` - Mongodb User Authority.