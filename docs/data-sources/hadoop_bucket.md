---
subcategory: "Hadoop"
---


# Data Source: ncloud_hadoop_bucket

This module can be useful for getting Object Storage buckets to create hadoop instance.

~> **NOTE:** This only supports VPC environment.

## Example Usage

The following example shows how to take buckets.

```terraform
data "ncloud_hadoop_bucket" "bucket" {

}

data "ncloud_hadoop_add_on" "bucket_output_file" {
  output_file = "hadoop_bucket.json"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `bucket_list` - The Object Storage bucket list of Hadoop.
