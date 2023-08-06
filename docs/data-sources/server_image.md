---
subcategory: "Server"
---


# Data Source: ncloud_server_image

To create a server instance (VM), you should select a server image. This data source gets a server image.

## Example Usage

* Filter by product name

```hcl
data "ncloud_server_image" "image" {
  product_name = "CentOS 7.3 (64-bit)"
}

// Use filter
data "ncloud_server_image" "image_by_filter" {
  filter {
    name = "product_name"
    values = ["CentOS 7.3 (64-bit)"]
  }
}
```

* Filter by product type

```hcl
data "ncloud_server_image" "image" {
  platform_type = "WINNT"
}

// Use filter
data "ncloud_server_image" "image_by_filter" {
  filter {
    name = "platform_type"
    values = ["WINNT"]
  }
}
```



## Argument Reference

The following arguments are supported:

* `product_code` - (Optional) Product code you want to view on the list. Use this when searching for 1 product.
* `platform_type` - (Optional) Values required for identifying platform.
    The available values are as follows: Linux 32Bit(LNX32) | Linux 64Bit(LNX64) | Windows 32Bit(WND32) | Windows 64Bit(WND64) | Ubuntu Desktop 64Bit(UBD64) | Ubuntu Server 64Bit(UBS64)
* `infra_resource_detail_type_code` - (Optional) infra resource detail type code.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of server image product.
* `product_name` - Product name
* `product_type` - Product type code.
* `product_description` - Product description.
* `infra_resource_type` - Infra resource type code.
* `base_block_storage_size` - Base block storage size.
* `os_information` - OS Information.
