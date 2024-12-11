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
    rule				= "public-read'
}
```

## Argument Reference

* `object_id` - (Required) Target object id to create. Has format of `bucket-name/key`.
* `rule` - (Required) Rule to apply. Value must be one of "private", "public-read", "public-read-write", "authenticated-read".

## Attribute Reference

* `id` - Unique ID for ACL. As same as `object_id`. Has format of `bucket-name/key`.
* `grants` - List of member who grants this rule.
  * `grantee` - The person being granted permissions.
    * `type` - Type of grantee
    * `display_name` - Screen name of the grantee.
    * `email_address` - Email address of the grantee.
    * `id` - The canonical user ID of the grantee.
    * `uri` - URI of the grantee group.
  * `permission` - Specifies the permission given to the grantee.
* `owner_id` - ID of owner.
* `owner_displayname` - Name of owner.

## Import

~> **NOTE:** When importing `ncloud_objectstorage_object_acl`, the `rule` value cannot be retrieved automatically. User need to manually set the `rule` value in your Terraform state file after import.

### `terraform import` command

* Object Storage Object ACL can be imported using the `object-id`. For example:

```console
$ terraform import ncloud_objectstorage_object_acl.rsc_name object-id
```

### `import` block

* In Terraform v1.5.0 and later, use a [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Object Storage Object ACL using the `object-id`. For example:

```terraform
import {
    to = ncloud_objectstorage_object_acl.rsc_name
    id = "object-id"
}
```