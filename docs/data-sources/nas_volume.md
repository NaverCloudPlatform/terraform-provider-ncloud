---
subcategory: "NAS Volume"
---


# Data Source: ncloud_nas_volume

Get NAS volume instance

## Example Usage

```hcl
variable "nas_volume_no" {}

data "ncloud_nas_volume" "vol" {
  id = var.nas_volume_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of specific NAS Volume to retrieve.
* `volume_allotment_protocol_type_code` - (Optional) Volume allotment protocol type code. All volume instances will be selected if the filter is not specified. (`NFS` | `CIFS`).
* `is_event_configuration` - (Optional) Indicates whether the event is set. All volume instances will be selected if the filter is not specified. (`true` | `false`).
* `is_snapshot_configuration` - (Optional) Indicates whether a snapshot volume is set. All volume instances will be selected if the filter is not specified. (`true` | `false`).
* `zone` - (Optional) Zone code. Get available values using the data source `ncloud_zones`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.
  
## Attributes Reference

* `nas_volume_no` - The ID of NAS Volume.
* `name` - Volume name.
* `volume_total_size` - Volume total size.
* `volume_size` - Volume size.
* `snapshot_volume_size` - Snapshot volume size.
* `custom_ip_list` - NAS volume instance custom IP list.
* `description` - NAS volume description.
* `is_encrypted_volume` - Volume encryption. 
* `mount_information` - Mount information for NAS volume.