---
subcategory: "Cloud Data Streaming Service"
---


# Data Source: ncloud_cdss_node_products

## Example Usage

```hcl
variable "subnet_no" {}

data "ncloud_subnet" "selected" {
  id = var.subnet_no
}

data "ncloud_cdss_os_image" "os_image_sample" {
  filter {
    name   = "product_name"
    values = ["CentOS 7.8 (64-bit)"]
  }
}

data "ncloud_cdss_node_products" "nodes_sample" {
  os_image  = data.ncloud_cdss_os_image.os_image_sample.id
  subnet_no = ncloud_subnet.selected.id
}
```

## Argument Reference

The following arguments are supported:

* `os_image` - (Required) OS type to be used.
* `subnet_no` - (Required) Subnet number where the node will be located.

## Attributes Reference

* `node_products` - A list of Server product
