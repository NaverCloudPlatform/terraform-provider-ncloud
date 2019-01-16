---
layout: "ncloud"
page_title: "NCLOUD: ncloud_member_server_image"
sidebar_current: "docs-ncloud-datasource-member-server-image"
description: |-
  Get member server image
---

# Data Source: ncloud_member_server_image

Gets a member server image.

## Example Usage

```hcl
data "ncloud_member_server_image" "test" {
  "most_recent" = "true"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the member server image list returned by ncloud
* `no_list` - (Optional) List of member server images to view
* `platform_type_code_list` - (Optional) List of platform codes of server images to view. Linux 32Bit (`LNX32`) | Linux 64Bit (`LNX64`) | Windows 32Bit (`WND32`) | Windows 64Bit (`WND64`) | Ubuntu Desktop 64Bit (`UBD64`) | Ubuntu Server 64Bit (`UBS64`)
* `region_code` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Conflicts with `region_no`. Only one of `region_no` and `region_code` can be used.
    Default: KR region.
* `region_no` - (Optional) Region number. Get available values using the data source `ncloud_regions`.
    Conflicts with `region_code`. Only one of `region_no` and `region_code` can be used.
    Default: KR region.
* `most_recent` - (Optional) If more than one result is returned, get the most recent created member server image.

## Attributes Reference

* `no` - Member server image no
* `name` - Member server image name
* `description` - Member server image description
* `original_server_instance_no` - Original server instance no
* `original_server_product_code` - Original server product code
* `original_server_name` - Original server name
* `original_base_block_storage_disk_type` - Original base block storage disk type
    * `code` - Original base block storage disk type code
    * `code_name` - Original base block storage disk type name
* `original_server_image_product_code` - Original server image product code
* `original_os_information` - Original os information
* `original_server_image_name` - Original server image name
* `status_name` - Member server image status name
* `status` - Member server image status
    * `code` - Member server image status code
    * `code_name` - Member server image status name
* `operation` - Member server image operation
    * `code` - Member server image operation code
    * `code_name` - Member server image operation name
* `platform_type` - Member server image platform type
    * `code` - Member server image platform type code
    * `code_name` - Member server image platform type name
* `create_date` - Creation date of the member server image
* `region` - Region info
    * `region_no` - region no
    * `region_code` - Region code
    * `region_name` - Region name
* `block_storage_total_rows` - Member server image block storage total rows
* `block_storage_total_size` - Member server image block storage total size
