# Data Source: ncloud_ses_versions

Provides list of available Search Engine Service versions.

## Example Usage

```hcl
data "ncloud_ses_versions" "versions" {}

data "ncloud_ses_versions" "opensearch_v133" {
  filter {
    name = "id"
    values = ["133"]
  }
}

data "ncloud_ses_versions" "opensearch_v133" {
  filter {
    name = "type"
    values = ["OpenSearch"]
  }
}

data "ncloud_ses_versions" "opensearch_v133" {
  filter {
    name = "version"
    values = ["1.3.3"]
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

* `versions` - A List of SES Version.

### Search Engine Service Version Reference
`versions` are also exported with the following attributes, when there are relevant: Each element supports the following:

* `id` - The Code of SES Version
* `name` - SES version name
* `type` - SES version type
* `version` - SES version