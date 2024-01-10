---
subcategory: "Redis"
---


# Resource: ncloud_redis_config_group

Provides a Redis Config Group resource.

~> **NOTE:** This resource only supports VPC environment.

## Example Usage

```terraform
resource "ncloud_redis_config_group" "example" {
  name          = "test"
  redis_version = "7.0.13-simple"
  description   = "example"
}
```

## Argument Reference
The following arguments are supported:

* `name` - (Required) Redis Config Group name. Composed of lowercase alphabets, numbers, hyphen (-). Must start with an alphabetic character, and the last character can only be an English letter or number. 3-15 characters.
* `redis_version` - (Required) Redis Service version. These values may change later. For example, `5.0.14-cluster` | `5.0.14-simple` | `7.0.13-cluster` | `7.0.13-simple`
* `description` - (Optional) Redis Config Group description. 1-255 characters.

## Attribute Reference
In addition to all arguments above, the following attributes are exported

* `id` - Redis Config Group instance number.

## Import

### `terraform import` command

* Redis Config Group can be imported using the `name`. For example:

```console
$ terraform import ncloud_redis_config_group.rsc_name test
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Redis Config Group using the `name`. For example:

```terraform
import {
  to = ncloud_redis_config_group.rsc_name
  id = "test"
}
```
