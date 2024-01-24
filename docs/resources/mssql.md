---
subcategory: "Mssql"
---


# Resource: ncloud_mssql

Provides a MSSQL instance resource.

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
  subnet_type    = "PUBLIC"
  usage_type     = "GEN"
}

resource "ncloud_mssql" "example-mssql" {
  subnet_no           = ncloud_subnet.example-sub.id 
  service_name        = "my-tf-mssql"
  is_ha               = true
  is_automatic_backup = true
  user_name           = "test"
  user_password       = "password1!"
}
```

## Argument Reference

The following arguments are supported:

* `subnet_no` - (Required) The ID of the associated Subnet.
* `service_name` - (Required) Service name to create. Only English alphabets, numbers, dash ( - ) and Korean letters can be entered. Min: 3, Max: 15
* `is_ha` - (Required) Whether is High Availability or not. If High Availability is selected, 2 servers including the Standby Master server will be created and additional charges will be incurred. Default : true.
* `user_name` - (Required) MSSQL access User ID. - Only English letters, numbers, and underscore characters ( _ ) are allowed, and must start with an English letter. Min: 4, Max: 16
* `user_password` - (Required) MSSQL access  User Password. Must be at least 8 characters in length and contain at least 1 each of English letter, special character, and number. The following characters cannot be used in the password: ` & \ " ' / and space. Min: 8, Max: 20
* `config_group_no` - (Optional) MSSQL config group Number. Already-created Config Group can be applied when creating a server. When you do not have any config groups, you can select from provided config groups by default. You can view through getCloudMssqlConfigGroupList API. Default: 0
* `image_product_code` - (Optional) Image product code to determine the MSSQL instance server image specification to create. If not entered, the instance is created for default value. It can be obtained through [`ncloud_mssql_image_products` data source](../data-sources/mssql_image_products.md)
* `product_code` - (Optional) Product code to determine the MSSQL instance server image specification to create. It can be obtained through [`ncloud_mssql_products` data source](../data-sources/mssql_products.md). Default : Minimum specifications(1 memory, 2 cpu)
* `data_storage_type` - (Optional) Data storage type. You can select `SSD|HDD`. Default: SSD
* `backup_file_retention_period` - (Optional) Backups are performed daily and backup files are stored in separate backup storage. Fees are charged based on the space used. Default : 1(1 day), Min: 1, Max: 30
* `backup_time` - (Optional, Required if `is_backup` is true and `is_automatic_backup` is false) You can set the time when backup is performed. it must be entered if backup status(is_backup) is true and automatic backup status(is_automatic_backup) is false.
* `is_automatic_backup` - (Optional) You can select whether to automatically set the backup time. if `is_automatic_backup` is true, backup_time cannot be entered. Default : true 
* `port` - (Optional) You can set TCP port to access the mssql instance. Default : 1433, Min: 10000, Max: 20000
* `character_set_name` - (Optional) DB character set can be selected from Korean and English collation. You can view through getCloudMssqlCharacterSetList API. Default: Korean_Wansung_CI_AS. Options: `Korean_Wansung_CI_AS`, `SQL_Latin1_General_CP1_CI_AS`

## Attributes Reference

In addition to all arguments above, the following attributes are exported

* `id` - MSSQL instance number.
* `engine_version` - MSSQL Engine version.
* `region_code` - Region code.
* `vpc_no` - The ID of the associated Vpc.
* `access_control_group_no_list` - The ID list of the associated Access Control Group.
* `mssql_server_list` - The list of the MySQL server.
  * `server_instance_no` - Server instance number.
  * `server_name` - Server name.
  * `server_role` - Member or Arbiter or Mongos or Config.
  * `zone_code` - Zone code.
  * `subnet_no` - The ID of the associated Subnet.
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

* MSSQL can be imported using the `id`. For example:

```console
$ terraform import ncloud_mssql.rsc_name 12345
```

### `import` block

* In Terraform v1.7.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import MSSQL using the `id`. For example:

```terraform
import {
  to = ncloud_mssql.rsc_name
  id = "12345"
}
```
