# Data Source: ncloud_nks_server_products

Provides list of available Kubernetes Nodepool ServerProducts.

## Example Usage

```hcl
data "ncloud_nks_server_products" "products" {}


data "ncloud_nks_server_images" "images"{
  filter {
    name = "label"
    values = ["ubuntu-20.04-64-server"]
  }
}

data "ncloud_nks_server_products" "product" {

  software_code = data.ncloud_nks_server_images.images.images[0].value
  zone = "KR-1"

  filter {
    name = "product_type"
    values = [ "STAND" ]
  }

  filter {
    name = "cpu_count"
    values = [ 2 ]
  }

  filter {
    name = "memory_size"
    values = [ "8GB" ]
  }

  filter {
    name = "product_code"
    values = [ "SSD" ]
    regex = true
  }
}

```

## Argument Reference

The following arguments are supported:

* `software_code` - (Required) NKS ServerImage code.
* `zone` - (Required) zone Code.

* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `products` - A list of ServerProduct
  * `label` - ServerProduct spec korean description
  * `value` - ServerProduct code
  * `detail`
    * `cpu_count` - Number of cpu
    * `gpu_count` - Number of gpu
    * `gpu_memory_size` - Size of GPU memory(GB)
    * `memory_size` - Size of memory(GB)
    * `product_code` -  ServerProduct code
    * `product_english_desc` - ServerProduct spec english description
    * `product_korean_desc` - ServerProduct spec korean description
    * `product_type` - ServerProduct Type