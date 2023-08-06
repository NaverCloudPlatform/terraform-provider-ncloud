---
subcategory: "Server"
---


# Data Source: ncloud_port_forwarding_rules

Gets a list of port forwarding rules.
When a server is created for the first time, a public IP address for port forwarding is given per account.

## Example Usage

```hcl
data "ncloud_port_forwarding_rules" "rules" {
  zone_code = "KR-1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) Region code. Get available values using the data source `ncloud_regions`.
    Default: KR region.
* `zone` - (Optional) Zone code. You can decide a zone where servers are created. You can decide which zone the product list will be requested in. Default : Select the first Zone in the specific region
    Get available values using the data source `ncloud_zones`.
* `port_forwarding_internal_port` - (Optional) Port forwarding internal port.
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.

## Attributes Reference

* `port_forwarding_configuration_no` - Port forwarding configuration number
* `port_forwarding_public_ip` - Port forwarding public ip
* `port_forwarding_rule_list` - Port forwarding rule list
    * `server_instance_no` - Server instance number
    * `port_forwarding_external_port` - Port forwarding external port.
    * `port_forwarding_internal_port` - Port forwarding internal port.
