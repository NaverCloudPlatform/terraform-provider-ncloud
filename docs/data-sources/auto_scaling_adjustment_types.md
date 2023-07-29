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
  adjustment_type_code = data.ncloud_auto_scaling_adjustment_types.test.types[2].code # 2. 조회된 Adjustment Type 참조
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

* `types` - (Optional) This is the list of Auto Scaling Adjustment Types. Each type consists of the following fields:
    * `code` - (Optional) - This is the code for the type of adjustment. The codes are defined as follows:
        * `CHANG` - This refers to a `Change in Capacity` adjustment.
        * `PRCNT` - This stands for `Percent Change in Capacity` adjustment.
        * `EXACT` - This code is for an `Exact Capacity` adjustment.


* `code_name` - (Optional) - This is a more descriptive name for each code. They are defined as follows:
    * `ChangeInCapacity` - This means the auto-scaling policy adjusts the number of instances by a specified absolute number.
    * `PercentChangeInCapacit` - The auto-scaling policy adjusts the number of instances by a specified percentage.
    * `ExactCapacity` - The auto-scaling policy adjusts the number of instances to a specified exact number.


* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.