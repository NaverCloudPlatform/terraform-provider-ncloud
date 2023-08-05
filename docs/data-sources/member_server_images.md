---
subcategory: "Member Server Image"
---


# Data Source: ncloud_member_server_images

Gets a list of member server images.

## Example Usage

```hcl
data "ncloud_member_server_images" "member_server_images" {}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the member server image list returned by ncloud
* `no_list` - (Optional) List of member server images to view
* `platform_type_code_list` - (Optional) List of platform codes of server images to view. Linux 32Bit (LNX32) | Linux 64Bit (LNX64) | Windows 32Bit (WND32) | Windows 64Bit (WND64) | Ubuntu Desktop 64Bit (UBD64) | Ubuntu Server 64Bit (UBS64)
* `region` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Default: KR region.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.

## Attributes Reference

* `member_server_images` - A list of Member server image no
