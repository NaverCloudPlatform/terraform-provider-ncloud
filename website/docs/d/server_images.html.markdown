---
layout: "ncloud"
page_title: "NCLOUD: ncloud_server_images"
sidebar_current: "docs-ncloud_server_images"
description: |-
  Get server image product list
---

# Data Source: ncloud_server_images

To create a server instance (VM), you should select a server image. This API gets a list of server images.

## Example Usage

```hcl
data "ncloud_server_images" "linux" {
}
```

## Argument Reference

The following arguments are supported:

* `product_name_regex` - (Optional) A regex string to apply to the server image list returned by ncloud.
* `exclusion_product_code` - (Optional) Product code you want to exclude from the list.
* `product_code` - (Optional) Product code you want to view on the list. Use this when searching for 1 product.
* `platform_type_code_list` - (Optional) Values required for identifying platforms in list-type.
    The available values are as follows: Linux 32Bit(LNX32) | Linux 64Bit(LNX64) | Windows 32Bit(WND32) | Windows 64Bit(WND64) | Ubuntu Desktop 64Bit(UBD64) | Ubuntu Server 64Bit(UBS64)
* `block_storage_size` - (Optional) Block storage size.
* `region_no` - (Optional) Region number. Get available values using the `data ncloud_regions`.
* `infra_resource_detail_type_code` - (Optional) infra resource detail type code.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.

## Attributes Reference

* `id` - ID of server images.
* `server_images` - A List of server image
    * `product_code` - Product code
    * `product_name` - Product name
    * `product_type` - Product type
        * `code` - Product type code
        * `name` - Product type name
    * `product_description` - Product description
    * `infra_resource_type` - Infra resource type
        * `code` - Infra resource type code
        * `name` - Infra resource type name
    * `cpu_count` - CPU count
    * `memory_size` - Memory size
    * `base_block_storage_size` - Base block storage size
    * `platform_type` - Platform type
        * `code` - Platform type code
        * `name` - Platform type name
    * `os_information` - OS Information
    * `add_block_storage_size` - additional block storage size