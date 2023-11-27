---
subcategory: "Server"
---


# ncloud_port_forwarding_rule

Provides a ncloud port forwarding rule resource.

~> **NOTE:** This resource only supports Classic environment.

## Example Usage

```hcl
resource "ncloud_port_forwarding_rule" "rule" {
   port_forwarding_configuration_no = "1222"
   server_instance_no = "812345"
   port_forwarding_external_port = "2022"
   port_forwarding_internal_port = "22"
}
```

## Argument Reference

The following arguments are supported:

* `server_instance_no` - (Required) Server instance number for which port forwarding is set
* `port_forwarding_external_port` - (Required) External port for port forwarding
* `port_forwarding_internal_port` - (Required) Internal port for port forwarding. Only the following ports are available. [Linux: `22` | Windows: `3389`]
* `port_forwarding_configuration_no` - (Optional) Port forwarding configuration number. You can get by calling `data ncloud_port_forwarding_rules`

## Attributes Reference

* `id` - The ID of Port forwarding rule.`port_forwarding_configuration_no:zone_no:port_forwarding_external_port` 
* `port_forwarding_public_ip` - Port forwarding Public IP
* `zone` - Zone code

## Import

### `terraform import` command

* Port Forwarding Rule can be imported using the `id`. For example:

```console
$ terraform import ncloud_port_forwarding_rule.rsc_name 1222:3:2022
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Port Forwarding Rule using the `id`. For example:

```terraform
import {
  to = ncloud_port_forwarding_rule.rsc_name
  id = "1222:3:2022"
}
```
