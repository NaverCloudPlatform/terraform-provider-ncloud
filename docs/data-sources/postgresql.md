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
* `image_product_code` - The image product code of the instance.
* `generation_code` - The generation code of the image.
* `engine_version` - The engien version.
* `ha` - Whether using high availability. (`true` or `false`)
* `multi_zone` - Whether using multi zone. (`true` or `false`)
* `data_storage_type` - The type of data storage.
* `storage_encryption` - Whether using storage encryption. (`true` or `false`)
* `backup` - Whether using backup. (`true` or `false`)
* `backup_file_retention_period` - Backup file retention period.
* `backup_time` - Backup time.
* `port` - Port of PostgreSQL instance.
* `access_control_group_no_list` - The list of access control group number.
* `postgresql_config_list` - The list of config.
* `postgresql_server_list` - The list of PostgreSQL server instance.
  * `server_instance_no` - Server instance number.
  * `server_name` - Server name.
  * `server_role` - Server role code. M(Primary), H(Secondary), S(Read Replica)
  * `product_code` - Product code.
  * `zone_code` - Zone code.
  * `subnet_no` - Number of the associated Subnet.
  * `public_subnet` - Public subnet status. (`true` or `false`)
  * `public_domain` - Public domain.
  * `private_domain` - Private domain.
  * `private_ip` - Private IP.
  * `data_storage_size` - Storage size.
  * `used_data_storage_size` - Size of data storage in use.
  * `cpu_count` - CPU count.
  * `memory_size` - Available memory size.
  * `uptime` - Running start time.
  * `create_date` - Server create date.
