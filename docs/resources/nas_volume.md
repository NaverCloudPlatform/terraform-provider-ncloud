---
subcategory: "NAS Volume"
---


# Resource: ncloud_nas_volume

Provides a NAS volume.

## Example Usage

```hcl
resource "ncloud_nas_volume" "test" {
	volume_name_postfix            = "vol"
	volume_size                    = "600"
	volume_allotment_protocol_type = "NFS"
}
```

## Argument Reference

The following arguments are supported:

* `volume_name_postfix` - (Required) Name of a NAS volume to create. Enter a volume name that is 3-20 characters in length after entering the name for user identification.
* `volume_size` - (Required) NAS volume size. The default capacity of a volume ranges from 500 GB to 10,000 GB. Additions can be made in units of 100 GB.
* `volume_allotment_protocol_type` - (Required) Volume allotment protocol type code. `NFS` | `CIFS`
    `NFS`: You can mount the volume in a Linux server such as CentOS and Ubuntu.
    `CIFS`: You can mount the volume in a Windows server.
* `server_instance_no_list` - (Optional) List of server instance numbers where you want to mount the NAS volume.
* `cifs_user_name` - (Optional) CIFS user name. The ID must contain a combination of English alphabet and numbers, which can be 6-19 characters in length.
* `cifs_user_password` - (Optional) CIFS user password. The password must contain a combination of at least 2 English letters, numbers and special characters,   which can be 8-14 characters in length.
* `description` - (Optional) NAS volume description. 1-1000 characters.
* `zone` - (Optional) Zone code. Zone in which you want to create a NAS volume. Default: The first zone of the region.  Get available values using the data      source `ncloud_zones`.
* `return_protection` - (Optional) Termination protection status. Default `false`

~> **NOTE:** Below arguments only support Classic environment.

* `custom_ip_list` - (Optional) To add a server of another account to the NAS volume, enter a private IP address of the server.

~> **NOTE:** Below arguments only support VPC environment.

* `is_encrypted_volume` - (Optional) Volume encryption. Default `false`.


## Attributes Reference

* `id` - The ID of NAS Volume.
* `nas_volume_no` - The ID of NAS Volume. (It is the same result as `id`)
* `name` - NAS volume name.
* `volume_total_size` - Volume total size, in GiB
* `snapshot_volume_size` - Snapshot volume size, in GiB
* `is_snapshot_configuration` - Indicates whether a snapshot volume is set.
* `is_event_configuration` - Indicates whether the event is set.
* `mount_information` - Mount information for NAS volume.
