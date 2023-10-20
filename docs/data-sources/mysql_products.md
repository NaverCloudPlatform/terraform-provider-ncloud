---
subcategory: "MySQL"
---


# Data Source: ncloud_mysql_products

Get a list of MySQL products.

## Example Usage

```hcl
data "ncloud_mysql_products" "all" {
  cloud_mysql_image_product_code = "SW.VDBAS.DBAAS.LNX64.CNTOS.0708.MYSQL.5732.B050"
  product_code = "SVR.VDBAS.STAND.C002.M008.NET.HDD.B050.G002"
  exclusion_product_code = "SVR.VDBAS.HICPU.C004.M008.NET.HDD.B050.G002"

  filter {
    name   = "product_type"
    values = ["STAND"]
  }

  output_file = "image.json"
}

output "product_list" {
  value = {
  for product in data.ncloud_mysql_products.all.product_list :
    image.product_name => image.product_code
  }
}
```

Outputs:
```hcl
list_image = {
  "vCPU 2EA, Memory 8GB" = "SVR.VDBAS.STAND.C002.M008.NET.HDD.B050.G002"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_mysql_image_product_code` - (Required) You can get one from `data ncloud_mysql_image_product`. This is a required value, and each available MySQL's specification varies depending on the MySQL image product.
* `product_code` - (Optional) Product code you want to view on the list. Use this for a single search.
* `exclusion_product_code` - (Optional) Product code you want to exclude on the list.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of MySQL product. (ID is UTC time when data source was created)
* `product_list` - List of MySQL product.

### MySQL Product Reference

`product_list` are also exported with the following attributes, when there are relevant: Each element supports the following:

* `product_name` - Product name.
* `product_code` - Product code.
* `product_type` - Product type code.
* `product_description` - Product description.
* `infra_resource_type` - Infra resource type code.
* `cpu_count` - CPU count.
* `memory_size` - Memory size.
* `disk_type` - Disk type.
