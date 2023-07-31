---
subcategory: "Meta Data Sources"
---


# Data Source: ncloud_regions

Gets a list of available regions.

## Example Usage

```hcl
data "ncloud_regions" "regions" {}
```

## Argument Reference

The following arguments are supported:

* `code` - (Optional) region code for filtering
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.

## Attributes Reference

* `regions` - A List of region
    * `region_no` - Region number
    * `region_code` - Region code
    * `region_name` - Region name