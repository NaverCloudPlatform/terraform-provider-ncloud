---
subcategory: "PostgreSQL"
---

# Resource: ncloud_postgresql

Provides a PostgreSQL instance resource.

~> **NOTE** This resource only supports VPC environment.

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
resource "ncloud_postgresql" "postgresql" {
  subnet_no = ncloud_subnet.test.id
  service_name = "tf-postgresql"
  server_name_prefix = "name-prefix"
  user_name = "username"
  user_password = "password1!"
  client_cidr = "0.0.0.0/0"
  database_name = "db_name"
}
```

## Argument Reference

The following arguments are supported:

* `service_name` - (Required) Service name to create. Only alphanumeric characters, numbers, hyphens (-), and Korean characters are allowed. Min: 3, Max: 30
* `server_name_prefix` - (Required) Server name prefix to create. In order to prevent overlapping host names, random text is added. It must only contain English letters (lowercase), numbers, and hyphens (-). It must start with an English letter and end with an English letter or a number. Min: 3, Max: 20
* `database_name` - (Required) Database name to create. Only English alphabets, numbers and special characters ( \ _ , - ) are allowed and must start with an English alphabet. Min: 1, Max: 30
* `user_name` - (Required) PostgreSQL User ID. Only English alphabets, numbers and special characters ( \ _ , - ) are allowed and must start with an English alphabet. Cannot include User ID. Min: 4, Max: 16
* `user_password` - (Required) PostgreSQL User Password. At least one English alphabet, number and special character must be included. Certain special characters ( ` & + \ " ' / space ) cannot be used. Min: 8, Max: 20
* `client_cidr` - (Required) Access Control (CIDR) of the client you want to connect to EX) Allow all access: 0.0.0.0/0, Allow specific IP access: 192.168.1.1/32, Allow IP band access: 192.168.1.0/24
* `image_product_code` - (Optional) Image product code to determine the PostgreSQL instance server image specification to create. If not entered, the instance is created for default value. It can be obtained through [`ncloud_postgresql_image_products` data source](../data-sources/postgresql_image_products.md)
* `product_code` - (Optional) Product code to determine the PostgreSQL instance server image specification to create. It can be obtained through [`ncloud_postgresql_products` data source](../data-sources/postgresql_products.md). Default: Minimum specifications(1 memory, 2 cpu)
* `data_storage_type_code` - (Optional) Data storage type. You can select `SSD|HDD`. Default: SSD
* `engine_version_code` - PostgreSQL Engine version code.
* `is_ha` - (Optional) High-availability (true/false). If high availability is selected, 2 servers including a Secondary server are generated, and additional fees are incurred. If the high availability status `is_ha` is false, `is_multi_zone` and `secondary_subnet_no` parameters are not used. Default: true.
* `is_multi_zone` - (Optional) Multi-zone (true/false). If the high availability status `is_ha` is true, multi-zone can be selected. If multi-zone is selected, Primary server and Secondary server are generated in mutually different zones, providing higher availability. Default: false 
* `is_storage_encryption` - (Optional) Whether data storage encryption is applied. If encryption is applied, DB data is encrypted and stored in storage. After Cloud DB for PostgreSQL instance is generated, storage encryption setting cannot be changed. Not available in Neurocloud environment.
* `is_backup` - (Optional) Backup status. If the high availability status `is_ha` is true then the backup setting status is fixed as true. Default: false
* `backup_time` - (Optional, Required if `is_backup` is true and `is_automatic_backup` is false) You can set the time when backup is performed. It must be entered if backup status(is_backup) is true and automatic backup status(is_automatic_backup) is false. EX) 01:15 
* `backup_file_retention_period` - (Optional) Backups are performed on a daily basis, and backup files are stored in a separate backup storage. Charges are based on the storage space used. Default: 1 (1 day)
* `backup_file_storage_count` - (Optional) Number of backup files kept. Min: 1, Max: 30
* `is_backup_file_compression` - (Optional) Whether to compress backup files. Default: true
* `is_automatic_backup` - (Optional) Select wheter to have backup times set automatically. If `is_automatic_backup` is true, `backup_time` cannot be entered. Default: true
* `port` - (Optional) TCP port to access the Cloud DB for PostgreSQL instance. Default: 5432, Min: 10000, Max: 20000
* `secondary_subnet_no` - (Optional, Required if `is_multi_zone` is true) `secondary_subnet_no` must be different from the Primary server's subnet and zone. And must be the same Public or Private. You can get it through the `getCloudPostgresqlTargetSubnetList` action. 

## Attributes Reference

In addition to all arguments above, the following attributes are exported

* `id` - PostgreSQL Instance Number. 
* `region_code` - Region code.
* `vpc_no` - The ID of the associated Vpc.
* `access_control_group_no_list` - The ID list of the associated Access Control Group.
* `postgresql_config_list` - The list of config.
* `postgresql_server_list` - The list of the PostgreSQL server.
  * `server_instance_no` - Server instance number.
  * `server_name` - Server name.
  * `server_role` - Server role code. M(Primary), H(Secondary), S(Read Replica)
  * `zone_code` - Zone code. 
  * `product_code` - Product code.
  * `is_public_subnet` - Public subnet status.
  * `public_domain` - Public domain.
  * `private_domain` - Private domain.
  * `private_ip` - Private IP.
  * `memory_size` - Available memory size.
  * `cpu_count` - CPU count.
  * `data_storage_size` - Storage size.
  * `used_data_storage_size` - Size of data storage in use.
  * `uptime` - Running start time.
  * `create_date` - Server create date.

## Import

### `terraform import` command

* PostgreSQL can be imported using the `id`. For example:

```console
$ terraform import ncloud_postgresql.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import PostgreSQL using the `id`. For example:

```terraform
import {
    to = ncloud_postgresql.rsc_name
    id = "12345"
}
```