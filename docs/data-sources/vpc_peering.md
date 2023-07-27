---
subcategory: "VPC"
---


# Data Source: ncloud_vpc_peering

This module can be useful for getting detail of VPC peering created before.

## Example Usage

```hcl
variable "vpc_peering_no" {}

data "ncloud_vpc_peering" "vpc_peering" {
  id = var.vpc_peering_no
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the specific VPC peering.
* `name` - (Optional) The name of the specific VPC Peering to retrieve.
* `source_vpc_name `- (Optional) The name of VPC to which the request to retrieve.
* `target_vpc_name `- (Optional) The name of the VPC that receives the request.
* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vpc_peering_no` - The ID of VPC peering. (It is the same result as `id`)
* `source_vpc_no` - The ID of VPC to which the request.
* `target_vpc_no `- The ID of VPC to receive requests.
* `target_vpc_login_id `- VPC Owner ID to receive requests
* `description` - Description of VPC Peering
* `has_reverse_vpc_peering` - Is existing Reverse VPC Peering.
* `is_between_accounts` - VPC Peering Between Accounts.