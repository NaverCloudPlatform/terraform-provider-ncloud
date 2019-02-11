---
layout: "ncloud"
page_title: "NCLOUD: ncloud_nas_volume"
sidebar_current: "docs-ncloud-datasource-nas-volume"
description: |-
  Get NAS volume instance
---

# Data Source: ncloud_nas_volume

Get NAS volume instance

## Example Usage

```hcl
data "ncloud_nas_volume" "vol" {}
```

## Argument Reference

The following arguments are supported:

* `volume_allotment_protocol_type_code` - (Optional) Volume allotment protocol type code. All volume instances will be selected if the filter is not specified. (`NFS` | `CIFS`)
* `is_event_configuration` - (Optional) Indicates whether the event is set. All volume instances will be selected if the filter is not specified. (`true` | `false`)
* `is_snapshot_configuration` - (Optional) Indicates whether a snapshot volume is set. All volume instances will be selected if the filter is not specified. (`true` | `false`)
* `no_list` - (Optional) List of nas volume instance numbers.
* `region` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Default: `KR` region.
* `zone` - (Optional) Zone code. Get available values using the data source `ncloud_zones`.

## Attributes Reference

* `instance_no` - NAS volume instance number
* `volume_name` - Volume name
* `instance_status` - NAS volume instance status code
* `create_date` - Creation date of the NAS Volume instance
* `volume_total_size` - Volume total size
* `volume_size` - Volume size
* `volume_use_size` - Volume use size
* `volume_use_ratio` - Volume use ratio
* `snapshot_volume_size` - Snapshot volume size
* `snapshot_volume_use_size` - Snapshot volume use size
* `snapshot_volume_use_ratio` - Snapshot volume use ratio
* `instance_custom_ip_list` - NAS volume instance custom IP list
* `description` - NAS volume description
