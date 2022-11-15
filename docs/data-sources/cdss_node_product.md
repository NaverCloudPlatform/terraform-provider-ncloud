# Data Source: ncloud_cdss_node_product

## Example Usage

```hcl
data "ncloud_cdss_os_product" "os_sample" {
  filter {
    name = "product_name"
    values = ["CentOS 7.8 (64-bit)"]
  }
}

data "ncloud_cdss_node_product" "node_sample" {
  os_product_code = data.ncloud_cdss_os_product.os_sample.id
  subnet_no       = "3438"
  
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
* `os_product_code` - (Required) OS type to be used.
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
    