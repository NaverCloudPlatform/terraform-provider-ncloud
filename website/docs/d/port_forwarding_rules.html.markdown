---
layout: "ncloud"
page_title: "NCLOUD: ncloud_port_forwarding_rules"
sidebar_current: "docs-ncloud-datasource-port-forwarding-rules"
description: |-
  Get port forwarding rule list
---

# Data Source: ncloud_port_forwarding_rules

Gets a list of port forwarding rules.
When a server is created for the first time, a public IP address for port forwarding is given per account.

## Example Usage

```hcl
data "ncloud_port_forwarding_rules" "rules" {}
```

## Argument Reference

The following arguments are supported:

* `internet_line_type_code` - (Optional) Internet line code. PUBLC(Public), GLBL(Global)
* `region_no` - (Optional) Region number. You can reach a state in which inout is possible by calling `data ncloud_regions`.
* `zone_no` - (Optional) Zone number. You can decide a zone where servers are created. You can decide which zone the product list will be requested at.
  You can get one by calling `data ncloud_zones`.
  default : Select the first Zone in the specific region
* `port_forwarding_internal_port` - (Optional) Port forwarding internal port.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.

## Attributes Reference

* `port_forwarding_configuration_no` - Port forwarding configuration number
* `port_forwarding_public_ip` - Port forwarding public ip
* `port_forwarding_rule_list` - Port forwarding rule list
    * `server_instance_no` - Server instance number
    * `port_forwarding_external_port` - Port forwarding external port.
    * `port_forwarding_internal_port` - Port forwarding internal port.