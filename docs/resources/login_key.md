---
subcategory: "Login Key"
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

* `private_key` - Generated private key
* `fingerprint` - Fingerprint of the login key
* `create_date` - Creation date of the login key