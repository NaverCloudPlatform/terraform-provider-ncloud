---
subcategory: "MongoDb"
---


# Data Source: ncloud_mongodb_image_products


Provides MongoDB a number of image products to create MongoDB instance.

## Example Usage

```hcl
data "ncloud_mongodb_image_products" "mongodb_products" {
  filter {
    name   = "product_name"
    values = ["MongoDB 4.2.23 Community Edition"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `product_code` - (Optional) Product code.
* `exclusion_product_code` - (Optional) Product code to exclude.
* `generation_code` - (Optional) Generation code. It can be `G2` or `G3`.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of mongodb products.
* `image_product_list` - MongoDb Image product list

The `image_product_list` object support the following:

* `product_code` - MongoDb product code.
* `generation_code` - Generation code.
* `product_name` - MongoDb product name.
* `product_type` - MongoDb Product type code.
* `infra_resource_type` - MongoDb infra resource type.
* `product_description` - MongoDb product description.
* `platform_type` - Platform type.
* `os_information` - OS information.

