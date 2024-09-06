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

The following arguments are required:

* `bucket` - (Required) Name of the bucket to read the object from. Bucket name must be between 3 and 63 characters long, can contain lowercase letters, numbers, periods, and hyphens. It must start and end with a letter or number, and cannot have consecutive periods.
* `key` - (Required) Full path to the object inside the bucket.
* `source` - (Required) Path to the file you want to upload. 

The following arguments are optional:

* `bucket_key_enabled` - (Optional) Whether this resource uses Ncloud KMS Keys for SSE.
* `content_encoding` - (Optional) Content encodings that have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. Read [w3c content encoding](https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.11) for further information.
* `checksum_algorithm` - (Optional) Indicates the algorithm used to create the checksum for the object. Valid values: `CRC32`, `CRC32C`, `SHA1`, `SHA256`.
* `content_language` - (Optional) Language the content is in e.g., en-US or en-GB.
* `content_type` - (Optional) Standard MIME type describing the format of the object data, e.g., application/octet-stream. All Valid MIME Types are valid for this input.
* `server_side_encryption` - (Optional) Server-side encryption of the object in Object Storage. Valid values are "AES256".
* `website_redirect_location` - (Optional) Target URL for website redirect.

## Attribute Reference.

~> **NOTE:** Since Ncloud Object Stroage uses S3 Compatible SDK, these arguments are served as best-effort.

This resource exports the following attributes in addition to the arguments above:

* `accept_ranges` - Indicates that a range of bytes was specified.
* `content_length` - Size of the body in bytes.
* `cache_control` - Specifies caching behavior along the request/reply chain.
* `checksum_crc32` - The base64-encoded, 32-bit CRC32 checksum of the object.
* `checksum_crc32c` - The base64-encoded, 32-bit CRC32C checksum of the object.
* `checksum_sha1` - The base64-encoded, 160-bit SHA-1 digest of the object.
* `checksum_sha256` - The base64-encoded, 256-bit SHA-256 digest of the object.
* `etag` - ETag generated for the object (an MD5 sum of the object content). For plaintext objects or objects encrypted with an AWS-managed key, the hash is an MD5 digest of the object data. For objects encrypted with a KMS key or objects created by either the Multipart Upload or Part Copy operation, the hash is not an MD5 digest, regardless of the method of encryption. More information on possible values can be found on [Common Response Headers](https://docs.aws.amazon.com/AmazonS3/latest/API/RESTCommonResponseHeaders.html). 
* `expiration` - the object expiration is configured, the response includes this header. It includes the expiry-date and rule-id key-value pairs providing object expiration information. The value of the rule-id is URL-encoded. 
* `last_modified` - Date and time when the object was last modified.
* `parts_count` -  The count of parts this object has. This value is only returned if you specify partNumber in your request and the object was uploaded as a multipart upload.
* `sse_customer_key_id` - If present, indicates the ID of the Key Management Service (KMS) symmetric encryption customer managed key that was used for the object.
* `version_id` - Unique version ID value for the object, if bucket versioning is enabled.

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