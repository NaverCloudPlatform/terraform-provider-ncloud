---
subcategory: "Hadoop"
---


# Data Source: ncloud_hadoop_images

To create a hadoop instance, you should select a hadoop image. This data source gets a list of hadoop image.

## Example Usage

The following example shows how to take a list of Hadoop image.

```hcl
data "ncloud_hadoop_images" "images" {
  product_code = "SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"
}

data "ncloud_hadoops" "hadoop_by_filter" {
  filter {
    name = "product_code"
    values = ["SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"]
  }
}
```

## Argument Reference

* `product_code` - (Optional) The product code you want to view on the list. Use this when searching for 1 product.
* `exclusion_product_code` - (Optional)
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`. 
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Reuired) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

The following attrivutes are exported:

* `id` - The ID of list of hadoop image.
* `images` - The list of Hadoop image.

### Hadoop Reference

`images` are also exported with the following attributes, where relevant: Each element supports the following:

* `product_code` - Product code of Hadoop image product.
* `product_name` - Product name of Hadoop image product.
* `product_type` - Product type of Hadoop image product.
* `product_description` - Product description of Hadoop image product.
* `infra_resource_type` - Infra resource type of Hadoop image product.
* `base_block_storage_size` - Base block storage size of Hadoop image product.
* `platform_type` - Platform type of Hadoop image product.
* `os_information` - OS information of Hadoop image product.
* `generation_code` - Generation code of Hadoop image product.