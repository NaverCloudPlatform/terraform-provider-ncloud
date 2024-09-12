---
subcategory: "Object Storage"
---


# Resource: ncloud_objectstorage_object_acl

Provides Object Storage Object ACL service resource.

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
    bucket_name			= "your-bucket-name"
}

resource "ncloud_objectstorage_object" "testing_object" {
    bucket 				= ncloud_objectstorage_bucket.testing_bucket.bucket_name
    key					= "your-object-key"
    source				= "path/to/file"
}

resource "ncloud_objectstorage_object_acl" "testing_acl" {
    object_id			= ncloud_objectstorage_object.testing_object.id
    rule				= "RULL_TO_APPLY" 
}
```

## Argument Reference

The following arguments are supported:

* `id` - Unique ID for ACL. Has format of `object_acl_${object_id}`.
* `object_id` - (Required) Target object id to create.
* `rule` - (Required) Rule to apply. Value must be one of "private", "public-read", "public-read-write", "authenticated-read".
* `grants` - List of member who grants this rule. Consists of `grantee`, `permission`. Individual `grantee` has `type`, `display_name`, `email-address`, `id`, `uri` attributes.
* `owner_id` - ID of owner.
* `owner_displayname` - Name of owner.

## Import

### `terraform import` command

* Object Storage Object ACL can be imported using the `id`. For example:

```console
$ terraform import ncloud_objectstorage_object_acl.rsc_name object_acl_objectID
```

### `import` block

* In Terraform v1.5.0 and later, use a [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Object Storage Bucket ACL using the `id`. For example:

```terraform
import {
    to = ncloud_objectstorage_object_acl.rsc_name
    id = "object_acl_objectID"
}
```