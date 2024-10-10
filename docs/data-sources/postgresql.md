---
subcategory: "PostgreSQL"
---

# Data Source: ncloud_postgresql

This moudule can be useful for getting detail of PostgreSQL created before.

~> **NOTE:** This only supports VPC environment.

## Example Usage

The following example shows how to take PostgreSQL instance ID and obtain the data.

```terraform
data "ncloud_postgresql" "by_id" {
    id = 12345
}
data "ncloud_postgresql" "by_name" {
    service_name = "example"
}
```

## Argument Reference

The following arguments are required:

* `id` - (Required) Postgresql instance number. Either `id` or `service_name` must be provided.
* `service_name` - (Required) Postgresql service name. Either `id` or `service_name` must be provided.

## Attributes Reference

This data source exports the following attributes in addition to the argument above:

* `region_code` - Region code.
* `vpc_no` - The ID of the associated VPC.
* `image_product_code` - The image product code of the PostgreSQL instance.
* `product_code` - Server specifications of the Cloud DB for PostgreSQL instance.
* `data_storage_type` - The type of data storage.
* `is_ha` - Whether using high availability of the specific PostgreSQL.
* `is_multi_zone` - Wheter using multi zone of the specific PostgreSQL.
* `is_storage_encryption` - Wheter data storage encryption is applied.
* `is_backup` - Wheter using backup of the specific PostgreSQL.
* `backup_file_retention_period` - The backup period of the PostgreSQL database.
* `backup_file_storage_count` -  Number of backup files kept.
* `backup_time` - The backup time fo the PostgreSQL database.
* `port` - Port of PostgreSQL instance.
* `engine_version_code` - The engien version of the specific PostgreSQL.
* `access_control_group_no_list` - The list of access control group number.
* `postgresql_config_list` - The list of config.
* `postgresql_server_list` - The list of PostgreSQL server instance.
  * `server_instance_no` - Server instance number.
  * `server_name` - Server name.
  * `server_role` - Server role code. M(Primary), H(Secondary)
  * `subnet_no` - Number of the associated Subnet.
  * `product_code` - Product code.
  * `is_public_subnet` - Public subnet status.
  * `public_domain` - Public domain.
  * `private_domain` - Private domain.
  * `private_ip` - Private IP.
  * `memory_size` - Available memory size.
  * `cpu_count` - CPU count.
  * `data_storage_size` - Storage size.
  * `used_data_storage_size` - Size of data storage in use.
  * `uptime` - Running start time.
  * `create_date` - Server create date.