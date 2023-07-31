---
subcategory: "VPC"
---


# Resource: ncloud_network_interface

Provides a Network Interface resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage

```hcl
resource "ncloud_vpc" "vpc" {
	ipv4_cidr_block    = "10.0.0.0/16"
}

resource "ncloud_subnet" "subnet" {
	vpc_no             = ncloud_vpc.vpc.vpc_no
	subnet             = cidrsubnet(ncloud_vpc.vpc.ipv4_cidr_block, 8, 1) // 10.0.1.0/24
	zone               = "KR-2"
	network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
	subnet_type        = "PUBLIC"
	usage_type         = "GEN"
}

resource "ncloud_network_interface" "nic" {
	name                  = "my-nic"
	description           = "for example"
	subnet_no             = ncloud_subnet.subnet.id
	private_ip            = "10.0.1.6"
	access_control_groups = [ncloud_vpc.vpc.default_access_control_group_no]
}
```

## Argument Reference

The following arguments are supported:

* `subnet_no` - (Required) The ID of the associated Subnet.
* `access_control_groups` - (Required) List of ACG ID to apply to network interfaces. A maximum of three ACGs can be
  applied.
* `name` - (Optional) The name to create. If omitted, Terraform will assign a random, unique name.
* `description` - (Optional) description to create.
* `private_ip` - (Optional) Set the IP addresses that you want to assign to the network interface. Must be in the IP
  address range of the subnet where the network interface is created. The last `0` to `5' IP address of the Subnet is
  not available and duplicate IP addresses are not available at the Subnet scope.
* `server_instance_no` - (Optional) The ID of server instance to assign network interface.

## Attributes Reference

* `id` - The ID of Network Interface.
* `network_interface_no` - The ID of Network Interface. (It is the same result as `id`)
* `status` - The status of Network Interface.
* `instance_type` - Type of server instance.
* `is_default` - Whether is default or not by Server instance creation.