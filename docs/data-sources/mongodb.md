---
subcategory: "MongoDB"
---


# Data Source: ncloud_mongodb

Provides a Database Service MongoDB data.

~> **NOTE:** This only supports VPC environment.

## Example Usage

```terraform
data "ncloud_mongodb" "example-id" {
  id = 1234567
}

data "ncloud_mongodb" "example-name" {
  service_name = "example"
}
```

## Argument Reference

The following arguments are required:

* `id` - (Required) MongoDb instance number. Either `id` or `service_name` must be provided.
* `service_name` - (Required) MongoDb service name. Either `id` or `service_name` must be provided.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `vpc_no` - The ID of the Vpc.
* `subnet_no` - The ID of the Subnet.
* `cluster_type_code` - The cluster Type of MongoDB.
* `engine_version` - The engine version of the specific MongoDB to retrieve.
* `image_product_code` - The image product code of MongoDB.
* `backup_file_retention_period` - The backup period of the MongoDB.
* `backup_time` - The backup time of the MongoDB.
* `shard_count` - The number of MongoDB Shards.
* `data_storage_type` - Type of data storage.
* `member_port` - TCP port number for access to the MongoDB Member Server.
* `arbiter_port` - TCP port number for access to the MongoDB Arbiter Server.
* `mongos_port` - TCP port number for access to the MongoDB Mongos Server.
* `config_port` - TCP port number for access to the MongoDB Config Server.
* `compress_code` - MongoDB Data Compression Algorithm Code.
* `region_code` - Region code.
* `zone_code` - Zone code.
* `access_control_group_no_list` - The ID list of the associated Access Control Group.
* `mongodb_server_list` - The list of the MongoDB server.
  * `server_instance_no` - Server instance number.
  * `server_name` - Name of the server.
  * `server_role` - Server role.
  * `cluster_role` - Cluster role.
  * `product_code` - Product code.
  * `private_domain` - Name of the private domain.
  * `public_domain` - Name of the public domain.
  * `replica_set_name` - Name of the replica set.
  * `memory_size` - Memory size.
  * `cpu_count` - the number of the virtual CPU.
  * `data_storage_size` - Data storage size.
  * `uptime` - Running start time.
  * `create_date` - Server create data.
