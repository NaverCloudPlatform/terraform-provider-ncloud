# Data Source: ncloud_server_images

To create a server instance (VM), you should select a server image. This data source gets a list of server images.

## Example Usage

```hcl
data "ncloud_server_images" "all" {
  filter {
    name = "product_name"
    values = ["CentOS 7.3 (64-bit)"]
  }

  output_file = "server_images.json"
}
```

## Argument Reference

The following arguments are supported:

* `product_code` - (Optional) Product code you want to view on the list. Use this when searching for 1 product.
* `platform_type_code_list` - (Optional) Values required for identifying platforms in list-type.
    The available values are as follows: Linux 32Bit(LNX32) | Linux 64Bit(LNX64) | Windows 32Bit(WND32) | Windows 64Bit(WND64) | Ubuntu Desktop 64Bit(UBD64) | Ubuntu Server 64Bit(UBS64)
* `infra_resource_detail_type_code` - (Optional) infra resource detail type code.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `server_images` - A List of server image product code
