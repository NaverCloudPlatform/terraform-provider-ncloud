# Data Source: ncloud_cdss_kafka_version

## Example Usage

```hcl
data "ncloud_cdss_kafka_version" "sample_01" {
  filter {
    name   = "id"
    values = ["2403005"]
  }
}

data "ncloud_cdss_kafka_version" "sample_02" {
  filter {
    name   = "name"
    values = ["Kafka 2.4.0"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Custom filter block as described below.
    * `name` - (Required) The name of the field to filter by.
    * `values` - (Required) Set of values that are accepted for the given field.
    * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `id` - The ID of kafka version.
* `name` - Kafka version name