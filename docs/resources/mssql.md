---
subcategory: "Mssql"
---


# Resource: ncloud_mssql

Provides a MSSQL instance resource.

## Example Usage

#### Basic (VPC)

```hcl
resource "ncloud_vpc" "test_vpc" {
    name               = "%[1]s"
	ipv4_cidr_block    = "10.0.0.0/16"
}
resource "ncloud_subnet" "test_subnet" {
	vpc_no             = ncloud_vpc.test_vpc.vpc_no
	name               = "%[1]s"
	subnet             = "10.0.0.0/24"
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.test_vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
}
resource "ncloud_mssql" "mssql" {
	vpc_no = ncloud_vpc.test_vpc.vpc_no
	subnet_no = ncloud_subnet.test_subnet.id
	service_name = "%[1]s"
	is_ha = true
	is_multi_zone = false
	is_automatic_backup = true
	user_name = "test"
	user_password = "qwer1234!"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_no` - (Required) The ID of the associated Vpc.
* `subnet_no` - (Required) The ID of the associated Subnet.
* `service_name` - (Required) Service name to create.
* `user_name` - (Required) MSSQL User ID
* `user_password` - (Required) MSSQL User Password
* `is_ha` - (Optional) Whether is High Availability or not. If High Availability is selected, 2 servers including the Standby Master server will be created and additional charges will be incurred. Default : true.
* `is_multi_zone` - (Optional) When High Availability is set to true, You can determine Whether is Multi Zone or not. Default : false
* `mirror_subnet_no` - (Optional) The ID of the Mirror Server Subnet. It will be required when `is_multi_zone` is true. Default : false
* `port` - (Optional) You can set TCP port to access the mssql instance. Default : 1433
* `data_storage_type_code` - (Optional) Data storage type. You can select `SSD|HDD`.
* `backup_file_retention_period` - (Optional) Backups are performed daily and backup files are stored in separate backup storage. Fees are charged based on the space used. Default : 1(1 day)
* `backup_time` - (Optional, Required if `is_backup` is true and `is_automatic_backup` is false) You can set the time when backup is performed. it must be entered if backup status(is_backup) is true and automatic backup status(is_automatic_backup) is false.
* `image_product_code` - (Optional) Image product code to determine the mssql instance server image specification to create. You can get it through the `getCloudMssqlImageProductList` action. If not entered, the instance is created for default value.
* `product_code` - (Optional) Product code to determine the mssql instance server image specification to create. You can get it through the `getCloudMssqlProductList` action. Default : Minimum specifications(1 memory, 2 cpu)
* `is_automatic_backup` - (Optional) You can select whether to automatically set the backup time. if `is_automatic_backup` is true, backup_time cannot be entered. Default : true 

## Attributes Reference

* `id` - MSSQL Instance No.

## Import

### `terraform import` command

* MSSQL can be imported using the `id`. For example:

```console
$ terraform import ncloud_mssql.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import MSSQL using the `id`. For example:

```terraform
import {
  to = ncloud_mssql.rsc_name
  id = "12345"
}
```
