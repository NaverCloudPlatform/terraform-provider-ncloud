---
subcategory: "Object Storage"
---


# Resource: ncloud_objectstorage_bucket

Provides Object Storage Bucket service resource.

~> **NOTE:** This resource is platform independent. Does not need VPC configuration.

## Example Usage

```terraform
provider "ncloud" {
    support_vpc = true
    access_key = var.access_key
    secret_key = var.secret_key
    region = var.region
}

resource "ncloud_objectstorage_bucket" "testing_bucket" {
    bucket_name				= "your-bucket-name"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required) Bucket name to create. Bucket name must be between 3 and 63 characters long, can contain lowercase letters, numbers, periods, and hyphens. It must start and end with a letter or number, and cannot have consecutive periods.

## Attribute Reference

* `id` - Unique ID for bucket. Since bucket name is already unique in specific region, ID is same as `bucket_name`.
* `creation_date` - Date of when this bucket created.

## Import

### `terraform import` command

* Object Storage Bucket can be imported using the `id`. For example:

```console
$ terraform import ncloud_objectstorage_bucket.rsc_name example-id
```

### `import` block

* In Terraform v1.5.0 and later, use a [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Object Storage Bucket using the `id`. For example:

```terraform
import {
    to = ncloud_objectstorage_bucket.rsc_name
    id = "example-id"
}
```