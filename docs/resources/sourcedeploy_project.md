# Resource: ncloud_sourcedeploy_project

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
* `name` - The name of Sourcedeploy project.
