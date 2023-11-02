---
subcategory: "MongoDB"
---


# Data Source: ncloud_mongodb

Provides a Database Service MongoDB data.

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

* `id` - (Optional) The ID of the specific MongoDB to retrieve.
* `service_name` - (Optional) The name of the specific MongoDB to retrieve.
* `cluster_type_code` - (Optional) The cluster Type of MongoDB.
* `image_product_code` - (Optional) The image product code of MongoDB.
* `engine_version` - (Optional) The engine version of the specific MongoDB to retrieve.
* `shard_count` - The number of MongoDB Shards.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.


## Attributes Reference

* `vpc_no` - The ID of the Vpc.
* `subnet_no` - The ID of the Subnet.
* `backup_file_retention_period` - The backup period of the MongoDB database.
* `backup_time` - The backup time of the MongoDB database.
* `image_product_code` - The image product code of the MongoDB instance.
* `server_instance_list` The list of MongoDB server instances.

The `server_instance_list` object supports the following:

* `server_instance_no` - Server instance number
* `server_name` - Name of the server
* `cluster_role` - Cluster role
* `server_role` - Server role
* `region_code` - Region code
* `vpc_no` - The ID of the vpc
* `subnet_no` - The ID of the subnet
* `uptime` - Uptime
* `zone_code` - Zone code
* `private_domain` - Name of the private domain
* `public_domain` - Name of the public domain
* `memory_size` - Memory size
* `cpu_count` - the number of the virtual CPU
* `data_storage_size` - Data storage size
* `used_data_storage_size` - Used data storage size
* `product_code` - Product code
* `replica_set_name` - Name of the replica set
* `data_storage_type` - type of data storage