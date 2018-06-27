---
layout: "ncloud"
page_title: "NCLOUD: ncloud_nas_volumes"
sidebar_current: "docs-ncloud-datasource-nas-volumes"
description: |-
  Get NAS volume instance list
---

# Data Source: ncloud_nas_volumes

Gets a list of NAS volume instances.

## Example Usage

```hcl
data "ncloud_nas_volumes" "nas_volumes" {}
```

## Argument Reference

The following arguments are supported:

* `volume_allotment_protocol_type_code` - (Optional) Volume allotment protocol type code. All volume instances will be selected if the filter is not specified. (`NFS` | `CIFS`)
* `is_event_configuration` - (Optional) Indicates whether the event is set. All volume instances will be selected if the filter is not specified. (`true` | `false`)
* `is_snapshot_configuration` - (Optional) Indicates whether a snapshot volume is set. All volume instances will be selected if the filter is not specified. (`true` | `false`)
* `nas_volume_instance_no_list` - (Optional) List of nas volume instance numbers.
* `region_no` - (Optional) Region number. Get available values using the `data ncloud_regions`.
* `zone_no` - (Optional) Zone number. Get available values using the `data ncloud_regions`.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.

## Attributes Reference

* `nas_volumes` - NAS Volume Instance list
    * `nas_volume_instance_no` - NAS volume instance number
    * `volume_name` - Volume name
    * `nas_volume_instance_status` - NAS volume instance status
        * `code` - NAS volume instance status code
        * `code_name` - NAS volume instance status name
    * `create_date` - Creation date of the NAS Volume instance
    * `volume_allotment_protocol_type` - Volume allotment protocol type
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
    * `nas_volume_description` - NAS volume description
    * `zone` - Zone info
        * `zone_no` - Zone number
        * `zone_code` - Zone code
        * `zone_name` - Zone name
    * `region` - Region info
        * `region_no` - Region number
        * `region_code` - Region code
        * `region_name` - Region name
