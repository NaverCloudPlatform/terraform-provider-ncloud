---
subcategory: "Hadoop"
---


# Data Source: ncloud_hadoop_add_on

This module can be useful for getting add-ons of Hadoop.

## Example Usage

#### Basic usage

The following example shows how to take add-ons of Hadoop.

```hcl
data "ncloud_hadoop_add_on" "addon" {
	image_product_code= "SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"
	cluster_type_code= "CORE_HADOOP_WITH_SPARK"
}

data "ncloud_hadoop_add_on" "addon_output_file" {
  image_product_code= "SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"
  cluster_type_code= "CORE_HADOOP_WITH_SPARK"
  
  output_file = "hadoop_add_on.json"
}
```

## Argument Reference

* `image_product_code` - (Required) The image product code of the specific Hadoop add-on to retrieve.
* `cluster_type_code` - (Required) The cluster type code of the specific Hadoop add-on to retrieve.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
 
## Attributes Reference

* `id` - The ID of add-on list.
* `add_on_list` - The add-on list of Hadoop.