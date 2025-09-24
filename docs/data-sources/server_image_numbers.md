---
subcategory: "Server"
---


# Data Source: ncloud_server_image_numbers

To create a server instance (VM), you should select a server image. This data source gets a list of server images. You can look up the image numbers of Gen2(XEN) and Gen3(KVM) servers.

~> **NOTE:** This only supports VPC environment.

## Example Usage

The following example shows how to take a list of Server image.

```terraform
data "ncloud_server_image_numbers" "example" {
  output_file = "image.json" 
}

output "image_list" {
  value = {
    for image in data.ncloud_server_image_numbers.example.image_number_list:
    image.server_image_number => [image.name, image.hypervisor_type]
  }
}
```

Outputs: 
```terraform
image_list = {
  "16187005" = [
    "ubuntu-20.04",
    "XEN",
  ]
  "16187007" = [
    "win-2019-64-en",
    "XEN",
  ]
  "17552318" = [
    "ubuntu-20.04-base",
    "KVM",
  ]
  "23789321" = [
    "ubuntu-22.04-gpu",
    "KVM",
  ]
  "25495367" = [
    "rocky-8.10-base",
    "KVM",
  ]
  "25623982" = [
    "rocky-8.10-gpu",
    "KVM",
  ]
  "25624115" = [
    "rocky-8.10-base",
    "XEN",
  ]
}
```

```terraform
data "ncloud_server_image_numbers" "example-xen" {
  server_image_name = "rocky-8.10-base"
  filter {
    name = "hypervisor_type"
    values = ["XEN"]
  }
}

data "ncloud_server_image_numbers" "example-kvm" {
  server_image_name = "rocky-8.10-base"
  filter {
    name = "hypervisor_type"
    values = ["KVM"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `server_image_name` - (Optional) Server image name.
* `hypervisor_type` - (Optional) Server image hypervisor type. Options: `XEN` | `KVM`
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `image_number_list` - List of SEVER image number.
  * `server_image_number` - Server image number.
  * `name` - Server image name.
  * `description` - Server image description.
  * `type` - Server image type. (`SELF` or `NCP`)
  * `hypervisor_type` - Server image hypervisor type. (`XEN` or `KVM`)
  * `cpu_architecture_type` - Server image cpu type.
  * `os_category_type` - Server image os category type. (`LINUX` or `WINDOWS`)
  * `os_type` - Server image os type. (`ROCKY` or `UBUNTU` or `CENTOS` or `WINDOWS`)
  * `product_code` - The code of image product.
  * `block_storage_mapping_list` - List of block storage allocated to the server image. Viewable after the server image is created.
    * `order` - Block storage order.
    * `block_storage_snapshot_instance_no` - Block storage snapshot instance number.
    * `block_storage_snapshot_name` - Block storage snapshot name.
    * `block_storage_size` - Block storage size (byte).
    * `block_storage_name` - Block storage name.
    * `block_storage_volume_type` - Block storage volume type. (`SSD`|`HDD`|`CB1`|`FB1`|`CB2`|`FB2`).
    * `iops` - IOPS.
    * `throughput` - Load balancing performance.
    * `is_encrypted_volume` - Volume encryption status. (`true` or `false`)
