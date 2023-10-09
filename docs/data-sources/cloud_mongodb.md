---
subcategory: "MongoDb"
---


# Data Source: ncloud_mongodb

Provides a Database Service MongoDb data.

## Example Usage

```hcl
data "ncloud_mongodb" "by_id" {
  id = ncloud_mongodb.mongodb.id
}

data "ncloud_mongodb" "by_filter" {
  filter {
    name = "instance_no"
    values = [ncloud_mongodb.mongodb.id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific MongoDb to retrieve.
* `service_name` - (Optional) The name of the specific MongoDb to retrieve.
* `cluster_type_code` - (Optional) The cluster Type of MongoDb.
* `image_product_code` - (Optional) The image product code of MongoDb.
* `engine_version` - (Optional) The engine version of the specific MongoDb to retrieve.
* `shard_count` - The number of MongoDb Shards.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.


## Attributes Reference

* `vpc_no` - The ID of the Vpc.
* `subnet_no` - The ID of the Subnet.
* `backup_file_retention_period` - The backup period of the MongoDb database.
* `backup_time` - The backup time of the MongoDb database.
* `image_product_code` - The image product code of the MongoDb instance.
* `server_instance_list` The list of MongoDb server instances.