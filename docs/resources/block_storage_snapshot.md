---
subcategory: "Server"
---


# Resource: ncloud_block_storage

Provides a ncloud block storage snapshot resource.

~> **NOTE:** This resource only supports Classic environment.

## Example Usage

```hcl
resource "ncloud_block_storage_snapshot" "snapshot" {
	block_storage_instance_no = "812345"
	name = "tf-test-snapshot1"
	description = "Terraform test snapshot1"
}
```

## Argument Reference

The following arguments are supported:

* `block_storage_instance_no` - (Required) Block storage instance Number for creating snapshot.
* `name` - (Optional) Block storage snapshot name to create. default : Ncloud assigns default values.
* `description` - (Optional) Descriptions on a snapshot to create.

## Attributes Reference

* `instance_no` - Block Storage Snapshot Instance Number
* `volume_size` - Block Storage Snapshot Volume Size
* `original_block_storage_instance_no` - Original Block Storage Instance Number
* `original_block_storage_name` - Original Block Storage Name
* `instance_status` - Block Storage Snapshot Instance Status code
* `instance_status_name` - Block Storage Snapshot Instance Status Name
* `instance_operation` - Block Storage Snapshot Instance Operation code
* `create_date` - Creation date of the block storage snapshot instance
* `server_image_product_code` - Server Image Product Code
* `os_information` - OS Information
