---
layout: "ncloud"
page_title: "NCLOUD: ncloud_nas_volume"
sidebar_current: "docs-ncloud-resource-nas-volume"
description: |-
  Provides a ncloud NAS volume.
---

# ncloud_nas_volume

Provides a ncloud NAS volume.

## Example Usage

```hcl
resource "ncloud_nas_volume" "test" {
	"volume_name_postfix" = "vol"
	"volume_size_gb" = "600"
	"volume_allotment_protocol_type_code" = "NFS"
}
```

## Argument Reference

The following arguments are supported:

* `volume_name_postfix` - (Required) Name of a NAS volume to create. Enter a volume name that can be 3-20 characters in length after the name already entered for user identification.
* `volume_size_gb` - (Required) Enter the nas volume size to be created. You can enter in GB units.
* `volume_allotment_protocol_type_code` - (Required) Volume allotment protocol type code. `NFS` | `CIFS`
    `NFS`: You can mount the volume in a Linux server such as CentOS and Ubuntu.
    `CIFS`: You can mount the volume in a Windows server.
* `server_instance_no_list` - (Optional) List of server instance numbers for which access to NFS is to be controlled
* `custom_ip_list` - (Optional) To add a server of another account to the NAS volume, enter a private IP address of the server.
* `cifs_user_name` - (Conditional) CIFS user name. The ID must contain a combination of English alphabet and numbers, which can be 6-20 characters in length.
* `cifs_user_password` - (Conditional) CIFS user password. The password must contain a combination of at least 2 English letters, numbers and special characters, which can be 8-14 characters in length.
* `nas_volume_description` - (Optional) NAS volume description
* `region_code` - (Optional) Region code. Get available values using the data source `ncloud_regions`. Default: KR region. Conflicts with `region_no`
* `region_no` - (Optional) Region number. Get available values using the data source `ncloud_regions`. Default: KR region. Conflicts with `region_code`
* `zone_code` - (Optional) Zone code. Zone in which you want to create a NAS volume. Default: The first zone of the region.
    Get available values using the data source `ncloud_zones`. Conflicts with `zone_no`
* `zone_no` - (Optional) Zone number. Zone in which you want to create a NAS volume. Default: The first zone of the region.
    Get available values using the data source `ncloud_zones`. Conflicts with `zone_code`

## Attributes Reference

* `volume_name` - NAS volume name.
* `nas_volume_instance_status` - NAS Volume instance status
    * `code` - NAS Volume instance status code
    * `code_name` - NAS Volume instance status name
* `create_date` - Creation date of the NAS volume
* `volume_allotment_protocol_type` - Volume allotment protocol type.
    * `code` - Volume allotment protocol type code
    * `code_name` - Volume allotment protocol type name
* `volume_total_size` - Volume total size
* `volume_size` - Volume size
* `volume_use_size` - Volume use size
* `volume_use_ratio` - Volume use ratio
* `snapshot_volume_size` - Snapshot volume size
* `snapshot_volume_use_size` - Snapshot volume use size
* `snapshot_volume_use_ratio` - Snapshot volume use ratio
* `is_snapshot_configuration` - Indicates whether a snapshot volume is set.
* `is_event_configuration` - Indicates whether the event is set.
* `nas_volume_instance_custom_ip_list` - NAS volume instance custom IP list
* `zone` - Zone info
    * `zone_no` - Zone number
    * `zone_code` - Zone code
    * `zone_name` - Zone name
* `region` - Region info
    * `region_no` - Region number
    * `region_code` - Region code
    * `region_name` - Region name