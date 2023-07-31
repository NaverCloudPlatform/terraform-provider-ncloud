---
subcategory: "Member Server Image"
---


# Data Source: ncloud_member_server_image

Gets a member server image.

## Example Usage

```hcl
data "ncloud_member_server_image" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the member server image list returned by ncloud
* `no_list` - (Optional) List of member server images to view
* `platform_type_code_list` - (Optional) List of platform codes of server images to view. Linux 32Bit (`LNX32`) | Linux 64Bit (`LNX64`) | Windows 32Bit (`WND32`) | Windows 64Bit (`WND64`) | Ubuntu Desktop 64Bit (`UBD64`) | Ubuntu Server 64Bit (`UBS64`)
* `region` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Default: `KR` region.

## Attributes Reference

* `no` - Member server image no
* `name` - Member server image name
* `description` - Member server image description
* `original_server_instance_no` - Original server instance no
* `original_server_product_code` - Original server product code
* `original_server_name` - Original server name
* `original_base_block_storage_disk_type` - Original base block storage disk type
* `original_server_image_product_code` - Original server image product code
* `original_os_information` - Original os information
* `original_server_image_name` - Original server image name
* `status_name` - Member server image status name
* `status` - Member server image status
* `operation` - Member server image operation
* `platform_type` - Member server image platform type
* `region` - Region info
* `block_storage_total_rows` - Member server image block storage total rows
* `block_storage_total_size` - Member server image block storage total size
