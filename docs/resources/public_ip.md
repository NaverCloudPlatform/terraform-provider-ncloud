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
  server_instance_no = "812345"
}
```

## Argument Reference

The following arguments are supported:

* `server_instance_no` - (Optional) Server instance No. to assign after creating a public IP. You can get one by calling getPublicIpTargetServerInstanceList.
* `description` - (Optional) Public IP description.
* `internet_line_type` - (Optional) Internet line code. PUBLC(Public), GLBL(Global)
* `zone` - (Optional) Zone code. You can decide a zone where servers are created. You can decide which zone the product list will be requested at. default : Select the first Zone in the specific region
    Get available values using the data source `ncloud_zones`.

## Attributes Reference

* `instance_no` - Public IP instance No.
* `public_ip` - Public IP Address.
* `create_date` - Creation date of the public IP instance
* `instance_status_name` - Public IP instance status name
* `instance_status` - Public IP instance status code
* `instance_operation` - Public IP instance operation code
* `kind_type` - Public IP kind type
