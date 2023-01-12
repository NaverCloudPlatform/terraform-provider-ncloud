# Data Source: ncloud_cdss_node_product

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

data "ncloud_cdss_node_product" "node_sample" {
  os_image  = data.ncloud_cdss_os_image.os_image_sample.id
  subnet_no = ncloud_subnet.selected.id
  
  filter {
    name   = "cpu_count"
    values = ["2"]
  }

  filter {
    name   = "memory_size"
    values = ["8GB"]
  }

  filter {
    name   = "product_type"
    values = ["STAND"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `os_image` - (Required) OS type to be used.
* `subnet_no` - (Required) Subnet number where the node will be located.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of server product.
* `cpu_count` - CPU count.
* `memory_size` - Memory size.
* `product_type` - Product type code.
    