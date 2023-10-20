---
subcategory: "MySQL"
---


# Data Source: ncloud_mysql_image_products

Get a list of MySQL image products.

## Example Usage

```hcl
data "ncloud_mysql_image_products" "images_by_code" {
  product_code = "SW.VDBAS.DBAAS.LNX64.CNTOS.0708.MYSQL.5732.B050"

  output_file = "image.json"
}

output "image_list" {
  value = {
    for image in data.ncloud_mysql_image_products.images.image_product_list:
    image.product_name => image.product_code
  }
}
```

```hcl
data "ncloud_mysql_image_products" "images_by_filter" {
  filter {
    name = "product_code"
    values = ["SW.VDBAS.DBAAS.LNX64.CNTOS.0708.MYSQL.5732.B050"]
  }

  output_file = "image.json"
}

output "image_list" {
  value = {
    for image in data.ncloud_mysql_image_products.images.image_product_list:
    image.product_name => image.product_code
  }
}
```

Outputs:
```hcl
image_list = {
  "mysql(5.7.32)": "SW.VDBAS.DBAAS.LNX64.CNTOS.0708.MYSQL.5732.B050",
}
```
## Argument Reference

The following arguments are supported:

* `product_code` - (Optional) Product code you want to view on the list. Use this for a single search.
* `exclusion_product_code` - (Optional) Product code you want to exclude on the list.
* `generation_code` - (Optional) Generation code. The available values are as follows: G2 | G3
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of MySQL image product. (ID is UTC time when data source was created)
* `image_product_list` - List of MySQL image product.

### MySQL Image Product Reference

`image_product_list` are also exported with the following attributes, when there are relevant: Each element supports the following:

* `product_code` - The ID of MySQL image product.
* `generation_code` - Generation code. (Only `G2` or `G3`)
* `product_name` - Product name.
* `product_type` - Product type code.
* `infra_resource_type` - Infra resource type code.
* `product_description` - Product description.
* `platform_type` - Platform type code.
* `os_information` - OS Information.