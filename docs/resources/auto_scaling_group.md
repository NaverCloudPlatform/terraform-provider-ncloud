# Resource: ncloud_auto_scaling_group

Provides a ncloud auto scaling group resource.

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Auto Scaling Group name to create. default : Ncloud assigns default values.
* `launch_configuration_no` - (Required) Launch Configuration Number for creating Auto Scaling Group.
* `desired_capacity` - (Optional) The number of servers is adjusted according to the desired capacity value.
This value must be more than the minimum capacity and less than the maximum capacity. If this value is not specified, initially create a server with a minimum capacity. valid from `0` to `30`.
* `min_size` - (Required) The minimum size of the Auto Scaling Group. valid from `0` to `30`.
* `max_size` - (Required) The maximum size of the Auto Scaling Group. valid from `0` to `30`.
* `ignore_capacity_changes` - (Optional) If set this `true`, any changes of `desired_capacity`, `min_size` and `max_size` will be ignored and not rolled back.
* `default_cooldown` - (Optional) The cooldown time is the period set to ignore even if the monitoring event alarm occurs after the actual scaling is being performed or is completed.
* `health_check_type_code` - (Optional) `SVR` or `LOADB`. Controls how health checking is done.
* `wait_for_capacity_timeout` - (Optional) The maximum amount of time Terraform should wait for an ASG instance to become healthy. Setting this to "0" causes Terraform to skip all Capacity Waiting behavior.
* `health_check_grace_period` - (Optional) Set the time to hold health check after the server instance is put into the service with the health check hold period.

~> **NOTE:** If the `health_check_type_code` is `LOADB`, `health_check_grace_period` is required.


~> **NOTE:** Below arguments only support Classic environment.

* `zone_no_list` - (Required) the list of zone numbers where server instances belonging to this group will exist.

~> **NOTE:** Below arguments only support VPC environment.

* `subnet_no` - (Required) The ID of the associated Subnet.
* `access_control_group_no_list` - (Required) The ID of the ACG.
* `target_group_list` - (Optional) - Target Group number list of Load Balancer.

~> **NOTE:** `target_group_list` is valid only if the `health_check_type_code` is `LOADB`.

* `server_name_prefix` - (Optional) Create name beginning with the specified prefix.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of Auto Scaling Group.
* `auto_scaling_group_no` - The ID of Auto Scaling Group (It is the same result as id)
* `server_instance_no_list` - List of server instances belonging to Auto Scaling Group.

~> **NOTE:** Below attributes only support VPC environment.

* `vpc_no` - The ID of the associated VPC.
