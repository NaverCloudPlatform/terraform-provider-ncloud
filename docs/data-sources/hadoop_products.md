---
subcategory: "Hadoop"
---


# Data Source: ncloud_hadoop_products

Provides the list of Hadoop server specification codes.

~> **NOTE:** This only supports VPC environment.

## Example Usage

```terraform
data "ncloud_hadoop_products" "all" {
  image_product_code = "SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"

  filter {
    name   = "product_type"
    values = ["STAND"]
  }

  output_file = "products.json"
}

output "product_list" {
  value = {
  for product in data.ncloud_hadoop_products.all.product_list :
    product.product_code => product.product_name
  }
}
```

Outputs:
```terraform
product_list = {
  "SVR.VCHDP.EDGND.STAND.C004.M016.NET.HDD.B050.G002" = "vCPU 4EA, Memory 16GB, Disk 50GB"
  "SVR.VCHDP.EDGND.STAND.C008.M032.NET.HDD.B050.G002" = "vCPU 8EA, Memory 32GB, Disk 50GB"
  "SVR.VCHDP.MSTDT.STAND.C004.M016.NET.HDD.B050.G002" = "vCPU 4EA, Memory 16GB, Disk 50GB"
  "SVR.VCHDP.MSTDT.STAND.C008.M032.NET.HDD.B050.G002" = "vCPU 8EA, Memory 32GB, Disk 50GB"
  "SVR.VCHDP.MSTDT.STAND.C016.M064.NET.HDD.B050.G002" = "vCPU 16EA, Memory 64GB, Disk 50GB"
  "SVR.VCHDP.MSTDT.STAND.C032.M128.NET.HDD.B050.G002" = "vCPU 32EA, Memory 128GB, Disk 50GB"
}
```


## Argument Reference

The following arguments are supported:

* `image_product_code` - (Required) You can get one from `data.ncloud_hadoop_images`. This is a required value, and each available Hadoop's specification varies depending on the hadoop image product.
* `product_code` - (Required) Product code you want to view on the list. Use this when searching for 1 product.
* `infra_resource_detail_type_code` - (Optional) Hadoop Other Server infra Detailed Product Code. Options : MSTDT | EDGND
* `exclusion_product_code` - (Optional) Product code you want to exclude on the list.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `product_list` - Hadoop product list
  * `product_name` - Product name.
  * `product_code` - Product code.
  * `product_type` - Product type code.
  * `product_description` - Product description.
  * `infra_resource_type` - Infra resource type code.
  * `infra_resource_detail_type` - Infra resource detail type code.
  * `cpu_count` - CPU count.
  * `memory_size` - Memory size.
  * `disk_type` - Disk type.
