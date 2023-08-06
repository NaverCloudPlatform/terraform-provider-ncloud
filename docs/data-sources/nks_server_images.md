---
subcategory: "Kubernetes Service"
---


# Data Source: ncloud_nks_server_images

Provides list of available Kubernetes Nodepool ServerImages.

## Example Usage

```hcl
data "ncloud_nks_server_images" "images" {}

data "ncloud_nks_server_images" "ubuntu20" {
  filter {
    name = "label"
    values = ["ubuntu-20.04-64-server"]
    regex = true
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

* `images` - A list of ServerImages
  * `label` - ServerImage name
  * `value` - ServerImage code
