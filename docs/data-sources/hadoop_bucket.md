---
subcategory: "Hadoop"
---


# Data Source: ncloud_hadoop_bucket

This module can be useful for getting buckets to create hadoop instance.

## Example Usage

#### Basic usage

The following example shows how to take buckets.

```hcl
data "ncloud_hadoop_bucket" "bucket" {

}

data "ncloud_hadoop_add_on" "bucket_output_file" {
  output_file = "hadoop_bucket.json"
}
```

## Argument Reference

* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.

## Attributes Reference

* `id` - The ID of bucket list.
* `bucket_list` - The bucket list of Hadoop.