---
layout: "ncloud"
page_title: "NCLOUD: ncloud_access_control_rules"
sidebar_current: "docs-ncloud-datasource-access-control-rules"
description: |-
  Get access control rule list
---

# Data Source: ncloud_access_control_rules

List of access configuration rules you want to get

## Example Usage

```hcl
data "ncloud_access_control_rules" "test" {
    // access_control_group_configuration_no : You can get one from `ncloud_access_control_group`
    //      or `ncloud_access_control_groups`
	"access_control_group_configuration_no" = "123"
}
```

## Argument Reference

The following arguments are supported:

* `access_control_group_configuration_no` - (Required) Access control group configuration number to search
* `source_name_regex` - (Optional) A regex string to apply to the ACG rule list returned by ncloud
* `output_file` - (Optional) The name of file that can save data source after running `terraform plan`.

## Attributes Reference

* `access_control_rules`
    * `configuration_no` - Access control group configuration number
    * `protocol_type` - Protocol type
        * `code` - Protocol type code
        * `code_name` - Protocol type name
    * `source_ip` - Source IP
    * `destination_port` - Destination Port
    * `source_configuration_no` - Source access control rule configuration no
    * `source_name` - Source access control rule name
    * `description` - Access control rule description
