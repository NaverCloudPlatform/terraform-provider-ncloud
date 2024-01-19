---
subcategory: "MySQL"
---


# Data Source: ncloud_mysql

This module can be useful for getting detail of MySQL created before.

~> **NOTE:** This only supports VPC environment.

## Example Usage

The following example shows how to take MySQL instance ID and obtain the data.

```terraform
data "ncloud_mysql" "by_id" {
  id = 1234567
}

data "ncloud_mysql" "by_filter" {
  service_name = "example"
}
```

## Argument Reference

The following arguments are required:

* `id` - (Required) Mysql instance number. Either `id` or `service_name` must be provided.
* `service_name` - (Required) Mysql service name. Either `id` or `service_name` must be provided.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `vpc_no` - The ID of the associated VPC. 
* `subnet_no` - The ID of the associated Subnet.
* `region_code` - Region code.
* `zone_code` - Zone code.
* `image_product_code` - The image product code of the MySQL instance.
* `data_storage_type` - The type of data storage.
* `is_ha` - Whether using high availability of the specific MySQL.
* `is_multi_zone` - Whether using multi zone of the specific MySQL.
* `is_storage_encryption` - Whether data storage encryption is applied.
* `is_backup` -  Whether using backup of the specific MySQL.
* `backup_file_retention_period` - The backup period of the MySQL database.
* `backup_time` - The backup time of the MySQL database.
* `port` - Port of MySQL instance.
* `engine_version_code` - The engine version of the specific MySQL.
* `access_control_group_no_list` - The list of access control group number.
* `mysql_config_list` - The list of config.
* `mysql_server_list` The list of MySQL server instance.
  * `server_instance_no` - Server instance number.
  * `server_name` - Name of the server.
  * `server_role` - Server role.
  * `product_code` - Product code.
  * `is_public_subnet` - Public subnet status.
  * `private_domain` - Name of the private domain.
  * `public_domain` - Name of the public domain.
  * `memory_size` - Memory size.
  * `cpu_count` - the number of the virtual CPU.
  * `data_storage_size` - Data storage size.
  * `used_data_storage_size` - Size of data storage in use.
  * `uptime` - Running start time.
  * `create_date` - Server create data.
