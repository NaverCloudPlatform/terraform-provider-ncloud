# Data Source: ncloud_ses_node_products

Provides list of available Server product.

## Example Usage

```hcl
variable "subnet_no" {}

data "ncloud_ses_node_os_images" "os_images" {}

data "ncloud_ses_node_products" "node_products" {
  os_image_code         = data.ncloud_ses_node_os_images.os_images.images.0.id
  subnet_no             = var.subnet_no
  
  filter {
    name   = "cpu_count"
    values = ["2"]
  }
}
```

## Argument Reference
The following arguments are supported:
* `os_image_code` - (Required) OS type to be used.
* `subnet_no` - (Required) Subnet number where the node will be located.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `codes` - A List of server product.

### Node Product Reference
`codes` are also exported with the following attributes, when there are relevant: Each element supports the following:

* `id` - The value of server product code.
* `cpu_count` - CPU count.
* `memory_size` - Memory size.
* `name` - Product name.
    