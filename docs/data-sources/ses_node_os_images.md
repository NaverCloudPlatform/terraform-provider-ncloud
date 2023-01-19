# Data Source: ncloud_ses_node_os_images

Provides list of available Server OS images.

## Example Usage

```hcl
data "ncloud_ses_node_os_images" "all_images" {}

data "ncloud_ses_node_os_images" "CentOS_7-8" {
  filter {
    name = "name"
    values = ["CentOS 7.8 (64-bit)"]
  }
}

data "ncloud_ses_node_os_images" "CentOS_7-8" {
  filter {
    name = "id"
    values = ["SW.VELST.OS.LNX64.CNTOS.0708.B050"]
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

* `images` - A List of OS image product.

### OS Image Product Reference
`images` are also exported with the following attributes, when there are relevant: Each element supports the following:

* `id` - The ID of OS image product.
* `name` - OS image name