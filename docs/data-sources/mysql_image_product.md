---
subcategory: "Mysql"
---


# Data Source: ncloud_mysql_image_product

You should select a mysql product (mysql specification) to create a mysql instance (SW).
To this end, we provide data source by which you can search a mysql product.

## Example Usage

```hcl
data "ncloud_mysql_image_product" "product" {
  product_code = "SW.VDBAS.DBAAS.LNX64.CNTOS.0708.MYSQL.5732.B050"

  filter {
    name   = "product_name"
    values = ["mysql(5.7.32)"]
  }
  
  filter {
    name   = "product_type"
    values = ["LINUX"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `product_code` - (Optional) Enter a product code to search from a mysql product list. Use it for a single search.
* `exclusion_product_code` - (Optional) Enter a product code to exclude.
* `generation_code` - (Optional) Enter a generation code. You can select `G2` or `G3`.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by
    * `values` - (Required) Set of values that are accepted for the given field.

## Attributes Reference

* `id` - The ID of mysql product.
* `product_name` - Mysql product name.
* `product_type` - Mysql Product type code.
* `product_description` - Mysql product description.
* `platform_type` - Platform type.
* `os_information` - OS information.
