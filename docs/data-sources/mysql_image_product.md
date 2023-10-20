---
subcategory: "MySQL"
---


# Data Source: ncloud_mysql_image_product

You should select a MySQL product (MySQL specification) to create a MySQL instance (SW).
To this end, we provide data source by which you can search a MySQL product.

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

* `product_code` - (Optional) Product code you want to view on the list. Use this for a single search.
* `exclusion_product_code` - (Optional) Product code you want to exclude on the list.
* `generation_code` - (Optional) Generation code. The available values are as follows: G2 | G3
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of MySQL product.
* `product_name` - MySQL product name.
* `product_type` - MySQL Product type code.
* `product_description` - MySQL product description.
* `platform_type` - Platform type.
* `os_information` - OS information.
