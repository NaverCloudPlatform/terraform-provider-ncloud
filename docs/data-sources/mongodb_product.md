---
subcategory: "MongoDb"
---


# Data Source: ncloud_mongodb_product

Provides the list of available Cloud DB for MongoDB server specification codes.

## Example Usage

```hcl
data "ncloud_mongodb_product" "product" {
  image_product_code = "SW.VMGDB.LNX64.CNTOS.0708.MNGDB.402.CE.B050"
  product_code = "SVR.VMGDB.MNGOS.STAND.C002.M008.NET.SSD.B050.G002"
  infra_resource_detail_type_code = "MNGOS"
  exclusion_product_code = "SVR.VMGDB.MNGOS.HICPU.C004.M008.NET.SSD.B050.G00"
  
  filter {
    name   = "product_type"
    values = ["STAND"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `image_product_code` - (Required) You can get one from `data ncloud_mongodb_image_product`. This is a required value, and each available mysql's specification varies depending on the mongodb image product.
* `product_code` - (Optional) Enter a product code to search from the list. Use it for a single search.
* `exclusion_product_code` - (Optional) Enter a product code to exclude.
* `infra_resource_detail_type_code` - (Optional) Cloud for MongoDb Other Server infra Detailed Product Code. Options : MNGOD | MNGOS | ARBIT | CFGSV
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by
    * `values` - (Required) Set of values that are accepted for the given field.

## Attributes Reference

* `id` - The ID of mongodb product.
* `product_name` - Product name.
* `product_type` - Product type code.
* `product_description` - Product description.
* `infra_resource_type` - Infra resource type code.
* `cpu_count` - CPU count.
* `memory_size` - Memory size.
* `disk_type` - Disk type.