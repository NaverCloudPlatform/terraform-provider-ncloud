---
subcategory: "Server"
---


# Resource: ncloud_public_ip

Provides a Public IP instance resource.

## Example Usage

```hcl
resource "ncloud_public_ip" "public_ip" {
  server_instance_no = "812345"
}
```

## Argument Reference

The following arguments are supported:

* `server_instance_no` - (Optional) Server instance number to assign after creating a public IP. You can get one by calling getPublicIpTargetServerInstanceList.
* `description` - (Optional) Public IP description.

~> **NOTE:** Below arguments only support Classic environment.

* `zone` - (Optional) Zone code. You can decide a zone where servers are created. You can decide in which zone the product list will be requested. default : Select the first Zone in the specific region
    Get available values using the data source `ncloud_zones`.

## Attributes Reference

* `id` - The ID of Public IP.
* `public_ip_no` - The ID of Public IP. (It is the same result as `id`)
* `public_ip` - Public IP Address.
* `kind_type` - Public IP kind type
