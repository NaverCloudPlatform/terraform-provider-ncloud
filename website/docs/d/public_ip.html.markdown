---
layout: "ncloud"
page_title: "NCLOUD: ncloud_public_ip"
sidebar_current: "docs-ncloud-datasource-public-ip"
description: |-
  Get public IP
---

# Data Source: ncloud_public_ip

Get public IP instance.


## Example Usage

```hcl
data "ncloud_public_ip" "public_ip" {
  "sorted_by" = "publicIp"
  "sorting_order" = "ascending"
  "most_recent" = "true"
}
```

## Argument Reference

The following arguments are supported:

* `most_recent` - (Optional) If more than one result is returned, get the most recent created Public IP.
* `internet_line_type_code` - (Optional) Internet line type code. `PUBLC` (Public), `GLBL` (Global)
* `is_associated` - (Optional) Indicates whether the public IP address is associated or not.
* `public_ip_instance_no_list` - (Optional) List of public IP instance numbers to get.
* `public_ip_list` - (Optional) List of public IP addresses to get.
* `search_filter_name` - (Optional) `publicIp` (Public IP) | `associatedServerName` (Associated server name)
* `search_filter_value` - (Optional) Filter value to search
* `region_code` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Conflicts with `region_no`. Only one of `region_no` and `region_code` can be used.
    Default: KR region.
* `region_no` - (Optional) Region number. Get available values using the data source `ncloud_regions`.
    Conflicts with `region_code`. Only one of `region_no` and `region_code` can be used.
    Default: KR region.
* `zone_code` - (Optional) Zone code. You can filter the list of public IP instances by zones. All the public IP addresses in the zone of the region will be selected if the filter is not specified.
    Get available values using the data source `ncloud_zones`.
    Conflicts with `zone_no`. Only one of `zone_no` and `zone_code` can be used.
* `zone_no` - (Optional) Zone number. You can filter the list of public IP instances by zones. All the public IP addresses in the zone of the region will be selected if the filter is not specified.
    Get available values using the data source `ncloud_zones`.
    Conflicts with `zone_code`. Only one of `zone_no` and `zone_code` can be used.
* `sorted_by` - (Optional) The column based on which you want to sort the list.
* `sorting_order` - (Optional) Sorting order of the list. `ascending` (Ascending) | `descending` (Descending) [case insensitive]. Default: `ascending` Ascending

## Attributes Reference

* `public_ip_instance_no` - Public IP instance number
* `public_ip` - Public IP
* `public_ip_description` - Public IP description
* `create_date` - Creation date of the public ip
* `internet_line_type` - Internet line type
* `public_ip_instance_status_name` - Public IP instance status name
* `public_ip_instance_status` - Public IP instance status
* `public_ip_instance_operation` - Public IP instance operation
* `public_ip_kind_type` - Public IP kind type
* `server_instance` - Associated server instance
    * `server_instance_no` - Associated server instance number
    * `server_name` - Associated server name
    * `create_date` - Creation date of the server instance
