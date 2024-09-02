---
subcategory: "Object Storage"
---

# Data Source: ncloud_objectstorage_object

Prvides information about a object.

~> **NOTE:** This resources operates in serverless environment. Does not need VPC configuration.

## Example Usage

```terraform
data "ncloud_objectstorage_object" "test" {
		object_id				= ncloud_objectstorage_object.testing_object.id
}
```

## Argument Reference

The following arguments are required:

* `object_id` - (Required) Object id to get. same as "\${bucket_name}/${object_key}".

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `bucket` - Name of bucket where object belongs.
* `key` - Key of object.
* `source` - Path that informs where does object is located in bucket.
* `content_length` - How long the object is.
* `content_type` - Type of the object.
* `last_modified` - Time that the object last modified.
* `body` - Saved content of the object.