---
subcategory: "Object Storage"
---

# Data Source: ncloud_objectstorage_bucket

Prvides information about a bucket.

~> **NOTE:** This resource is platform independent. Does not need VPC configuration.

## Example Usage

```terraform
data "ncloud_objectstorage_bucket" "test-bucket" {
    bucket_name = "your-bucket"
}
```

## Argument Reference

The following arguments are required:

* `id` - Unique ID for bucket. Since bucket name is already unique in specific region, ID is same as `bucket_name`.
* `bucket_name` - (Required) Bucket name to create. Bucket name must be between 3 and 63 characters long, can contain lowercase letters, numbers, periods, and hyphens. It must start and end with a letter or number, and cannot have consecutive periods.
* `creation_date` - Date of when this bucket created.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `owner_id` - ID of target bucket owner.
* `owner_displayname` - Display name of target bucket owner.
* `creation_date` - Date information of when this bucket created.