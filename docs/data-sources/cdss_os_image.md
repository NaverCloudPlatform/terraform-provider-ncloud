---
subcategory: "Cloud Data Streaming Service"
---


# Data Source: ncloud_cdss_os_image

## Example Usage

```hcl
data "ncloud_cdss_os_image" "sample_01" {
  filter {
    name   = "id"
    values = ["SW.VCDSS.OS.LNX64.CNTOS.0708.B050"]
  }
}

data "ncloud_cdss_os_image" "sample_02" {
  filter {
    name   = "image_name"
    values = ["CentOS 7.8 (64-bit)"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of server image product.
* `image_name` - Os image name
