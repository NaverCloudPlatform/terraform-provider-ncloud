---
subcategory: "Hadoop"
---


# Data Source: ncloud_hadoops

This resource is useful for look up the list of Hadoop in the region.

## Example Usage

```hcl
data "ncloud_hadoops" "hadoops_by_id" {
  id = ncloud_hadoop.hadoop.id
}

data "ncloud_hadoops" "hadoops_by_filter" {
  filter {
    name = "id"
    values = [ncloud_hadoop.hadoop.id]
  }
}
```

## Argumnet Reference

The following arguments are supported:

* `id` - (Optional) The ID of the Hadoop to retrieve.
* `zone_code` - (Optional) The zone code of the specific Hadoop to retrieve.
* `vpc_no` - (Optional) The VPC ID of the specific Hadoop to retrieve.
* `subnet_no` - (Optional) The subnet ID of the specific Hadoop to retrieve.
* `cluster_name` - (Optional) The name of the specific Hadoop to retrieve.
* `server_name` - (Optional) The server name in server list of specific Hadoop to retrieve.
* `server_instance_no` - (Optional) The server ID in server list of the specific Hadoop to retrieve.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Reuired) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

The following attrivutes are exported:

* `id` - The ID of hadoops.
* `hadoops` - The list of Hadoop.

### Hadoop Reference 

`hadoops` are also exported with the following attributes, where relevant: Each element supports the following:

* `id` - The ID of Hadoop.
* `cluster_name` - The name of Hadoop.
* `cluster_type_code` - The type code of Hadoop.
* `version` - The version of Hadoop.
* `image_product_code` - The image product code of Hadoop.
* `hadoop_server_instance_list` - The server instance list of Hadoop.