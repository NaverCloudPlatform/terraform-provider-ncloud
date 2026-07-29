---
subcategory: "Sub Account"
---


# Resource: ncloud_subaccount_access_key

Provides an API access key for a Sub Account.

~> **Note:** This resource requires main account API keys. Sub account credentials cannot manage sub accounts.

~> **Note:** The secret key is returned by the API only at creation time and will be stored in the raw state as plain-text.
[Read more about sensitive data in state](https://developer.hashicorp.com/terraform/language/state/sensitive-data).

~> **Note:** A sub account can have at most two access keys, and issuing keys requires `can_api_gateway_access` enabled on the sub account.

## Example Usage

```hcl
resource "ncloud_subaccount" "ci" {
  login_id               = "ci-deployer"
  name                   = "CI Deployer"
  can_api_gateway_access = true
}

resource "ncloud_subaccount_access_key" "ci" {
  sub_account_id = ncloud_subaccount.ci.id
}

output "ci_access_key" {
  value = ncloud_subaccount_access_key.ci.access_key
}

output "ci_secret_key" {
  value     = ncloud_subaccount_access_key.ci.secret_key
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `sub_account_id` - (Required) ID of the sub account to issue the access key for. Changing this creates a new access key.

## Attributes Reference

* `id` - The access key ID.
* `access_key` - Access key ID.
* `secret_key` - Secret key. Available only at creation time.
* `create_time` - Creation time of the access key.

## Import

Import is not supported, because the secret key can never be retrieved after creation.
To rotate a key, create a new `ncloud_subaccount_access_key` and remove the old one.
