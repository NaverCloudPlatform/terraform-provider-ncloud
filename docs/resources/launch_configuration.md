---
subcategory: "Auto Scaling"
---


# Resource: ncloud_launch_configuration

Provides a ncloud launch configuration resource.

## Example Usage
```hcl
resource "ncloud_launch_configuration" "lc" {
  name = "my-lc"
  server_image_product_code = "SW.VSVR.OS.LNX64.ROCKY.0810.B050"
  server_product_code = "SVR.VSVR.HICPU.C002.M004.NET.SSD.B050.G002"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Launch Configuration name to create. default : Ncloud assigns default values.
* `server_image_product_code` - (Optional) Server image product code to determine which server image to create. It can be obtained through data ncloud_server_images. You are required to select one between two parameters: server image product code (server_image_product_code) and member server image number member_server_image_no) 
* `server_product_code` - (Optional) Server product code to determine the server specification to create. It can be obtained through the getServerProductList action. Default : Selected as minimum specification. The minimum standards are 1. memory 2. CPU 3. basic block storage size 4. disk type (NET,LOCAL)
* `member_server_image_no` - (Optional) Required value when creating a server from a manually created server image. It can be obtained through the getMemberServerImageList action.
* `login_key_name` - (Optional) The login key name to encrypt with the public key. Default : Uses the login key name most recently created.
* `init_script_no` - (Optional) Set init script ID, The server can run a user-set initialization script at first boot.
* `is_encrypted_volume` - (Optional) you can set whether to encrypt basic block storage if server image is RHV. Default false.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of Launch Configuration.
* `launch_configuration_no` - The ID of Launch Configuration (It is the same result as id)

## Import

### `terraform import` command

* Launch Configuration can be imported using the `id`. For example:

```console
$ terraform import ncloud_launch_configuration.rsc_name 12345
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Launch Configuration using the `id`. For example:

```terraform
import {
  to = ncloud_launch_configuration.rsc_name
  id = "12345"
}
```
