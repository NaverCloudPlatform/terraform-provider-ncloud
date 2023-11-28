---
subcategory: "Hadoop"
---


# Data Source: ncloud_hadoop_image

To create a hadoop instance, you should select a hadoop image. This data source gets a hadoop image.

## Example Usage

The following example shows how to take Hadoop image.

```hcl
data "ncloud_hadoop_image" "image" {
  product_code = "SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"
}

data "ncloud_hadoop" "hadoop_by_filter" {
  filter {
    name = "product_code"
    values = ["SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"]
  }
}
```

## Argument Reference

* `product_code` - (Optional) The product code you want to view on the list. Use this when searching for 1 product.
* `exclusion_product_code` - (Optional)
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Reuired) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of Hadoop image product.
* `product_code` - Product code of Hadoop image product.
* `product_name` - Product name of Hadoop image product.
* `product_type` - Product type of Hadoop image product.
* `product_description` - Product description of Hadoop image product.
* `infra_resource_type` - Infra resource type of Hadoop image product.
* `base_block_storage_size` - Base block storage size of Hadoop image product.
* `platform_type` - Platform type of Hadoop image product.
* `os_information` - OS information of Hadoop image product.
* `generation_code` - Generation code of Hadoop image product.