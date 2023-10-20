---
subcategory: "MySQL"
---


# Resource: ncloud_mysql

Provides a MySQL instance resource.

## Example Usage

#### Basic (VPC)

```hcl
resource "ncloud_login_key" "loginkey" {
  key_name = "test-key"
}

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
  name_prefix = "name-prefix"
  user_name = "username"
  user_password = "password1!"
  host_ip = "192.168.0.1"
  database_name = "db_name"
}
```

## Argument Reference

The following arguments are supported:

* `image_product_code` - (Optional) Image product code to determine the MySQL instance server image specification to create. If not entered, the instance is created for default value. It can be obtained through `data.ncloud_mysql_image_product`.
    - [Docs mysql(5.7) Centos 7.3-64 Image Products](https://github.com/NaverCloudPlatform/terraform-ncloud-docs/blob/main/docs/vpc_products/mysql(5.7)-centos-7.3-64.md)
    - [Docs mysql(5.7) ubuntu 16.04-64 Image Products](https://github.com/NaverCloudPlatform/terraform-ncloud-docs/blob/main/docs/vpc_products/mysql(5.7)-ubuntu-16.04-64-server.md)
    - [Docs mysql(8.0) centos 7.3-64 Image Products](https://github.com/NaverCloudPlatform/terraform-ncloud-docs/blob/main/docs/vpc_products/mysql(5.7)-centos-7.3-64.md)
    - [Docs mysql(8.0) centos 7.8-64 Image Products](https://github.com/NaverCloudPlatform/terraform-ncloud-docs/blob/main/docs/vpc_products/mysql(8.0)-centos-7.8-64.md)
    - [`ncloud_mysql_image_product` data source](../data-sources/mysql_image_product.md)


* `product_code` - (Optional) Product code to determine the MySQL instance server image specification to create. It can be obtained through `data.ncloud_mysql_product`. Default : Minimum specifications(1 memory, 2 cpu)
    - [Docs mysql(5.7) Centos 7.3-64 Image Products](https://github.com/NaverCloudPlatform/terraform-ncloud-docs/blob/main/docs/vpc_products/mysql(5.7)-centos-7.3-64.md)
    - [Docs mysql(5.7) ubuntu 16.04-64 Image Products](https://github.com/NaverCloudPlatform/terraform-ncloud-docs/blob/main/docs/vpc_products/mysql(5.7)-ubuntu-16.04-64-server.md)
    - [Docs mysql(8.0) centos 7.3-64 Image Products](https://github.com/NaverCloudPlatform/terraform-ncloud-docs/blob/main/docs/vpc_products/mysql(5.7)-centos-7.3-64.md)
    - [Docs mysql(8.0) centos 7.8-64 Image Products](https://github.com/NaverCloudPlatform/terraform-ncloud-docs/blob/main/docs/vpc_products/mysql(8.0)-centos-7.8-64.md)
    - [`ncloud_mysql_product` data source](../data-sources/mysql_product.md)

* `service_name` - (Required) Service name to create.
* `name_prefix` - (Required) Server name prefix to create.
* `user_name` - (Required) MySQL User ID
* `user_password` - (Required) MySQL User Password
* `host_ip` - (Required) MySQL host Ip
* `database_name` - (Required) Database name to create. 
* `subnet_no` - (Required) The ID of the associated Subnet.
* `engine_version_code` - (Optional) Engine version code to determine the MySQL engine specification to create. It can be obtained through the `data.(s)` action. Default : Selected as the latest version.
* `data_storage_type_code` - (Optional) Data storage type. If `generationCode` is `G2`, You can select `SSD|HDD`, else if `generationCode` is `G3`, you can select CB1. Default : SSD in G2, CB1 in G3
* `is_ha` - (Optional) Whether is High Availability or not. If High Availability is selected, 2 servers including the Standby Master server will be created and additional charges will be incurred. Default : true.
* `is_multi_zone` - (Optional) When High Availability is set to true, You can determine Whether is Multi Zone or not. Default : false
* `is_storage_encryption` - (Optional) When High Availability is set to true, You can determine whether to encrypt the database storage or not. Default : false
* `is_backup` - (Optional) When High Availability is set to true, it is fixed true. You can determine whether to back up. Default : true
* `backup_file_retention_period` - (Optional) Backups are performed daily and backup files are stored in separate backup storage. Fees are charged based on the space used. Default : 1(1 day)
* `backup_time` - (Optional, Required if `is_backup` is true and `is_automatic_backup` is false) You can set the time when backup is performed. it must be entered if backup status(is_backup) is true and automatic backup status(is_automatic_backup) is false.
* `is_automatic_backup` - (Optional) You can select whether to automatically set the backup time. if `is_automatic_backup` is true, backup_time cannot be entered. Default : true 
* `port` - (Optional) You can set TCP port to access the MySQL instance. Default : 3306
* `standby_master_subnet_no` - (Optional, Required if `is_multi_zone` is true) if `is_multi_zone` is false, input is not accepted. if `is_multi_zone` is true, input must be entered. `standby_master_subnet_no` must be different from the master server's subnet and zone. And must be the same Public or Private. You can get it through the `getCloudMysqlTargetSubnetList` action.

## Attributes Reference

* `instance_no` - The ID of serve****r instance.
* `vpc_no` - The ID of the associated Vpc.
* `mysql_server_list` - The list of the MySQL server.

## Import

MySQL instance can be imported using id, e.g.,

``` 
$ terraform import ncloud_mysql.my_mysql id
```