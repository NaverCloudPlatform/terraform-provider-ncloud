---
subcategory: "MySQL"
---


# Data Source: ncloud_mysql_products

Get a list of MySQL products.

~> **NOTE:** This only supports VPC environment.

## Example Usage

```terraform
data "ncloud_mysql_products" "all" {
  image_product_code = "SW.VDBAS.DBAAS.LNX64.CNTOS.0708.MYSQL.8025.B050"

  filter {
    name   = "product_type"
    values = ["STAND"]
  }

  output_file = "products.json"
}

output "product_list" {
  value = {
  for product in data.ncloud_mysql_products.all.product_list :
    product.product_name => product.product_code
  }
}
```

Outputs:
```terraform
list_image = {
  "vCPU 2EA, Memory 8GB" = "SVR.VDBAS.STAND.C002.M008.NET.HDD.B050.G002"
}
```

## Argument Reference

The following arguments are supported:

* `image_product_code` - (Required) You can get one from `data.ncloud_mysql_image_products`. This is a required value, and each available MySQL's specification varies depending on the MySQL image product.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `product_list` - List of MySQL product.
  * `product_name` - Product name.
  * `product_code` - Product code.
  * `product_type` - Product type code.
  * `product_description` - Product description.
  * `infra_resource_type` - Infra resource type code.
  * `cpu_count` - CPU count.
  * `memory_size` - Memory size.
  * `disk_type` - Disk type.
