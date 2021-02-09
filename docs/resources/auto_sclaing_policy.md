# Resource: ncloud_auto_scaling_policy

Provides a ncloud auto scaling policy resource.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `name` - (Required) Auto Scaling Policy name to create.
* `adjustment_type_code` - (Required) Determines how the number of servers is scaled when the scaling policy is performed. Valid values are `CHANG`, `EXACT`, and `PRCNT`
* `scaling_adjustment` - (Required) Specify the adjustment value for the adjustment type. Enter a negative value to decrease when adjustTypeCode is `CHANG` or `PRCNT`.
* `cooldown` - (Optional) The amount of time, in seconds, after a scaling activity completes and before the next scaling activity can start.
* `min_adjustment_step` - (Optional) Change the number of server instances by the minimum adjustment width.
* `auto_scaling_group_no` - (Required) The number of the auto scaling group.

## Attributes Reference

* `id` - The ID of Auto Scaling Policy
* `name` - The ID of Auto Scaling Policy (It is the same result as id)