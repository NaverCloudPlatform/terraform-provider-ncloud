---
subcategory: "MongoDb"
---


# Data Source: ncloud_mongodb_image_product

You should select a mysql product (mysql specification) to create a mysql instance (SW).
To this end, we provide data source by which you can search a mysql product.

## Example Usage

```hcl
data "ncloud_mongodb_image_product" "product" {
  product_code = "SW.VDBAS.VRDS.LNX64.CNTOS.0703.MongoDb.4014.B050"
  generation_code = "G2"

  filter {
    name   = "product_name"
    values = ["MongoDB 4.2 Community Edition"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `product_code` - (Optional) Enter a product code to search from a mongodb product list. Use it for a single search.
* `exclusion_product_code` - (Optional) Enter a product code to exclude.
* `generation_code` - (Optional) Enter a generation code. You can select `G2` or `G3`.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by
    * `values` - (Required) Set of values that are accepted for the given field.

## Attributes Reference

* `id` - The ID of mongodb product.
* `product_name` - MongoDb product name.
* `product_type` - MongoDb Product type code.
* `product_description` - MongoDb product description.
* `platform_type` - Platform type.
* `os_information` - OS information.
* `add_block_storage_size` - Additional block storage size.