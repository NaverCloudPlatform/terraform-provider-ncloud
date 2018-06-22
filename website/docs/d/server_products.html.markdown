---
layout: "ncloud"
page_title: "NCLOUD: ncloud_server_products"
sidebar_current: "docs-ncloud-datasource-server-products"
description: |-
Searching a server product list
---

# Data Source: ncloud_server_products

You should select a server product (server specification) to create a server instance (VM).
To this end, we provide API by which you can search a server product.

## Example Usage

```hcl
data "ncloud_server_products" "all" {
  # server_image_product_code: You can get one from `data ncloud_server_images`
  "server_image_product_code" = "SPSW0LINUX000032"
}
```

## Argument Reference

The following arguments are supported:

* `product_name_regex` - (Optional) A regex string to apply to the Server Product list returned.

* `exclusion_product_code` - (Optional) Enter a product code to exclude from the list.

* `product_code` - (Optional) Enter a product code to search from the list. Use it for a single search.

* `server_image_product_code` - (Required) You can get one from `data ncloud_server_images`. This is a required value, and each available server's specification varies depending on the server image product.

* `region_no` - (Optional) You can reach a state in which inout is possible by calling `data ncloud_regions`.

* `zone_no` - (Optional) You can decide a zone where servers are created. You can decide which zone the product list will be requested at.
  You can get one by calling `data ncloud_zones`.
  default : Select the first Zone in the specific region

* `internet_line_type_code` - (Optional) Internet line code. PUBLC(Public), GLBL(Global)

## Attributes Reference

* `server_products` - A List of Server Product
    * `product_code` - Product code
    * `product_name` - Product name
    * `product_type` - Product type
        * `code` - Product type code
        * `name` - Product type name
    * `product_description` - Product description
    * `infra_resource_type` - Infra resource type
        * `code` - Infra resource type code
        * `code_name` - Infra resource type name
    * `cpu_count` - CPU count
    * `memory_size` - Memory size
    * `base_block_storage_size` - Base block storage size
    * `platform_type` - Platform type
        * `code` - Platform type code
        * `code_name` - Platform type name
    * `os_information` - OS Information
    * `add_block_storage_size` - additional block storage size