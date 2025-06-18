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

## Attributes Reference

* `public_ip_no` - The ID of Public IP. (It is the same result as `id`)
* `description` - Public IP description
* `kind_type` - Public IP kind type
* `server_instance_no` - Associated server instance number
* `server_name` - Associated server name
