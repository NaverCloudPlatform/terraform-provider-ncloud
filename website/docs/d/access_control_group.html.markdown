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
	"access_control_group_name" = "acg-name"
}
```

* Filter by most recent ACG

```hcl
data "ncloud_access_control_group" "test" {
    # use the most recent ACG
	"most_recent" = "true"
}
```


## Argument Reference

The following arguments are supported:

* `access_control_group_configuration_no` - (Optional) List of ACG configuration numbers you want to get
* `is_default_group` - (Optional) Indicates whether to get default groups only
* `access_control_group_name` - (Optional) Name of the ACG you want to get
* `most_recent` - (Optional) If more than one result is returned, get the most recent created ACG.

## Attributes Reference

* `access_control_group_description` - ACG description
* `create_date` - Creation date of ACG