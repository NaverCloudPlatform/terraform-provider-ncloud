---
subcategory: "Redis"
---


# Data Source: ncloud_redis_products

Get a list of Redis products.

~> **NOTE:** This only supports VPC environment.

## Example Usage

```terraform
data "ncloud_redis_products" "all" {
  redis_image_product_code = "SW.VRDS.OS.LNX64.ROCKY.0810.REDIS.B050"

  filter {
    name   = "product_type"
    values = ["STAND"]
  }

  output_file = "products.json"
}

output "product_list" {
  value = {
  for product in data.ncloud_redis_products.all.product_list :
    product.product_name => product.product_code
  }
}
```

Outputs:
```terraform
product_list = {
  "Memory 1.5GB": "SVR.VRDS.STAND.C004.M001.NET.SSD.B050.G002",
  "Memory 12.8GB": "SVR.VRDS.STAND.C004.M012.NET.SSD.B050.G002",
  "Memory 25.6GB": "SVR.VRDS.STAND.C004.M025.NET.SSD.B050.G002",
  "Memory 3GB": "SVR.VRDS.STAND.C004.M003.NET.SSD.B050.G002",
  "Memory 6.4GB": "SVR.VRDS.STAND.C004.M006.NET.SSD.B050.G002"
}
```

## Argument Reference

The following arguments are supported:

* `redis_image_product_code` - (Required) You can get one from `data.ncloud_redis_image_products`. This is a required value, and each available Redis's specification varies depending on the Redis image product.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `product_list` - List of Redis product.
  * `product_name` - Product name.
  * `product_code` - Product code.
  * `product_type` - Product type code.
  * `product_description` - Product description.
  * `infra_resource_type` - Infra resource type code.
  * `cpu_count` - CPU count.
  * `memory_size` - Memory size.
  * `disk_type` - Disk type.
