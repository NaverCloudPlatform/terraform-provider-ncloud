---
subcategory: "Mssql"
---


# Data Source: ncloud_mssql

This module can be useful for getting detail of MSSQL created before.

~> **NOTE:** This only supports VPC environment.

## Example Usage

The following example shows how to take MSSQL instance ID and obtain the data.

```terraform
data "ncloud_mssql" "by_id" {
  id = 1234567
}

data "ncloud_mssql" "by_name" {
  service_name = "example"
}
```

## Argument Reference

The following arguments are required:

* `id` - (Required) MSSQL instance number. Either `id` or `service_name` must be provided.
* `service_name` - (Required) MSSQL service name. Either `id` or `service_name` must be provided.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `region_code` - Region code.
* `vpc_no` - The ID of the associated VPC. 
* `image_product_code` - The image product code of the MSSQL instance.
* `data_storage_type` - The type of data storage.
* `is_ha` - Whether using high availability of the specific MSSQL.
* `backup_file_retention_period` - The backup period of the MSSQL database.
* `backup_time` - The backup time of the MSSQL database.
* `config_group_no` - MSSQL config group number.
* `port` - Port of MSSQL instance.
* `engine_version` - The engine version of the specific MSSQL.
* `character_set_name` - DB character set name.
* `access_control_group_no_list` - The list of access control group number.
* `mssql_server_list` The list of MSSQL server instance.
  * `server_instance_no` - Server instance number.
  * `server_name` - Name of the server.
  * `server_role` - Server role code. ex) M(Principal), H(Mirror)
  * `zone_code` - Zone code.
  * `subnet_no` - The ID of the associated Subnet.
  * `product_code` - Product code.
  * `is_public_subnet` - Public subnet status.
  * `private_domain` - Name of the private domain.
  * `public_domain` - Name of the public domain.
  * `memory_size` - Memory size.
  * `cpu_count` - the number of the virtual CPU.
  * `data_storage_size` - Data storage size.
  * `used_data_storage_size` - Size of data storage in use.
  * `uptime` - Running start time.
  * `create_date` - Server create date.
