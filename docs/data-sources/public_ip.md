---
subcategory: "Server"
---


# Data Source: ncloud_public_ip

Get public IP instance.

## Example Usage

```hcl
variable "public_ip_no" {}
 
data "ncloud_public_ip" "public_ip" {
  id = var.public_ip_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific Public IP instance to retrieve.
* `is_associated` - (Optional) Indicates whether the public IP address is associated or not.

~> **NOTE:** Below arguments only support Classic environment.

* `zone` - (Optional) Zone code. You can filter the list of public IP instances by zone. All the public IP addresses in the zone of the region will be selected if the filter is not specified.
    Get available values using the data source `ncloud_zones`.

## Attributes Reference

* `public_ip_no` - The ID of Public IP. (It is the same result as `id`)
* `description` - Public IP description
* `kind_type` - Public IP kind type
* `server_instance_no` - Associated server instance number
* `server_name` - Associated server name