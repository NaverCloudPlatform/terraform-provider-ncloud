---
subcategory: "Developer Tools"
---


# Data Source: ncloud_sourcecommit_repository

~> **Note:** This data source only supports 'public' site.

~> **Note:** This data source is a beta release. Some features may change in the future.

This data source is  useful for getting detail of Sourcecommit repository which is already created, such as getting git address of Sourcecommit repository.

## Example Usage

In the example below, Get git https address information of the repository by specific Sourcecommit repository name.

```hcl
data "ncloud_sourcecommit_repository" "test-repo" {
  name = "test-repo"
}

output "test-repo-git-address" {
    value = data.ncloud_sourcecommit_repository.test-repo.git_https_url
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the specific Repository to retrieve.
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Sourcecommit repository ID.
* `repository_no` - Sourcecommit repository ID.
* `description` - The description of reposiory.
* `creator` - The creator of repository.
* `git_https_url` - The https git address of repository.
* `git_ssh_url` - The ssh git address of repository.
* `file_safer` - whether to use the [File Safer](https://www.ncloud.com/product/security/fileSafer) service 