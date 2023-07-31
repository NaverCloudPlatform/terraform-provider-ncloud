---
subcategory: "Developer Tools"
---


# Resource: ncloud_sourcedeploy_project

~> **Note** This resource only supports 'public' site.

-> **Note:** This resource is a beta release. Some features may change in the future.

Provides a Sourcedeploy project resource.

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

SourceDeploy project can be imported using the project_id, e.g.,

$ terraform import ncloud_sourcedeploy_project.my_project project_id