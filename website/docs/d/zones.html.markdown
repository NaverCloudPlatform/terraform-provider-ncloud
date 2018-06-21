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

* `region_no` - (Optional) region no for filtering
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.

## Attributes Reference

* `zones` - A List of region
    * `zone_no` - zone no
    * `zone_code` - zone code
    * `zone_name` - zone name
    * `zone_description` - zone description
    * `region_no` - region no