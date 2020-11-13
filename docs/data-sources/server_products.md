# Data Source: ncloud_server_products

You should select a server product (server specification) to create a server instance (VM).
To this end, we provide data source by which you can search a server product.

## Example Usage

```hcl
# Classic
data "ncloud_server_products" "product_ids" {
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
* `zone` - (Optional) Zone code. You can decide a zone where servers are created. You can decide which zone the product list will be requested at. default : Select the first Zone in the specific region
    Get available values using the data source `ncloud_zones`.
* `internet_line_type_code` - (Optional) Internet line code. PUBLC(Public), GLBL(Global)
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.


## Attributes Reference

* `server_products` - A List of server product code
