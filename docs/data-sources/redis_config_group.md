---
subcategory: "Redis"
---


# Data Source: ncloud_redis_config_group

Provides information about a Redis Config Group.

~> **NOTE:** This only supports VPC environment.

## Example Usage

```terraform
data "ncloud_redis_config_group" "example" {
  name          = "test"
}
```

## Argument Reference

The following arguments are required:

* `name` - (Required) Redis Config Group name. Composed of lowercase alphabets, numbers, hyphen (-). Must start with an alphabetic character, and the last character can only be an English letter or number. 3-15 characters.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `id` - Redis Config Group instance number.
* `redis_version` - Redis Service version.
* `description` - Redis Config Group description. 1-255 characters.
