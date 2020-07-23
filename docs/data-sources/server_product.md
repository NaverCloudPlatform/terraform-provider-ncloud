---
layout: "ncloud"
page_title: "NCLOUD: ncloud_server_product"
sidebar_current: "docs-ncloud-datasource-server-product"
description: |-
  Get server product
---

# Data Source: ncloud_server_product

ou should select a server product (server specification) to create a server instance (VM).
To this end, we provide data source by which you can search a server product.

## Example Usage

```hcl
data "ncloud_server_product" "product" {
	server_image_product_code = "SPSW0LINUX000032"
}
```

## Argument Reference

The following arguments are supported:

* `server_image_product_code` - (Required) You can get one from `data ncloud_server_images`. This is a required value, and each available server's specification varies depending on the server image product.
* `product_name_regex` - (Optional) A regex string to apply to the Server Product list returned.
* `exclusion_product_code` - (Optional) Enter a product code to exclude from the list.
* `product_code` - (Optional) Enter a product code to search from the list. Use it for a single search.
* `region` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Default: KR region.
* `zone` - (Optional) Zone code. You can decide a zone where servers are created. You can decide which zone the product list will be requested at. default : Select the first Zone in the specific region
    Get available values using the data source `ncloud_zones`.
* `internet_line_type_code` - (Optional) Internet line code. PUBLC(Public), GLBL(Global)

## Attributes Reference

* `product_name` - Product name
* `product_type` - Product type code
* `product_description` - Product description
* `infra_resource_type` - Infra resource type code
* `cpu_count` - CPU count
* `memory_size` - Memory size
* `base_block_storage_size` - Base block storage size
* `platform_type` - Platform type code
* `os_information` - OS Information
* `add_block_storage_size` - Additional block storage size
