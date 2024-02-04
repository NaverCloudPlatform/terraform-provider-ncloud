---
subcategory: "Hadoop"
---


# Data Source: ncloud_hadoop

This module can be useful for getting detail of Hadoop instance created before.

## Example Usage

#### Basic usage

The following example shows how to take Hadoop instance ID and obtain the data.

```hcl
data "ncloud_hadoop" "hadoop_by_id" {
  id = ncloud_hadoop.hadoop.id
}

data "ncloud_hadoop" "hadoop_by_filter" {
  filter {
    name = "id"
    values = [ncloud_hadoop.hadoop.id]
  }
}
```

## Argument Reference

* `id` - (Optional) The ID of the specific Hadoop to retrieve.
* `zone_code` - (Optional) The zone code of the specific Hadoop to retrieve.
* `vpc_no` - (Optional) The vpc ID of the specific Hadoop to retrieve.
* `subnet_no` - (Optional) The subnet ID of the specific Hadoop to retrieve.
* `cluster_name` - (Optional) The name of the specific Hadoop to retrieve.
* `server_name` - (Optional) The server name in server list of specific Hadoop to retrieve.
* `server_instance_no` - (Optional) The server ID in server list of the specific Hadoop to retrieve.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Reuired) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of Hadoop.
* `cluster_name` - The name of Hadoop.
* `cluster_type_code` - The type code of Hadoop.
* `version` - The version of Hadoop.
* `image_product_code` - The image product code of Hadoop.
* `hadoop_server_instance_list` - The server instance list of Hadoop.