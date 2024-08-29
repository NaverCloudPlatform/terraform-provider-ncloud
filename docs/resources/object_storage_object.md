---
subcategory: "Object Storage"
---


# Resource: ncloud_objectstorage_object

Provides Object Storage Object service resource.

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

resource "ncloud_objectstorage_object" "testing_object" {
    bucket				= ncloud_objectstorage_bucket.testing_bucket.bucket_name
    key 				= "your-object-key"
    source				= "path/to/file"	
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) Name of the bucket to read the object from. Bucket name must be between 3 and 63 characters long, can contain lowercase letters, numbers, periods, and hyphens. It must start and end with a letter or number, and cannot have consecutive periods.
* `key` - (Required) Full path to the object inside the bucket.
* `source` - (Required) Path to the file you want to upload. 

## Import

### `terraform import` command

* Object Storage Object can be imported using the `id`. For example:

```console
$ terraform import ncloud_objectstorage_object.rsc_name bucket-name/key
```

### `import` block

* In Terraform v1.5.0 and later, use a [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Object Storage Object using the `id`. For example:

```terraform
import {
    to = ncloud_objectstorage_object.rsc_name
    id = "bucket-name/key"
}
```