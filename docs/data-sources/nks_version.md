# Data Source: ncloud_nks_version

Gets a list of available Kubernetes Service versions.

## Example Usage

```hcl
data "ncloud_version" "versions" {}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Custom filter block as described below.
  * `name` - (Required) The name of the field to filter by.
  * `values` - (Required) Set of values that are accepted for the given field.
  * `regex` - (Optional) is `values` treated as a regular expression.

## Attributes Reference

* `versions` - A list of verions
  * `label` - Version label
  * `value` - Version value
