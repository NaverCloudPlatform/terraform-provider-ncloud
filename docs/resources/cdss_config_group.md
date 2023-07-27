---
subcategory: "Cloud Data Streaming Service"
---


# Resource: ncloud_cdss_config_group

## Example Usage

``` hcl
resource "ncloud_cdss_config_group" "config-group" {
  name = "from-tf-config"
  kafka_version_code = "2823006"
  description = "test"
}
```

## Argument Reference
The following arguments are supported:

* `name` - (Required) ConfigGroup name.
* `kafka_version_code` - (Required) Cloud Data Streaming Service version to be used.
* `description` - (Required) ConfigGroup description.

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `id` - ConfigGroup id.
