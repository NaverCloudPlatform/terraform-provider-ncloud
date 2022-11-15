# Data Source: ncloud_cdss_config_group

## Example Usage

``` hcl
variable "cdss_config_group_uuid" {}

data "ncloud_cdss_config_group" "config_group"{
  id = var.cdss_config_group_uuid
}
```


## Argument Reference
The following arguments are supported

* `id` - (Required) ConfigGroup uuid.
* `kafka_version_code` - (Required) Cloud Data Streaming Service version to be used.
## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `name` - (Required) ConfigGroup name.
* `description` - (Required) ConfigGroup description.