---
subcategory: "Mssql"
---


# Data Source: ncloud_mssql_image_products

Get a list of MSSQL image products.

~> **NOTE:** This only supports VPC environment.

## Example Usage

```terraform
data "ncloud_mssql_image_products" "example" {
  output_file = "image.json"
}

output "image_list" {
  value = {
    for image in data.ncloud_mssql_image_products.example.image_product_list:
    image.product_name => image.product_code
  }
}
```

```terraform
data "ncloud_mssql_image_products" "example" {
  filter {
    name = "product_code"
    values = ["SW.VMSSL.OS.WND64.WINNT.SVR2016.MSSQL.15042981.SE.B100"]
  }
}
```

Outputs:
```terraform
image_list = {
  "MSSQL 15.0.4298.1 Standard Edition": "SW.VMSSL.OS.WND64.WINNT.SVR2016.MSSQL.15042981.SE.B100",
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `image_product_list` - List of MSSQL image product.
  * `product_code` - The ID of MSSQL image product.
  * `generation_code` - Generation code. (Only `G2` or `G3`)
  * `product_name` - Product name.
  * `product_type` - Product type code.
  * `platform_type` - Platform type code.
  * `os_information` - OS Information.
