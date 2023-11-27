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
* `description` - (Optional) ConfigGroup description.

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `id` - ConfigGroup id.

## Import

### `terraform import` command

* CDSS Config Group can be imported using the `id`:`kafka_version_code`. For example:

```console
$ terraform import ncloud_cdss_config_group.rsc_name 123:2823006
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import CDSS Config Group using the `id`:`kafka_version_code`. For example:

```terraform
import {
  to = ncloud_cdss_config_group.rsc_name
  id = "123:2823006"
}
```
