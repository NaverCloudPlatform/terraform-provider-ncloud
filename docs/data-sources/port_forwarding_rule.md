---
subcategory: "Server"
---


# Data Source: ncloud_port_forwarding_rule

Get a port forwarding rule.
When a server is created for the first time, a public IP address for port forwarding is given per account.

## Example Usage

```hcl
data "ncloud_port_forwarding_rule" "test" {
  port_forwarding_external_port = "4088"
}
```
ncloud_nas_volume
## Argument Reference

The following arguments are supported:

* `region` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Default: KR region.
* `zone` - (Optional) Zone code. You can decide a zone where servers are created. You can decide which zone the product list will be requested in. Default : Select the first Zone in the specific region
    Get available values using the data source `ncloud_zones`.
* `server_instance_no` - Filter by server instance number
* `port_forwarding_internal_port` - (Optional) Port forwarding internal port.
* `port_forwarding_external_port` - Port forwarding external port.

## Attributes Reference

* `port_forwarding_configuration_no` - Port forwarding configuration number
* `port_forwarding_public_ip` - Port forwarding public ip
