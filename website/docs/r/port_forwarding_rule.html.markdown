---
layout: "ncloud"
page_title: "NCLOUD: ncloud_port_forwarding_rule"
sidebar_current: "docs-ncloud-resource-port-forwarding-rule"
description: |-
  Provides a ncloud port forwarding rule resource.
---

# ncloud_port_forwarding_rule

Provides a ncloud port forwarding rule resource.

## Example Usage

```hcl
resource "ncloud_port_forwarding_rule" "rule" {
   "port_forwarding_configuration_no" = "1222"
   "server_instance_no" = "812345"
   "port_forwarding_external_port" = "2022"
   "port_forwarding_internal_port" = "22"
}
```

## Argument Reference

The following arguments are supported:

* `port_forwarding_configuration_no` - (Optional) Port forwarding configuration number. You can get by calling `data ncloud_port_forwarding_rules`
* `server_instance_no` - (Required) Server instance number for which port forwarding is set
* `port_forwarding_external_port` - (Required) External port for port forwarding
* `port_forwarding_internal_port` - (Required) Internal port for port forwarding. Only the following ports are available. [Linux: `22` | Windows: `3389`]

## Attributes Reference

* `port_forwarding_public_ip` - Port forwarding Public IP