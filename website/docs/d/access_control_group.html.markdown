---
layout: "ncloud"
page_title: "NCLOUD: ncloud_access_control_group"
sidebar_current: "docs-ncloud-datasource-access-control-group"
description: |-
  Get access control group
---

# Data Source: ncloud_access_control_group

When creating a server instance (VM), you can add an access control group (ACG) that you specified to set firewalls. `ncloud_access_control_group` provides details about a specific access control group (ACG) information.


## Example Usage

* Filter by ACG name

```hcl
data "ncloud_access_control_group" "test" {
    # filter by ACG name
	"name" = "acg-name"
}
```


## Argument Reference

The following arguments are supported:

* `configuration_no` - (Optional) List of ACG configuration numbers you want to get
* `name` - (Optional) Name of the ACG you want to get
* `is_default_group` - (Optional) Indicates whether to get default group only

Conditional: Requires `configuration_no` or` name` or `is_default_group`.

## Attributes Reference

* `description` - ACG description
