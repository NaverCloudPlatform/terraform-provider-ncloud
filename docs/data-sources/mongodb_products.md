---
subcategory: "MongoDB"
---


# Data Source: ncloud_mongodb_products

Provides the list of available Cloud DB for MongoDB server specification codes.

## Example Usage

```hcl
data "ncloud_mongodb_products" "all" {
  image_product_code = "SW.VMGDB.LNX64.CNTOS.0708.MNGDB.4223.CE.B050"
}
```

## Argument Reference

The following arguments are supported:

* `image_product_code` - (Required) You can get one from `data ncloud_mongodb_image_product`. This is a required value, and each available MongoDB's specification varies depending on the mongodb image product.
* `product_code` - (Optional) Product code to search from the list. Use it for a single search.
* `exclusion_product_code` - (Optional) Product code to exclude.
* `infra_resource_detail_type_code` - (Optional) Cloud for MongoDB Other Server infra Detailed Product Code. Options : MNGOD | MNGOS | ARBIT | CFGSV
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of mongodb product.
* `product_list` - MongoDB product list

The `product_list` object support following:

* `product_code` - Product code.
* `product_name` - Product name.
* `product_type` - Product type code.
* `product_description` - Product description.
* `infra_resource_type` - Infra resource type code.
* `cpu_count` - CPU count.
* `memory_size` - Memory size.
* `disk_type` - Disk type.