---
subcategory: "Redis"
---


# Data Source: ncloud_redis

Provides information about a Redis.

~> **NOTE:** This only supports VPC environment.

## Example Usage

```terraform
data "ncloud_redis" "example" {
  id = 1234567
}
```

## Argument Reference

The following arguments are required:

* `id` - (Required) Redis instance number.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `service_name` - Name of the Redis Service.
* `server_name_prefix` - Prefix name of the Redis Server.
* `vpc_no` - VPC number.
* `subnet_no` - Number of the associated Subnet.
* `config_group_no` - Redis Config Group number.
* `mode` - Configuration of Cloud DB for Redis. `CLUSTER` or `SIMPLE`.
* `image_product_code` - Image product code.
* `product_code` - Server specifications of the Cloud DB for Redis instance.
* `is_ha` - Whether is High Availability or not.
* `is_backup` - Backup status.
* `backup_file_retention_period` - Backup retention period.
* `backup_time` - The time of the backup runs.
* `backup_schedule` - Automatic or User-defined.
* `port` - Cloud Redis port.
* `region_code` - Region code.
* `access_control_group_no_list` - The ID list of the associated Access Control Group.
* `redis_server_list` - The list of the Redis server.
  * `redis_server_instance_no` - Redis Server instance number.
  * `redis_server_name` - Redis Server name.
  * `redis_server_role` - Stand Alone or Master or Slave.
  * `private_domain` - Private domain.
  * `memory_size` - Available memory size.
  * `os_memory_size` - OS memory size.
  * `uptime` - Running start time.
  * `create_date` - Server create date.
