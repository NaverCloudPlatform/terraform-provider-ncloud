---
subcategory: "MongoDB"
---


# Resource: ncloud_mongodb

Provides a Database Service MongoDB resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage

```terraform
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
  vpc_no = ncloud_vpc.vpc.id
  subnet_no = ncloud_subnet.subnet.id
  service_name = "sample-mongodb"
  server_name_prefix = "tf-svr"
  user_name = "username"
  user_password = "password1!"
  cluster_type_code = "STAND_ALONE"
}
```


## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name to create. Enter group name of DB server. Specify the replica set name with the entered DB service name. Only alphanumeric characters, numbers, hyphens (-), and Korean characters are allowed. Duplicate names and changes after creation are prohibited. Min: 3, Max: 15
* `server_name_prefix` - (Required) Enter the name prefix of the MongoDb Server. It is created with random text added after the transferred cloudMongoDbServerNamePrefix value to avoid duplicated host names. It must only contain English letters (lowercase), numbers, and hyphens (-). It must start with an English letter and end with an English letter or a number. Min: 3, Max: 15
* `user_name` - (Required) Username for access. Must assign username in the role of DB admin. Only English letters, numbers, underscores (_), and hyphens (-) are allowed and it must start with an English letter. Min: 4, Max: 16
* `user_password` - (Required) Password for access. Must assign password of the username in the role of DB admin. It must have at least 1 English letter, 1 number, and 1 special character. The following characters cannot be used in the password: ` & + \ " ' / space. Min: 8, Max: 20
* `vpc_no` - (Required) The ID of the associated Vpc.
* `subnet_no` - (Required) The ID of the associated Subnet.
* `cluster_type_code` - (Required) MongoDB cluster type code determines the cluster type of MongoDB. Options: STAND_ALONE | SINGLE_REPLICA_SET | SHARDED_CLUSTER
* `image_product_code` - (Optional) MongoDB image product code. If not entered, it is created as a default value. It can be obtained through [`data.ncloud_mongodb_image_products`](../data-sources/mongodb_image_products.md).
* `engine_version_code` - (Optional) MongoDB engine version code. Only entered when `generation_code` is G3, If not entered, generate with the latest version currently available.
* `member_product_code` - (Optional) Member server product code. It can be obtained through [`data.ncloud_mongodb_products`](../data-sources/mongodb_products.md). Default: select the minimum specifications and must be based on 1. Memory and 2. CPU
* `arbiter_product_code` - (Optional) Arbiter server product code. It can be obtained through [`data.ncloud_mongodb_products`](../data-sources/mongodb_products.md). Default: select the minimum specifications and must be based on 1. Memory and 2. CPU
* `mongos_product_code` - (Optional) Mongos server product code. It can be obtained through [`data.ncloud_mongodb_products`](../data-sources/mongodb_products.md). Default: select the minimum specifications and must be based on 1. Memory and 2. CPU
* `config_product_code` - (Optional) Config server product code. It can be obtained through [`data.ncloud_mongodb_products`](../data-sources/mongodb_products.md). Default: select the minimum specifications and must be based on 1. Memory and 2. CPU
* `shard_count` - (Optional, Changeable) The number of MongoDB Shards. The number of shards can be defined for sharding. Only 2 or 3 are allowed for the initial configuration. Only enter when `cluster_type_code` is SHARDED_CLUSTER. Default: 2, Min: 2, Max: 5 
* `member_server_count` - (Optional, Changeable) The number of MongoDB Member Servers. The number of member servers per replica set (or per shard if sharding) can be defined. Selectable between 3 to 7, including arbiter servers. Default : 3, Min: 2, Max: 7
* `arbiter_server_count` - (Optional, Changeable) The number of MongoDB Arbiter servers. You can select whether to use the Arbiter server per Replica Set (for each shard in the case of Sharding). Up to one Arbiter server can be selected. The Arbiter server is provided with a minimum configurable spec. Default: 0, Min: 0, Max: 1
* `mongos_server_count` - (Optional, Changeable) The number of MongoDB Mongos servers. If sharding is used, the number of mongos servers can be selected. Default: 2, Min: 2, Max: 5
* `config_server_count` - (Optional, Changeable) The number of MongoDB Config servers. If sharding is used, the config server's logarithm can be selected. Only 3 are allowed for the initial configuration. Default: 3, Min: 3, Max: 7 
* `backup_file_retention_period` - (Optional) Backups are performed daily and backup files are stored in separate backup storage. Fees are charged based on the space used. Default: 1(1 day), Min: 1, Max: 30
* `backup_time` - (Optional) You can set the time when backup is performed. Default: 02:00. HHMM format. You must enter in 15-minute increments.
* `data_storage_type` - (Optional) Data storage type. If `generationCode` is `G2`, You can select `SSD|HDD`, else if `generationCode` is `G3`, you can select CB1. Default : SSD in G2, CB1 in G3
* `member_port` - (Optional) TCP port number for access to the MongoDB Member Server. Default: 17017, Min: 10000, Max: 65535
* `mongos_port` - (Optional) TCP port number for access to the MongoDB Mongos Server.  Default: 17017, Min: 10000, Max: 65535
* `config_port` - (Optional) TCP port number for access to the MongoDB Config Server.  Default: 17017, Min: 10000, Max: 65535
* `compress_code` - (Optional) MongoDB Data Compression Algorithm Code allows you to select data compression algorithms provided by MongoDB. Default: SNPP,  Options: SNPP | ZLIB | ZSTD | NONE

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - MondoDb instance number. 
* `arbiter_port` - TCP port number for access to the MongoDB Arbiter Server.
* `region_code` - Region code.
* `zone_code` - Zone code.
* `access_control_group_no_list` - The ID list of the associated Access Control Group.
* `mongodb_server_list` - The list of the MongoDB server.
  * `server_instance_no` - Server instance number.
  * `server_name` - Server name.
  * `server_role` - Member or Arbiter or Mongos or Config.
  * `cluster_role` - STAND_ALONE or SINGLE_REPLICA_SET or SHARD or CONFIG or MONGOS.
  * `product_code` - Product code.
  * `private_domain` - Private domain.
  * `public_domain` - Public domain.
  * `replica_set_name` - Replica set name.
  * `memory_size` - Available memory size.
  * `cpu_count` - CPU count.
  * `data_storage_size` - Storage size.
  * `uptime` - Running start time.
  * `create_date` - Server create date.

## Import

### `terraform import` command

* MongoDB can be imported using the `id`. For example:

```console
$ terraform import ncloud_mongodb.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import MongoDB using the `id`. For example:

```terraform
import {
  to = ncloud_mongodb.rsc_name
  id = "12345"
}
```
