---
subcategory: "Server"
---


# Data Source: ncloud_access_control_rule

Access configuration rule you want to get.

## Example Usage

```hcl
data "ncloud_access_control_rule" "test" {
  is_default_group = "true"
  destination_port = "22"
}
```

## Argument Reference

The following arguments are supported:

* `access_control_group_configuration_no` - (Optional) Access control group number to search
* `access_control_group_name` - (Optional) Access control group name to search
* `is_default_group` - (Optional) Whether default group
* `source_name_regex` - (Optional) A regex string to apply to the source access control rule list returned by ncloud

## Attributes Reference

* `source_ip` - Source IP
* `destination_port` - Destination Port
* `protocol_type_code` - Protocol type code
* `configuration_no` - Access control rule configuration no
* `protocol_type` - Protocol type code
* `source_configuration_no` - Source access control rule configuration no
* `source_name` - Source access control rule name
* `description` - Access control rule description
