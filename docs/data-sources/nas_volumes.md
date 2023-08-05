---
subcategory: "NAS Volume"
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
* `no_list` - (Optional) List of nas volume instance numbers.
* `region` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Default: KR region.
* `zone` - (Optional) Zone code. Get available values using the data source `ncloud_zones`.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.
  
## Attributes Reference

* `ids` - A list of NAS Volume ID.
* `nas_volumes` - A list of NAS Volume Instance.
