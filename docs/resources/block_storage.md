---
subcategory: "Server"
---


# Resource: ncloud_block_storage

Provides a Block Storage resource.

## Example Usage

```terraform
resource "ncloud_block_storage" "storage" {
	server_instance_no = "812345"
	name = "tf-test-storage1"
	size = "10"
}
```

## Argument Reference

The following arguments are supported:

* `size` - (Required) The size of the block storage to create. Automatically determined if created using XEN type block storage snapshots. If created using a KVM type block storage snapshot, must be greater than or equal to the snapshot size. Enter in 10 GB increments. XEN type Min: 10GB, Max: 2000 GB. KVM type Min: 10GB, Max : 16380 GB.
* `server_instance_no` - **(Required) When first created**. (Optional) When changing the value after creation. Server instance ID to which you want to assign the block storage.
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name. Min: 3, Max: 30. Only English letters, numbers, and the special character "-" can be used. It must start with an English letter. It must end with an English letter or number.
* `description` - (Optional) description to create. Min: 0, Max: 1000 Bytes.
* `disk_detail_type` - (Optional) Type of block storage disk detail to create. Valid for XEN type only. Conflicts with `volume_type`. Default `SSD`. Accepted values: `SSD` | `HDD` 
* `stop_instance_before_detaching` - (Optional, Boolean) Set this to true to ensure that the target instance is stopped before trying to detach the block storage. It stops the instance, if it is not already stopped.
	> If `stop_instance_before_detaching` is `true`, server will be stopped and **will not start automatically**. User must start server instance manually via NCLOUD console or API.

~> **NOTE:** Below arguments only support VPC environment.

* `zone` - (Optional, Required if to select KVM type) The availability zone in which the block storage instance will be created. It must be the same zone code as the server..
* `snapshot_no` - (Optional) Create the block storage from the snapshots you take.
* `hypervisor_type` - (Optional) Hypervisor type. Requied with `volume_type`. (`XEN` or `KVM`)
* `volume_type` - (Optional) Decides the volume type of the block storage to be created. Required for KVM block storage. Conflicts with `disk_detail_type`. Required with `hypervisor_type`. Options : `XEN` type(` SSD` | `HDD`), `KVM`type(`FB1` | `CB1`)
* `is_return_protection` - (Optional) Enable return protection. Default: `false`. Options: `true`| `false`

## Attributes Reference

* `id` - The ID of Block storage instance.
* `block_storage_no` - The ID of Block storage instance. (It is the same result as `id`)
* `server_name` - Server name.
* `type` - Block storage type code.
* `device_name` - Device name.
* `product_code` - Block storage product code.
* `status` - Block storage instance status code.
* `disk_type` - Disk type code.
* `max_iops` - Maximum IOPS.
* `is_encrypted_volume` - Volume encryption status. (`true` or `false`)

## Import

### `terraform import` command

* Block Storage can be imported using the `id`. For example:

```console
$ terraform import ncloud_block_storage.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Block Storage using the `id`. For example:

```terraform
import {
  to = ncloud_block_storage.rsc_name
  id = "12345"
}
```
