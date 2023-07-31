---
subcategory: "Auto Scaling"
---


# Resource: ncloud_auto_scaling_policy

Provides a ncloud auto scaling policy resource.

## Example Usage
### Classic environment
```hcl
resource "ncloud_launch_configuration" "lc" {
  name = "my-lc"
  server_image_product_code = "SPSW0LINUX000046"
  server_product_code = "SPSVRSSD00000003"
}

resource "ncloud_auto_scaling_group" "asg" {
  name = "my-auto"
  launch_configuration_no = ncloud_launch_configuration.lc.launch_configuration_no
  min_size = 1
  max_size = 1
  zone_no_list = ["2"]
  wait_for_capacity_timeout = "0"
}

resource "ncloud_auto_scaling_policy" "policy" {
  name = "my-policy"
  adjustment_type_code = "CHANG"
  scaling_adjustment = 2
  auto_scaling_group_no = ncloud_auto_scaling_group.asg.auto_scaling_group_no
}
```
### VPC environment
```hcl
resource "ncloud_launch_configuration" "lc" {
  name = "my-lc"
  server_image_product_code = "SW.VSVR.OS.LNX64.CNTOS.0703.B050"
  server_product_code = "SVR.VSVR.HICPU.C002.M004.NET.SSD.B050.G002"
}

resource "ncloud_vpc" "example" {
  ipv4_cidr_block    = "10.0.0.0/16"
}

resource "ncloud_subnet" "example" {
  vpc_no             = ncloud_vpc.example.vpc_no
  subnet             = "10.0.0.0/24"
  zone               = "KR-2"
  network_acl_no     = ncloud_vpc.example.default_network_acl_no
  subnet_type        = "PUBLIC"
  usage_type         = "GEN"
}

resource "ncloud_auto_scaling_group" "auto" {
  access_control_group_no_list = [ncloud_vpc.example.default_access_control_group_no]
  subnet_no = ncloud_subnet.example.subnet_no
  launch_configuration_no = ncloud_launch_configuration.lc.launch_configuration_no
  min_size = 1
  max_size = 1
}

resource "ncloud_auto_scaling_policy" "test-policy-CHANG" {
  name = "tf-policy"
  adjustment_type_code = "CHANG"
  scaling_adjustment = 2
  auto_scaling_group_no = ncloud_auto_scaling_group.auto.auto_scaling_group_no
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Auto Scaling Policy name to create. Only lowercase alphanumeric characters and non-consecutive hyphens (-) allowed. First character must be a letter, but the last character may be a letter or a number.
* `adjustment_type_code` - (Required) Determines how the number of servers is scaled when the scaling policy is performed. Valid values are `CHANG`, `EXACT`, and `PRCNT`.
* `scaling_adjustment` - (Required) Specify the adjustment value for the adjustment type. Enter a negative value to decrease when adjustTypeCode is `CHANG` or `PRCNT`.
* `cooldown` - (Optional) The cooldown time is the period set to ignore even if the monitoring event alarm occurs after the actual scaling is being performed or is completed.
* `min_adjustment_step` - (Optional) Change the number of server instances by the minimum adjustment width.
* `auto_scaling_group_no` - (Required) The number of the auto scaling group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of Auto Scaling Policy.
* `name` - The ID of Auto Scaling Policy (It is the same result as id).
