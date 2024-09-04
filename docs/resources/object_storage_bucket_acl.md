---
subcategory: "Object Storage"
---


# Resource: ncloud_objectstorage_bucket_acl

Provides Object Storage Bucket ACL service resource.

~> **NOTE:** This resources operates in serverless environment. Does not need VPC configuration.

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

resource "ncloud_objectstorage_bucket_acl" "testing_acl" {
    bucket_name				= ncloud_objectstorage_bucket.testing_bucket.bucket_name
    rule					= "RULL_TO_APPLY"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required) Target bucket id to create(same as bucket name). Bucket name must be between 3 and 63 characters long, can contain lowercase letters, numbers, periods, and hyphens. It must start and end with a letter or number, and cannot have consecutive periods.
* `rule` - (Required) Rule to apply. Value must be one of "private", "public-read", "public-read-write", "authenticated-read".

## Import

### `terraform import` command

* Object Storage Bucket ACL can be imported using the `bucket_name`. For example:

```console
$ terraform import ncloud_objectstorage_bucket_acl.rsc_name bucket_acl_bucket-name
```

### `import` block

* In Terraform v1.5.0 and later, use a [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Object Storage Bucket ACL using the `id`. For example:

```terraform
import {
    to = ncloud_objectstorage_bucket_acl.rsc_name
    bucket_name = "bucket_acl_bucket-name"
}
```