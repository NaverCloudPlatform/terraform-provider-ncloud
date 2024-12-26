---
subcategory: "Auto Scaling"
---


# Data Source : ncloud_auto_scaling_adjustment_types
To create an Auto Scaling policy, it's necessary to select an Auto Scaling Adjustment Type. This data source provides a list of available Auto Scaling Adjustment Types


## Example Usage
```hcl
resource "ncloud_launch_configuration" "lc" {
  name                      = "my-lc"
  server_image_product_code = "SW.VSVR.OS.LNX64.ROCKY.0810.B050"
  server_product_code       = "SVR.VSVR.HICPU.C002.M004.NET.SSD.B050.G002"
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
  subnet_no                    = ncloud_subnet.example.subnet_no
  launch_configuration_no      = ncloud_launch_configuration.lc.launch_configuration_no
  min_size                     = 1
  max_size                     = 1
}

resource "ncloud_auto_scaling_policy" "test-policy-CHANG" {
  name                  = "tf-policy"
  adjustment_type_code  = "CHANG"
  scaling_adjustment    = 2
  auto_scaling_group_no = ncloud_auto_scaling_group.auto.auto_scaling_group_no
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

