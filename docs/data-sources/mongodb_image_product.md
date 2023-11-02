---
subcategory: "MongoDB"
---


# Data Source: ncloud_mongodb_image_product

Provides MongoDB image product to create MongoDB instance.

## Example Usage

```hcl
data "ncloud_mongodb_image_product" "mongodb_image_product" {
  product_code = "SW.VMGDB.LNX64.CNTOS.0708.MNGDB.4223.CE.B050"
  generation_code = "G2"

  filter {
    name   = "product_name"
    values = ["MongoDB 4.2.23 Community Edition"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `product_code` - (Optional) Product code to search from a mongodb product list. Use it for a single search.
* `exclusion_product_code` - (Optional) Product code to exclude.
* `generation_code` - (Optional) Generation code. It can be `G2` or `G3`.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of mongodb product.
* `product_name` - MongoDB product name.
* `product_type` - MongoDB Product type code.
* `infra_resource_type` - MongoDB infra resource type.
* `product_description` - MongoDB product description.
* `platform_type` - Platform type.
* `os_information` - OS information.