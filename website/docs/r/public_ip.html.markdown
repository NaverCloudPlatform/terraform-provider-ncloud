---
layout: "ncloud"
page_title: "NCLOUD: ncloud_public_ip"
sidebar_current: "docs-ncloud-resource-public-ip"
description: |-
  Provides a ncloud public IP instance resource.
---

# ncloud_public_ip

Provides a ncloud public IP instance resource.

## Example Usage

```hcl
resource "ncloud_public_ip" "public_ip" {
  "server_instance_no" = "812345"
  "region_no"          = "1"
  "zone_no"            = "3"
}
```

## Argument Reference

The following arguments are supported:

* `server_instance_no` - (Optional) Server instance No. to assign after creating a public IP. You can get one by calling getPublicIpTargetServerInstanceList.
* `description` - (Optional) Public IP description.
* `internet_line_type_code` - (Optional) Internet line code. PUBLC(Public), GLBL(Global)
* `region_code` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Conflicts with `region_no`. Only one of `region_no` and `region_code` can be used.
    Default: KR region.
* `region_no` - (Optional) Region number. Get available values using the data source `ncloud_regions`.
    Conflicts with `region_code`. Only one of `region_no` and `region_code` can be used.
    Default: KR region.
* `zone_code` - (Optional) Zone code. You can decide a zone where servers are created. You can decide which zone the product list will be requested at. default : Select the first Zone in the specific region
    Get available values using the data source `ncloud_zones`.
    Conflicts with `zone_no`. Only one of `zone_no` and `zone_code` can be used.
* `zone_no` - (Optional) Zone number. You can decide a zone where servers are created. You can decide which zone the product list will be requested at. default : Select the first Zone in the specific region
    Get available values using the data source `ncloud_zones`.
    Conflicts with `zone_code`. Only one of `zone_no` and `zone_code` can be used.

## Attributes Reference

* `instance_no` - Public IP instance No.
* `public_ip` - Public IP Address.
* `description` - Public IP description.
* `create_date` - Creation date of the public IP instance
* `internet_line_type` - Internet line type
    * `code` - Internet line type code
    * `code_name` - Internet line type code name
* `instance_status_name` - Public IP instance status name
* `instance_status` - Public IP instance status
    * `code` - Public IP instance status code
    * `code_name` - Public IP instance status code name
* `instance_operation` - Public IP instance operation
    * `code` - Public IP instance operation code
    * `code_name` - Public IP instance operation code name
* `kind_type` - Public IP kind type
* `zone` - Zone info
    * `zone_no` - Zone number
    * `zone_code` - Zone code
    * `zone_name` - Zone name
* `region` - Region info
    * `region_no` - Region number
    * `region_code` - Region code
    * `region_name` - Region name
