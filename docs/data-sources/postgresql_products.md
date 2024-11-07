---
subcategory: "PostgreSQL"
---

# Data Soruce: ncloud_postgresql_products

Get a list of PostgreSQL products.

~> **NOTE:** This only supports VPC environment.

## Example Usage

```terraform
data "ncloud_postgresql_products" "all" {
    image_product_code = "SW.VPGSL.OS.LNX64.ROCKY.0810.PGSQL.B050"
    filter {
        name = "product_type" 
        values = ["STAND"]
    }
    output_file = products.json"
}
output "product_list" {
    value = {
        for product in data.ncloud_postgresql_products.all.product.list :
        image.product_name => image.product_code
    }
}
```

Outputs:
`terraform
list_image = {
    "vCPU 2EA, Memory 8GB": "SVR.VPGSL.DBSVR.STAND.C002.M008.NET.SSD.B050.G002",
    "vCPU 4EA, Memory 16GB": "SVR.VPGSL.DBSVR.STAND.C004.M016.NET.SSD.B050.G002",
    "vCPU 8EA, Memory 32GB": "SVR.VPGSL.DBSVR.STAND.C008.M032.NET.SSD.B050.G002",
    "vCPU 16EA, Memory 64GB": "SVR.VPGSL.DBSVR.STAND.C016.M064.NET.SSD.B050.G002",
    "vCPU 32EA, Memory 128GB": "SVR.VPGSL.DBSVR.STAND.C032.M128.NET.SSD.B050.G002"
}

## Argument Reference

The following arguments are supported:

* `image_product_code` - (Required) Youc can get one from `data.ncloud_postgresql_image_products`, This is a required value, and each available PostgreSQL's specification varies depending on the PostgreSQL image product.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributees Reference

This data source exports the following attributes in addition to the argument above:

* `product_list` - List of PostgreSQL product.
  * `product_name` - Product name.
  * `product_code` - Product code.
  * `product_type` - Product type code.
  * `product_description` - Product description.
  * `engine_version_code` - The engine version of the specific PostgreSQL.
  * `infra_resource_type` - Infra resource type code.
  * `cpu_count` - CPU count.
  * `memory_size` - Memory size.
  * `disk_type` - Disk type.