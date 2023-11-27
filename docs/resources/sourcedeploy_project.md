---
subcategory: "Developer Tools"
---


# Resource: ncloud_sourcedeploy_project

~> **Note** This resource only supports 'public' site.

-> **Note:** This resource is a beta release. Some features may change in the future.

Provides a SourceDeploy Project resource.

## Example Usage

```hcl
resource "ncloud_sourcedeploy_project" "test-deploy-project" {
  name = "test-deploy-project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of Sourcedeploy project.

## Attributes Reference

* `id` - The ID of Sourcedeploy project.

## Import

### `terraform import` command

* SourceDeploy Project can be imported using the `name`. For example:

```console
$ terraform import ncloud_sourcedeploy_project.rsc_name test-deploy
```

### `import` block

* In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import SourceDeploy Project using the `name`. For example:

```terraform
import {
  to = ncloud_sourcedeploy_project.rsc_name
  id = "test-deploy"
}
```
