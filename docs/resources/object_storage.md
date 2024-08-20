---
subcategory: "Object Storage"
---


# Resource: ncloud_objectstorage

Provides Object Storage service resource.

~> **NOTE:** This resources operates in serverless environment. Does not need VPC configuration.

## Example Usage

```terraform
provider "ncloud" {
    support_vpc = true
    access_key = var.access_key
    secret_key = var.secret_key
    region = var.region
}

resource "ncloud_objectstorage_bucket" "bucket" {
    bucket_name             = "bucket"
}

resource "ncloud_objectstorage_object" "object" {
    bucket                  = ncloud_objectstorage.bucket.bucket_name
    key                     = "hello.md"
    source                  = "path/to/file"
}

resource "ncloud_objectstorage_bucket_acl" "bucket_acl" {
    bucket_id               = ncloud_objectstorage.bucket.id
    rule                    = "private"
}

resource "ncloud_objectstorage_object_acl" "object_acl" {
    object_id               = ncloud_objectstorage.object.id
    rule                    = "public-read-write"
}
```

## Argument Reference

The following arguments are supported:

TODO: TBD

## Import

### `terraform import` command

* Object Storage can be imported using the `id`. For example:

```console
$ terraform import ncloud_objectstorage.bucket https://kr.object.ncloudstorage.com/bucket-name
```

### `import` block

* In Terraform v1.5.0 and later, use a [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Object Storage using the `id`. For example:

```terraform
import {
    to = ncloud_objectstorage.bucket
    id = "https://kr.object.ncloudstorage.com/bucket-name"
}
```