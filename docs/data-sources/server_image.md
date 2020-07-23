---
layout: "ncloud"
page_title: "NCLOUD: ncloud_server_image"
sidebar_current: "docs-ncloud-datasource-server-image"
description: |-
  Get server image product
---

# Data Source: ncloud_server_image

To create a server instance (VM), you should select a server image. This data source get a server image.

## Example Usage

* Filter by product name

```hcl
data "ncloud_server_image" "image" {
  product_name_regex = "^Windows Server 2012(.*)"
}
```

* Filter by product type

```hcl
data "ncloud_server_image" "image" {
  product_type_code = "WINNT"
}
```

## Argument Reference

The following arguments are supported:

* `product_name_regex` - (Optional) A regex string to apply to the server image list returned by ncloud.
* `exclusion_product_code` - (Optional) Product code you want to exclude from the list.
* `product_code` - (Optional) Product code you want to view on the list. Use this when searching for 1 product.
* `product_type` - (Optional) Product type code
* `platform_type_code_list` - (Optional) Values required for identifying platforms in list-type.
    The available values are as follows: Linux 32Bit(LNX32) | Linux 64Bit(LNX64) | Windows 32Bit(WND32) | Windows 64Bit(WND64) | Ubuntu Desktop 64Bit(UBD64) | Ubuntu Server 64Bit(UBS64)
* `block_storage_size` - (Optional) Block storage size.
* `region` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Default: KR region.
* `infra_resource_detail_type_code` - (Optional) infra resource detail type code.

## Attributes Reference

* `product_name` - Product name
* `product_description` - Product description
* `infra_resource_type` - Infra resource type code
* `cpu_count` - CPU count
* `memory_size` - Memory size
* `base_block_storage_size` - Base block storage size
* `platform_type` - Platform type code
* `os_information` - OS Information
* `add_block_storage_size` - Additional block storage size
