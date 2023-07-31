---
subcategory: "Server"
---


# Data Source: ncloud_server_product

You should select a server product (server specification) to create a server instance (VM).
To this end, we provide data source by which you can search a server product.

## Example Usage

```hcl
data "ncloud_server_product" "product" {
  server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"  // Search by 'CentOS 7.3 (64-bit)' image vpc
  // server_image_product_code = "SPSW0LINUX000032"  // Search by 'CentOS 7.3 (64-bit)' image classic

  filter {
    name   = "product_code"
    values = ["SSD"]
    regex  = true
  }

  filter {
    name   = "cpu_count"
    values = ["2"]
  }

  filter {
    name   = "memory_size"
    values = ["8GB"]
  }

  filter {
    name   = "base_block_storage_size"
    values = ["50GB"]
  }

  filter {
    name   = "product_type"
    values = ["STAND"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `server_image_product_code` - (Required) You can get one from `data ncloud_server_images`. This is a required value, and each available server's specification varies depending on the server image product.
* `product_code` - (Optional) Enter a product code to search from the list. Use it for a single search.
* `zone` - (Optional) Zone code. You can decide a zone where servers are created. You can decide which zone the product list will be requested in. default : Select the first Zone in the specific region.
    Get available values using the data source `ncloud_zones`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of server product.
* `product_name` - Product name.
* `product_type` - Product type code.
* `product_description` - Product description.
* `infra_resource_type` - Infra resource type code.
* `cpu_count` - CPU count.
* `memory_size` - Memory size.
* `disk_type` - Disk type.
* `base_block_storage_size` - Base block storage size.
* `generation_code` - Generation Code.
