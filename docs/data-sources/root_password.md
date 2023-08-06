---
subcategory: "Server"
---


# Data Source: ncloud_root_password

Gets the password of a root account with the server's login key.

~> **Note:** All arguments including the private key will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "ncloud_root_password" "default" {
  server_instance_no = "server_instance_no" # ${ncloud_server.vm.id}
  private_key = "private_key" # ${ncloud_login_key.key.private_key}
}
```

## Argument Reference

The following arguments are supported:

* `server_instance_no` - (Required) Server instance number
* `private_key` - (Required) Server’s login key (auth key)

## Attributes Reference


* `root_password` - password of a root account
