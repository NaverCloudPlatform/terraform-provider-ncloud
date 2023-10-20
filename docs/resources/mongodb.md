---
subcategory: "Database Service"
---


# Resource: ncloud_mongodb

Provides a Database Service MongoDB resource.

## Example Usage

```hcl
resource "ncloud_vpc" "vpc" {
  name            = "vpc"
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.0.1.0/24"
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "subnet-01"
  usage_type     = "GEN"
}

resource "ncloud_mongodb" "mongodb" {
  subnet_no = ncloud_subnet.subnet.id
  service_name = "sample-mongodb"
  user_name = "username"
  user_password = "password1!"
  cluster_type_code = "STAND_ALONE"
}
```


## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name to create.
* `user_name` - (Required) MongoDB User ID
* `user_password` - (Required) MongoDB User Password
* `vpc_no` - (Required) The ID of the associated Vpc.
* `subnet_no` - (Required) The ID of the associated Subnet.
* `cluster_type_code` - (Required) MongoDB cluster type code determines the cluster type of MongoDB. Options: STAND_ALONE | SINGLE_REPLICA_SET | SHARDED_CLUSTER
* `image_product_code` - (Optional) MongoDB image product code, cloudMongoDBImageProductCode can be acquired as a productCode in getCloudMongoDBImageProductList action If not entered, it is created as a default value.
* `member_product_code` - (Optional) 
* `arbiter_product_code` - (Optional) 
* `mongos_product_code` - (Optional) 
* `config_product_code` - (Optional) 
* `data_storage_type_code` - (Optional) 
* `shard_count` - (Optional) The number of MongoDB Shards. If sharding is used, the number of shards can be selected. For initial configurations, only two and three are selectable.
You can enter the ClusterType only when it is Sharding. Default: 2
* `member_server_count` - (Optional) The number of MongoDB Member Servers, it is possible to select the number of member servers per Replica Set (for each shard in the case of Sharding).
It can be selected between 3 and 7 units including the Arbiter server. Default: 3
* `arbiter_server_count` - (Optional) The number of MongoDB Arbiter servers. You can select whether to use the Arbiter server per Replica Set (for each shard in the case of Sharding). Up to one Arbiter server can be selected. The Arbiter server is provided with a minimum configurable spec. Default: 0
* `mongos_server_count` - (Optional) The number of MongoDB Mongos servers. If sharding is used, the number of mongos servers can be selected. Default: 2
* `config_server_count` - (Optional) The number of MongoDB Config servers. If sharding is used, the config server's logarithm can be selected. Default: 3
* `backup_file_retention` - (Optional) Backups are performed daily and backup files are stored in separate backup storage. Fees are charged based on the space used. Default: 1(1 day)
* `backup_time` - (Optional) You can set the time when backup is performed. Default: 02:00
* `data_storage_type_code` - (Optional) Data storage type. If `generationCode` is `G2`, You can select `SSD|HDD`, else if `generationCode` is `G3`, you can select CB1. Default : SSD in G2, CB1 in G3
* `arbiter_port` - (Optional) TCP port number for access to the MongoDB Arbiter Server.  Default: 17017
* `member_port` - (Optional) TCP port number for access to the MongoDB Member Server.  Default: 17017
* `mongos_port` - (Optional) TCP port number for access to the MongoDB Mongos Server.  Default: 17017
* `config_port` - (Optional) TCP port number for access to the MongoDB Config Server.  Default: 17017
* `compress_code` - (Optional) MongoDB Data Compression Algorithm Code allows you to select data compression algorithms provided by MongoDB.
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of server instance.
* `access_control_group_no_list` - The ID of the ACG.

## Import

MongoDB Instance can be imported using the id(service_name), e.g.,

```
$ terraform import ncloud_mongodb.mongodb id
```