---
subcategory: "PostgreSQL"
---

# Data Source: ncloud_postgresql_image_products

Get a list of PostgreSQL image products.

~> **NOTE:** This only supports VPC environment.

## Example Usage

```terraform
data "ncloud_postgresql_image_products" "example" {
    output_file = "image.json"
}
output "image_list" {
    value = {
        for image in data.ncloud_postgresql_image_products.example.image_product_list:
        image.engine_version_code => image.product_code
    }
}
```

```terraform
data "ncloud_postgresql_image_products" "example" {
    filter {
        name = "engine_Version_code"
        values = ["13.13"]
    }
}
```

Outputs:
```terraform
image_list = {
    "13.13": "SW.VPGSL.OS.LNX64.ROCKY.0810.PGSQL.B050"
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

* `image_product_list` - List of PostgreSQL image product.
  * `product_code` - The ID of PostgreSQL image product.
  * `generation_code` - Generation code. (Only `G2` or `G3`)
  * `product_name` - Product name.
  * `product_type` - Product type code.
  * `platform_type` - Platform type code.
  * `os_information` - OS Information.
  * `engine_version_code` - Engine version code.