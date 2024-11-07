---
subcategory: "MongoDB"
---


# Data Source: ncloud_mongodb_products

Provides the list of available Cloud DB for MongoDB server specification codes.

~> **NOTE:** This only supports VPC environment.

## Example Usage

```terraform
data "ncloud_mongodb_products" "all" {
  image_product_code = "SW.VMGDB.OS.LNX64.ROCKY.0810.MNGDB.CE.B050"

  filter {
    name   = "product_type"
    values = ["STAND"]
  }

  output_file = "products.json"
}

output "product_list" {
  value = {
  for product in data.ncloud_mongodb_products.all.product_list :
    product.product_name => product.product_code
  }
}
```

Outputs:
```terraform
product_list = {
  "vCPU 2EA, Memory 8GB": "SVR.VMGDB.MNGOD.STAND.C002.M008.NET.SSD.B050.G002",
  "vCPU 2EA, Memory 8GB": "SVR.VMGDB.CFGSV.STAND.C002.M008.NET.SSD.B050.G002",
  "vCPU 2EA, Memory 8GB": "SVR.VMGDB.MNGOS.STAND.C002.M008.NET.SSD.B050.G002"
}
```


## Argument Reference

The following arguments are supported:

* `image_product_code` - (Required) You can get one from `data.ncloud_mongodb_image_products`. This is a required value, and each available MongoDB's specification varies depending on the mongodb image product.
* `infra_resource_detail_type_code` - (Optional) Cloud for MongoDB Other Server infra Detailed Product Code. Options : MNGOD | MNGOS | ARBIT | CFGSV
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `product_list` - MongoDB product list
  * `product_code` - Product code.
  * `product_name` - Product name.
  * `product_type` - Product type code.
  * `product_description` - Product description.
  * `infra_resource_type` - Infra resource type code.
  * `infra_resource_detail_type` - Infra resource detail type code.
  * `cpu_count` - CPU count.
  * `memory_size` - Memory size.
  * `disk_type` - Disk type.
