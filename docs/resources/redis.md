---
subcategory: "Redis"
---


# Resource: ncloud_redis

Provides a Redis resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage

```terraform
resource "ncloud_vpc" "example-v" {
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "example-sub" {
  vpc_no         = ncloud_vpc.example-v.vpc_no
  subnet         = cidrsubnet(ncloud_vpc.example-v.ipv4_cidr_block, 8, 1)
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.example-v.default_network_acl_no
  subnet_type    = "PRIVATE"
  usage_type     = "GEN"
}

resource "ncloud_redis_config_group" "example-rcg" {
  name          = "test-rcg"
  redis_version = "7.0.13-simple"
  description   = "example"
}

resource "ncloud_redis" "example-redis" {
  service_name       = "my-tf-redis"
  server_name_prefix = "tf-svr"
  vpc_no             = ncloud_vpc.example-v.vpc_no
  subnet_no          = ncloud_subnet.example-sub.id 
  config_group_no    = ncloud_redis_config_group.example-rcg.config_group_no
  mode               = "SIMPLE"
}
```

## Argument Reference
The following arguments are supported:

* `service_name` - (Required) Service name to create. Enter the group name of the Redis server (e.g., NAVER-HOME). You cannot double-use the Redis service name. Only alphanumeric characters, numbers, hyphens (-), and Korean characters are allowed. Min: 3, Max: 15
* `server_name_prefix` - (Required) Enter the name prefix of the Redis Server. The Redis server name is created with a 3-digit number, which is automatically created. You cannot double-use the Redis Server name. It must only contain English letters (lowercase), numbers, and hyphens (-). It must start with an English letter and end with an English letter or a number. Min: 3, Max: 15
* `user_name` - (Optional, Required if `gov` site) Redis User ID. Available only `gov` site. Only English alphabets, numbers and special characters ( \ _ , - ) are allowed and must start with an English alphabet. Min: 4, Max: 16
* `user_password` - (Optional, Required if `gov` site) Redis User Password. Available only `gov` site. At least one English alphabet, number and special character must be included. Certain special characters ( ` & + \ " ' / space ) cannot be used. Min: 8, Max: 20
* `vpc_no` - (Required) VPC number. Determining the VPC in which the Cloud DB for Redis instance will be created.
* `subnet_no` - (Required) The ID of the associated Subnet. Subnet transfer is not possible after a Cloud DB for Redis instance has been created.
* `config_group_no` - (Required) Redis Config Group number. Config groups are provided, and one cluster group uses the same config. A new config group must be created if none exists. It can be changed online after creation.
* `mode` - (Required) Determines the configuration of Cloud DB for Redis. When the CLUSTER setting is used, the `is_ha` setting is ignored. Options: `CLUSTER`, `SIMPLE`.
* `image_product_code` - (Optional) Image product code to determine the Redis instance server image specification to create. If not entered, the instance is created for default value. It can be obtained through [`data.ncloud_redis_image_products`](../data-sources/redis_image_products.md).
* `product_code` - (Optional) Sets the server specifications of the Cloud DB for Redis instance to be created. It can be obtained through [`data.ncloud_redis_products`](../data-sources/redis_products.md). Default : Minimum specifications(1 memory, 2 cpu)
* `shard_count` - (Optional) Number of shards to be created.  3 to 10 Number of master nodes. Necessary only if the `mode` is CLUSTER. Default: 3
* `shard_copy_count` - (Optional) Replicas per shard Redis Cluster consists of the master node and slave node. A slave node is necessary for HA. When adding a replica, one slave node is assigned to each master node.  For example, 3 shards, 1 replica per shard -> Master node: 3, Slave node: 3. You can enter 0 to 4 replica(s) for each shard. If the number of replicas per shard is set to 0, then high availability can't be supported. Necessary only if the `mode` is CLUSTER. Default: 0
* `is_ha` - (Optional) Whether is High Availability or not. The Cloud DB for Redis product supports automatic failure recovery using the Standby master. When high availability is supported, additional charges are incurred and backup is automatically configured. Default: false
* `is_backup` - (Optional) Backup status. If the high availability status `is_ha` is true, then the backup setting status is fixed as true. Default : false
* `backup_file_retention_period` - (Optional) Backups are performed on a daily basis, and backup files are stored in a separate backup storage. Charges are based on the storage space used. Default: 1 (1 day)
* `is_automatic_backup` - (Optional) Select whether to have backup times set automatically. When the automatic backup is true, then any `backup_time` entered is ignored and the backup time is configured automatically.
* `backup_time` - (Optional, Required if `is_backup` is true and `is_automatic_backup` is false) You can set the time when backup is performed. it must be entered if backup status(is_backup) is true and automatic backup status(is_automatic_backup) is false. EX) 01:15
* `port` - (Optional) Cloud Redis port. You need to enter the TCP port number of Redis access. Value range:	6379 or Min: 10000, Max: 20000. Default: 6379

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `id` - Redis instance number.
* `backup_schedule` - Automatic or User-defined.
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

## Import

### `terraform import` command

* Redis can be imported using the `id`. For example:

```console
$ terraform import ncloud_redis.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Redis using the `id`. For example:

```terraform
import {
  to = ncloud_redis.rsc_name
  id = "12345"
}
```
