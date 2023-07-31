---
subcategory: "Developer Tools"
---


# Resource: ncloud_sourcecommit_repository

~> **Note:** This resource only supports 'public' site.

~> **Note:** This resource is a beta release. Some features may change in the future.

Provides a Sourcecommit repository resource.

## Example Usage

### Basic Usage

```hcl
resource "ncloud_sourcecommit_repository" "test-repo-basic" {
	name = "repository"
	description = "repository description"
	file_safer = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to create. If omitted, Terraform will force to create new repository and delete previous one.
* `description` - (Optional) description to create.
* `file_safer` - (Optional) A boolean value that determines whether to use the [File Safer](https://www.ncloud.com/product/security/fileSafer) service . Default `false`, Accepted values: `true` | `false` (You must agree to the terms and conditions for use).


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Sourcecommit repository ID.
* `repository_no` - Sourcecommit repository ID.
* `creator` - Sourcecommit repository creator.
* `git_https_url` - Sourcecommit repository https git address.
* `git_ssh_url` - Sourcecommit repository ssh git address.