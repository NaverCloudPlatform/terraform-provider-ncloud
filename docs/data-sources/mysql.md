---
subcategory: "Mysql"
---


# Data Source: ncloud_mysql

This module can be useful for getting detail of MySQL created before.

## Example Usage

#### Basic (VPC)

The following example shows how to take MySQL instance ID and obtain the data.

```hcl
data "ncloud_mysql" "by_id" {
  id = ncloud_mysql.mysql.id
}

data "ncloud_mysql" "by_filter" {
  filter {
    name = "instance_no"
    values = [ncloud_mysql.mysql.id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific MySQL to retrieve.
* `service_name` - (Optional) The name of the specific MySQL to retrieve.
* `engine_version` - (Optional) The engine version of the specific MySQL to retrieve.
* `is_ha` - (Optional) Whether using high availability of the specific MySQL to retrieve.
* `is_multi_zone` - (Optional) Whether using multi zone of the specific MySQL to retrieve.
* `is_backup` - Whether using backup of the specific MySQL to retrieve.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.
## Attributes Reference

* `subnet_no` - The ID of the Subnet.
* `backup_file_retention_period` - The backup period of the MySQL database.
* `backup_time` - The backup time of the MySQL database.
* `image_product_code` - The image product code of the MySQL instance.
* `cloud_mysql_server_instance_list` The list of Mysql server instances.