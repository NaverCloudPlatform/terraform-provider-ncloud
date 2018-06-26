---
layout: "ncloud"
page_title: "NCLOUD: ncloud_login_key"
sidebar_current: "docs-ncloud-resource-login-key"
description: |-
  Provides an ncloud login key resource.
---

# ncloud_login_key

Provides an ncloud login key resource.

## Example Usage

```hcl
resource "ncloud_login_key" "loginkey" {
  "key_name" = "sample key name"
}
```

## Argument Reference

The following arguments are supported:

* `key_name` - (Required) Key name to generate. If the generated key name exists, an error occurs.


## Attributes Reference

* `private_key` - generated private key
* `fingerprint` - fingerprint of the login key
* `create_date` - creation date of the login key
