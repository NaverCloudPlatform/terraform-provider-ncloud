---
layout: "ncloud"
page_title: "NCLOUD: ncloud_public_ip"
sidebar_current: "docs-ncloud-resource-public-ip"
description: |-
  Provides an ncloud public IP instance resource.
---

# ncloud_public_ip

Provides an ncloud public IP instance resource.

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
* `public_ip_description` - (Optional) Public IP description.
* `internet_line_type_code` - (Optional) Internet line code. PUBLC(Public), GLBL(Global)
* `region_no` - (Optional) You can reach a state in which inout is possible by calling `data ncloud_regions`.
* `zone_no` - (Optional) You can decide a zone where servers are created. You can decide which zone the product list will be requested at.
  You can get one by calling `data ncloud_zones`.
  default : Select the first Zone in the specific region

## Attributes Reference

* `public_ip_instance_no` - Public IP instance No.
* `public_ip` - Public IP Address.
* `public_ip_description` - Public IP description.
* `create_date` - Creation date of the public IP instance
* `internet_line_type` - Internet line type
    * `code` - Internet line type code
    * `code_name` - Internet line type code name
* `public_ip_instance_status_name` - Public IP instance status name
* `public_ip_instance_status` - Public IP instance status
    * `code` - Public IP instance status code
    * `code_name` - Public IP instance status code name
* `public_ip_instance_operation` - Public IP instance operation
    * `code` - Public IP instance operation code
    * `code_name` - Public IP instance operation code name
* `public_ip_kind_type` - Public IP kind type
* `zone` - Zone info
    * `zone_no` - Zone number
    * `zone_code` - Zone code
    * `zone_name` - Zone name
* `region` - Region info
    * `region_no` - Region number
    * `region_code` - Region code
    * `region_name` - Region name
