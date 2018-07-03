---
layout: "ncloud"
page_title: "NCLOUD: ncloud_port_forwarding_rule"
sidebar_current: "docs-ncloud-datasource-port-forwarding-rule"
description: |-
  Get port forwarding rule
---

# Data Source: ncloud_port_forwarding_rule

Get a port forwarding rule.
When a server is created for the first time, a public IP address for port forwarding is given per account.

## Example Usage

```hcl
data "ncloud_port_forwarding_rule" "test" {
  "port_forwarding_external_port" = "4088"
}
```
ncloud_nas_volume
## Argument Reference

The following arguments are supported:

* `internet_line_type_code` - (Optional) Internet line code. PUBLC(Public), GLBL(Global)
* `region_code` - (Optional) Region code. Get available values using the data source `ncloud_regions`. Conflicts with `region_no`
* `region_no` - (Optional) Region number. Get available values using the data source `ncloud_regions`. Conflicts with `region_code`
* `zone_code` - (Optional) Zone code. You can decide a zone where servers are created. You can decide which zone the product list will be requested at. default : Select the first Zone in the specific region
    Get available values using the data source `ncloud_zones`. Conflicts with `zone_no`
* `zone_no` - (Optional) Zone number. You can decide a zone where servers are created. You can decide which zone the product list will be requested at. default : Select the first Zone in the specific region
    Get available values using the data source `ncloud_zones`. Conflicts with `zone_code`
* `server_instance_no` - Filter by server instance number
* `port_forwarding_internal_port` - (Optional) Port forwarding internal port.
* `port_forwarding_external_port` - Port forwarding external port.

## Attributes Reference

* `port_forwarding_configuration_no` - Port forwarding configuration number
* `port_forwarding_public_ip` - Port forwarding public ip
* `server_instance_no` - Server instance number
* `port_forwarding_external_port` - Port forwarding external port.
* `port_forwarding_internal_port` - Port forwarding internal port.