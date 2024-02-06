---
subcategory: "Hadoop"
---


# Data Source: ncloud_hadoop_image_products 

To create a hadoop instance, you should select a hadoop image product code. This data source gets a list of hadoop images.

~> **NOTE:** This only supports VPC environment.

## Example Usage

The following example shows how to take a list of Hadoop image.

```terraform
data "ncloud_hadoop_image_products" "example" {
  output_file = "image.json"
}

output "image_list" {
  value = {
    for image in data.ncloud_hadoop_image_products.example.image_product_list:
    image.product_name => image.product_code
  }
}
```

```terraform
data "ncloud_hadoop_image_products" "example" {
  filter {
    name = "product_name"
    values = ["Cloud Hadoop 2.1"]
  }
}
```

Outputs:
```terraform
image_list = {
  "Cloud Hadoop 1.8": "SW.VCHDP.LNX64.CNTOS.0708.HDP.18.B050",
  "Cloud Hadoop 1.9": "SW.VCHDP.LNX64.CNTOS.0708.HDP.19.B050",
  "Cloud Hadoop 2.0": "SW.VCHDP.LNX64.CNTOS.0708.HDP.20.B050",
  "Cloud Hadoop 2.1": "SW.VCHDP.LNX64.CNTOS.0708.HDP.21.B050",
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`. 
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Reuired) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `image_product_list` - List of HADOOP image product.
  * `product_code` - The code of image product.
  * `generation_code` - Generation code. (Only `G2` or `G3`)
  * `product_name` - Product name.
  * `product_type` - Product type code.
  * `platform_type` - Platform type code.
  * `os_information` - OS Information.
