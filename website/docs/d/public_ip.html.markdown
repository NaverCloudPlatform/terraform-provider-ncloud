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
}
```

## Argument Reference

The following arguments are supported:

* `internet_line_type` - (Optional) Internet line type code. `PUBLC` (Public), `GLBL` (Global)
* `is_associated` - (Optional) Indicates whether the public IP address is associated or not.
* `instance_no_list` - (Optional) List of public IP instance numbers to get.
* `list` - (Optional) List of public IP addresses to get.
* `search_filter_name` - (Optional) `publicIp` (Public IP) | `associatedServerName` (Associated server name)
* `search_filter_value` - (Optional) Filter value to search
* `region` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Default: KR region.
* `zone` - (Optional) Zone code. You can filter the list of public IP instances by zones. All the public IP addresses in the zone of the region will be selected if the filter is not specified.
    Get available values using the data source `ncloud_zones`.
* `sorted_by` - (Optional) The column based on which you want to sort the list.
* `sorting_order` - (Optional) Sorting order of the list. `ascending` (Ascending) | `descending` (Descending) [case insensitive]. Default: `ascending` Ascending

## Attributes Reference

* `instance_no` - Public IP instance number
* `public_ip` - Public IP
* `description` - Public IP description
* `create_date` - Creation date of the public ip
* `internet_line_type` - Internet line type
* `instance_status_name` - Public IP instance status name
* `instance_status` - Public IP instance status
* `instance_operation` - Public IP instance operation
* `kind_type` - Public IP kind type
* `server_instance` - Associated server instance
    * `server_instance_no` - Associated server instance number
    * `server_name` - Associated server name
    * `create_date` - Creation date of the server instance
