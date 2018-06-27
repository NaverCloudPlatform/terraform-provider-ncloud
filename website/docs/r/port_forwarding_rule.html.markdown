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

```

## Argument Reference

The following arguments are supported:

* `port_forwarding_configuration_no` - (Optional) Port forwarding configuration number.
* `server_instance_no` - (Required) Server instance number for which port forwarding is set
* `port_forwarding_external_port` - (Required) External port for port forwarding
* `port_forwarding_internal_port` - (Required) Internal port for port forwarding. Only the following ports are available. [Linux: `22` | Windows: `3389`]

## Attributes Reference

* `port_forwarding_public_ip` - Port forwarding Public IP