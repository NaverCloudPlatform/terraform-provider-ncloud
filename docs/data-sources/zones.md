---
subcategory: "Meta Data Sources"
---


# Data Source: ncloud_zones

Gets a list of available zones.

## Example Usage

```hcl
data "ncloud_zones" "zones" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Default: KR region.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.

## Attributes Reference

* `zones` - A List of region
    * `zone_no` - Zone number
    * `zone_code` - Zone code
    * `zone_name` - Zone name
    * `zone_description` - Zone description
    * `region_no` - Region number
