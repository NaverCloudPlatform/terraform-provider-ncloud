---
subcategory: "Mssql"
---


# Data Source: ncloud_mssql

This module can be useful for getting detail of MSSQL created before.

## Example Usage

#### Basic (VPC)

The following example shows how to take MSSQL instance ID and obtain the data.

```hcl
data "ncloud_mssql" "test" {
	id = ncloud_mssql.mssql.id
	service_name = ncloud_mssql.mssql.service_name
	is_ha = ncloud_mssql.mssql.is_ha
	is_multi_zone = ncloud_mssql.mssql.is_multi_zone
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific MSSQL to retrieve.
* `service_name` - (Optional) The name of the specific MSSQL to retrieve.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.
## Attributes Reference

* `vpc_no` - The ID of the VPC.
* `subnet_no` - The ID of the Subnet.
* `is_ha` - (Required) Whether using high availability of the specific MSSQL to retrieve.
* `is_multi_zone` - (Computed) Whether using multi zone of the specific MSSQL to retrieve.
* `is_backup` - (Computed) Whether using backup of the specific MSSQL to retrieve.
* `backup_file_retention_period` - The backup period of the MSSQL database.
* `backup_time` - The backup time of the MSSQL database.
* `image_product_code` - The image product code of the MSSQL instance.
* `cloud_mssql_server_instance_list` The list of MSSQL server instances.
