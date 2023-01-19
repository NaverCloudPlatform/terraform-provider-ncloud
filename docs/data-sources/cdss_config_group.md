# Data Source: ncloud_cdss_config_group

## Example Usage

```hcl
data "ncloud_cdss_config_group" "config_sample" {
  kafka_version_code = data.ncloud_cdss_kafka_version.kafka_version_sample.id

  filter {
    name   = "name"
    values = ["YOUR_CONFIG_GROUP_NAME"]
  }
}
```


## Argument Reference
The following arguments are supported

* `kafka_version_code` - (Required) Cloud Data Streaming Service version to be used.
* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.
## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `id` - ConfigGroup id.
* `config_group_no` - Config group number.
* `name` - ConfigGroup name.
* `description` - ConfigGroup description.