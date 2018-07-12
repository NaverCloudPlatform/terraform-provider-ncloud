---
layout: "ncloud"
page_title: "NCLOUD: ncloud_zones"
sidebar_current: "docs-ncloud-datasource-zones"
description: |-
  Get region list
---

# Data Source: ncloud_zones

Gets a list of available zones.

## Example Usage

```hcl
data "ncloud_zones" "zones" {}
```

## Argument Reference

The following arguments are supported:

* `region_code` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Conflicts with `region_no`. Only one of `region_no` and `region_code` can be used.
    Default: KR region.
* `region_no` - (Optional) Region number. Get available values using the data source `ncloud_regions`.
    Conflicts with `region_code`. Only one of `region_no` and `region_code` can be used.
    Default: KR region.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.

## Attributes Reference

* `zones` - A List of region
    * `zone_no` - Zone number
    * `zone_code` - Zone code
    * `zone_name` - Zone name
    * `zone_description` - Zone description
    * `region_no` - Region number