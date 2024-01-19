---
subcategory: "MySQL"
---


# Resource: ncloud_mysql

Provides a MySQL instance resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage

```terraform
resource "ncloud_vpc" "test" {
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "test" {
  vpc_no         = ncloud_vpc.test.vpc_no
  subnet         = cidrsubnet(ncloud_vpc.test.ipv4_cidr_block, 8, 1)
  zone           = "KR-2"
  network_acl_no = ncloud_vpc.test.default_network_acl_no
  subnet_type    = "PUBLIC"
}

resource "ncloud_mysql" "mysql" {
  subnet_no = ncloud_subnet.test.id
  service_name = "my-tf-mysql"
  server_name_prefix = "name-prefix"
  user_name = "username"
  user_password = "password1!"
  host_ip = "192.168.0.1"
  database_name = "db_name"
}
```


## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name to create. Only English alphabets, numbers, dash ( - ) and Korean letters can be entered. Min: 3, Max: 20
* `server_name_prefix` - (Required) Server name prefix to create. In order to prevent overlapping host names, random text is added. Can comprise only lower-case English alphabets, numbers and dash ( - ). The first letter must be an English alphabet and the last letter must be an English alphabet or a number. Min: 3, Max: 30
* `user_name` - (Required) MySQL User ID. Only English alphabets, numbers and special characters ( \ _ , - ) are allowed and must start with an English alphabet. Min: 4, Max: 16
* `user_password` - (Required) MySQL User Password. At least one English alphabet, number and special character must be included. Certain special characters ( ` & + \ " ' / space ) cannot be used. Min: 8, Max: 20
* `host_ip` - (Required) MySQL user host. ex) Overall connection permitted: %, Connection by specific IPs permitted: 1.1.1.1, IP band connection permitted: 1.1.1.%
* `database_name` - (Required) Database name to create. Only English alphabets, numbers and special characters ( \ _ , - ) are allowed and must start with an English alphabet. Min: 1, Max: 30
* `subnet_no` - (Required) The ID of the associated Subnet. Public domain can only be used on a DB server generated on Public Subnet. 
* `image_product_code` - (Optional) Image product code to determine the MySQL instance server image specification to create. If not entered, the instance is created for default value. It can be obtained through [`ncloud_mysql_image_products` data source](../data-sources/mysql_image_products.md)
* `product_code` - (Optional) Product code to determine the MySQL instance server image specification to create. It can be obtained through [`ncloud_mysql_products` data source](../data-sources/mysql_products.md). Default : Minimum specifications(1 memory, 2 cpu)
* `data_storage_type` - (Optional) Data storage type. If `generationCode` is `G2`, You can select `SSD|HDD`, else if `generationCode` is `G3`, you can select CB1. Default : SSD in G2, CB1 in G3
* `is_ha` - (Optional) High-availability (True/False).  If high availability is selected, 2 servers including a Standby Master server are generated, and additional fees are incurred. If the high availability status `is_ha` is false, `is_multi_zone` and `standby_master_subnet_no` parameters are not used. Default : true.
* `is_multi_zone` - (Optional) Multi-zone (True/False). If the high availability status `is_ha` is true, multi-zone can be selected. If multi-zone is selected, the Master server and Standby Master server are generated in mutually different zones, providing higher availability. Default : false
* `is_storage_encryption` - (Optional) Whether data storage encryption is applied. If encryption is applied, DB data is encrypted and stored in the storage after Cloud DB for MySQL instance is generated, storage encryption setting cannot be changed. Encryption can be applied only if the high availability status `is_ha` is true. Default : false
* `is_backup` - (Optional) When High Availability is set to true, it is fixed true. You can determine whether to back up. Default : true
* `backup_file_retention_period` - (Optional) Backups are performed daily and backup files are stored in separate backup storage. Fees are charged based on the space used. Default : 1(1 day), Min: 1, Max: 30
* `backup_time` - (Optional, Required if `is_backup` is true and `is_automatic_backup` is false) You can set the time when backup is performed. it must be entered if backup status(is_backup) is true and automatic backup status(is_automatic_backup) is false.
* `is_automatic_backup` - (Optional) You can select whether to automatically set the backup time. if `is_automatic_backup` is true, backup_time cannot be entered. Default : true 
* `port` - (Optional) You can set TCP port to access the MySQL instance. Default : 3306, Min: 10000, Max: 20000
* `standby_master_subnet_no` - (Optional, Required if `is_multi_zone` is true) if `is_multi_zone` is false, input is not accepted. if `is_multi_zone` is true, input must be entered. `standby_master_subnet_no` must be different from the master server's subnet and zone. And must be the same Public or Private. You can get it through the `getCloudMysqlTargetSubnetList` action.

## Attributes Reference

In addition to all arguments above, the following attributes are exported

* `id` - MySQL Instance Number.
* `engine_version_code` - MySQL Engine version code.
* `region_code` - Region code.
* `zone_code` - Zone code.
* `vpc_no` - The ID of the associated Vpc.
* `access_control_group_no_list` - The ID list of the associated Access Control Group.
* `mysql_config_list` - The list of config.
* `mysql_server_list` - The list of the MySQL server.
  * `server_instance_no` - Server instance number.
  * `server_name` - Server name.
  * `server_role` - Member or Arbiter or Mongos or Config.
  * `product_code` - Product code.
  * `is_public_subnet` - Public subnet status.
  * `private_domain` - Private domain.
  * `public_domain` - Public domain.
  * `memory_size` - Available memory size.
  * `cpu_count` - CPU count.
  * `data_storage_size` - Storage size.
  * `used_data_storage_size` - Size of data storage in use.
  * `uptime` - Running start time.
  * `create_date` - Server create date.

## Import

### `terraform import` command

* MySQL can be imported using the `id`. For example:

```console
$ terraform import ncloud_mysql.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import MySQL using the `id`. For example:

```terraform
import {
  to = ncloud_mysql.rsc_name
  id = "12345"
}
```
