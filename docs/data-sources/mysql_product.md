---
subcategory: "Mysql"
---


# Data Source: ncloud_server_product

You should select a mysql server product (mysql server specification) to create a mysql instance (VM).
To this end, we provide data source by which you can search a mysql server product.

## Example Usage

```hcl
data "ncloud_mysql_product" "product" {
  cloud_mysql_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"  // Search by 'CentOS 7.3 (64-bit)' image vpc
  product_code = "SVR.VDBAS.STAND.C002.M008.NET.HDD.B050.G002"
  exclusion_product_code = "SVR.VDBAS.HICPU.C004.M008.NET.HDD.B050.G002"
  
  filter {
    name   = "product_type"
    values = ["STAND"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `cloud_mysql_image_product_code` - (Required) You can get one from `data ncloud_mysql_image_product`. This is a required value, and each available mysql's specification varies depending on the mysql image product.
* `product_code` - (Optional) Enter a product code to search from the list. Use it for a single search.
* `exclusion_product_code` - (Optional) Enter a product code to exclude.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by
    * `values` - (Required) Set of values that are accepted for the given field.
    
## Attributes Reference

* `id` - The ID of mysql product.
* `product_name` - Product name.
* `product_type` - Product type code.
* `product_description` - Product description.
* `infra_resource_type` - Infra resource type code.
* `cpu_count` - CPU count.
* `memory_size` - Memory size.
* `disk_type` - Disk type.
