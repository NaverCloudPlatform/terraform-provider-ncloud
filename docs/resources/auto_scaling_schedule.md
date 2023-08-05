---
subcategory: "Auto Scaling"
---


# Resource: ncloud_auto_scaling_schedule

Provides a ncloud auto scaling schedule resource.

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

resource "ncloud_auto_scaling_schedule" "schedule" {
  name = "my-schedule"
  min_size = 1
  max_size = 1
  desired_capacity = 1
  start_time = "" # 2021-02-02T15:00:00+0900
  end_time = "" # 2021-02-02T17:00:00+0900
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

resource "ncloud_auto_scaling_schedule" "schedule" {
  name = "tf-schedule"
  min_size = 1
  max_size = 1
  desired_capacity = 1
  start_time = "" # 2021-02-02T15:00:00+0900
  end_time = "" # 2021-02-02T17:00:00+0900
  auto_scaling_group_no = ncloud_auto_scaling_group.auto.auto_scaling_group_no
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Auto Scaling Schedule name to create.
* `desired_capacity` - (Required) The number of servers is adjusted according to the desired capacity value. Valid from `0` to `30`.
* `min_size` - (Required) The minimum size of the Auto Scaling Group. Valid from `0` to `30`.
* `max_size` - (Required) The maximum size of the Auto Scaling Group. Valid from `0` to `30`.
* `start_time` - (Optional) You can determine the date and time when the schedule first starts. If you don't enter `recurrence`, be sure to enter startTime. It cannot be duplicated with the startTime of another schedule and must be later than the current time, before endTime. Format : `yyyy-MM-ddTHH:mm:ssZ` format in UTC/KST only (for example, 2021-02-02T15:00:00+0900).
* `end_time` - (Optional) You can determine the date and time when the schedule end. If you don't enter `recurrence`, be sure to enter startTime. 
It must be a time later than the current time and a time later than the startTime. Format : `yyyy-MM-ddTHH:mm:ssZ` format in UTC/KST only (for example, 2021-02-02T18:00:00+0900).
* `recurrence` - (Optional) Repeat Settings. You can specify a recurring schedule in crontab format.
* `auto_scaling_group_no` - (Required) The number of the auto scaling group.

~> **NOTE:** Below arguments only support VPC environment.

* `time_zone` - (Optional) the time band for the repeat settings. Valid values are `KST` and `UTC`. Default : `KST`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of Auto Scaling Schedule.
* `name` - The ID of Auto Scaling Schedule (It is the same result as id).
