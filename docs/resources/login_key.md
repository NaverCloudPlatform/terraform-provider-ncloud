---
subcategory: "Server"
---


# Resource: ncloud_login_key

Provides a Login key resource.

~> **Note:** All arguments including the private key will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "ncloud_login_key" "loginkey" {
  key_name = "sample key name"
}
```

## Argument Reference

The following arguments are supported:

* `key_name` - (Required) Key name to generate. If the generated key name exists, an error occurs.


## Attributes Reference

* `id` - The ID of login key.
* `private_key` - Generated private key
* `fingerprint` - Fingerprint of the login key

## Import

### `terraform import` command

* Login Key can be imported using the `key_name`. For example:

```console
$ terraform import ncloud_login_key.rsc_name test 
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Login Key using the `key_name`. For example:

```terraform
import {
  to = ncloud_login_key.rsc_name
  id = "test"
}
```
