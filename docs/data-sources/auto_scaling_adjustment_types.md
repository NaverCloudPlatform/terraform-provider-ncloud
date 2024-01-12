---
subcategory: "Auto Scaling"
---


# Data Source : ncloud_auto_scaling_adjustment_types
To create an Auto Scaling policy, it's necessary to select an Auto Scaling Adjustment Type. This data source provides a list of available Auto Scaling Adjustment Types


## Example Usage
```hcl
resource "ncloud_launch_configuration" "lc" {
  name = "my-lc"
  server_image_product_code = "SPSW0LINUX000046"
  server_product_code = "SPSVRSSD00000003"
}

resource "ncloud_auto_scaling_group" "asg" {
  launch_configuration_no = ncloud_launch_configuration.lc.launch_configuration_no
  min_size = 1
  max_size = 1
  zone_no_list = ["2"]
  wait_for_capacity_timeout = "0"
}

resource "ncloud_auto_scaling_policy" "policy" {
  name = "my-policy"
  adjustment_type_code = data.ncloud_auto_scaling_adjustment_types.test.types[2].code
  scaling_adjustment = 2
  auto_scaling_group_no = ncloud_auto_scaling_group.asg.auto_scaling_group_no
}

data "ncloud_auto_scaling_adjustment_types" "test" {
}
```

```hcl
data "ncloud_auto_scaling_adjustment_types" "test" {
  filter {
    name   = "code"
    values = ["EXACT"]
  }
}

output "filtered_types" {
  value = data.ncloud_auto_scaling_adjustment_types.test.types
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

This data source exports the following attributes in addition to the arguments above:

* `types` - This is the list of Auto Scaling Adjustment Types.
    * `code` - This is the code for the type of adjustment. </br>
    Valid Values :</br>
        `CHANG` - This refers to a `Change in Capacity` adjustment.</br>
        `PRCNT` - This stands for `Percent Change in Capacity` adjustment.</br>
        `EXACT` - This code is for an `Exact Capacity` adjustment.
    * `code_name` - This is a more descriptive name for each code.</br>
        Valid Values :</br>
            `ChangeInCapacity` - This means the auto-scaling policy adjusts the number of instances by a specified absolute number.</br>
            `PercentChangeInCapacit` - The auto-scaling policy adjusts the number of instances by a specified percentage.</br>
            `ExactCapacity` - The auto-scaling policy adjusts the number of instances to a specified exact number.

