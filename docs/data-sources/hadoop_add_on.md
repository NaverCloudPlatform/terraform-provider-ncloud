---
subcategory: "Hadoop"
---


# Data Source: ncloud_hadoop_add_on

This module can be useful for getting add-on list of Hadoop.

~> **NOTE:** This only supports VPC environment.

## Example Usage

The following example shows how to take add-ons of Hadoop.

```terraform
data "ncloud_hadoop_add_on" "addon" {
	image_product_code= "SW.VCHDP.OS.LNX64.ROCKY.0810.HDP.B050"
	cluster_type_code= "CORE_HADOOP_WITH_SPARK"
}

data "ncloud_hadoop_add_on" "addon_output_file" {
  image_product_code= "SW.VCHDP.OS.LNX64.ROCKY.0810.HDP.B050"
  cluster_type_code= "CORE_HADOOP_WITH_SPARK"
  
  output_file = "hadoop_add_on.json"
}
```

## Argument Reference

The following arguments are supported:

* `image_product_code` - (Required) The image product code of the specific Hadoop add-on to retrieve.
* `cluster_type_code` - (Required) The cluster type code of the specific Hadoop add-on to retrieve.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
 
## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `add_on_list` - The add-on list of Hadoop.
