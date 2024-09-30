---
subcategory: "Object Storage"
---

# Data Source: ncloud_objectstorage_object

Prvides information about a object.

~> **NOTE:** This resource is platform independent. Does not need VPC configuration.

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

~> **NOTE:** Since Ncloud Object Stroage uses S3 Compatible SDK, these arguments are served as best-effort.

This data source exports the following attributes in addition to the arguments above:

* `bucket` - Name of bucket where object belongs.
* `key` - Key of object.
* `source` - Path that informs where does object is located in bucket.
* `content_language` - Language the content is in e.g., en-US or en-GB.
* `content_length` - How long the object is.
* `content_type` - Type of the object.
* `body` - Saved content of the object.
* `content_encoding` - Content encodings that have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. Read [w3c content encoding](https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.11) for further information.
* `accept_ranges` - Indicates that a range of bytes was specified.
* `etag` - ETag generated for the object (an MD5 sum of the object content). For plaintext objects or objects encrypted with an AWS-managed key, the hash is an MD5 digest of the object data. For objects encrypted with a KMS key or objects created by either the Multipart Upload or Part Copy operation, the hash is not an MD5 digest, regardless of the method of encryption. More information on possible values can be found on [Common Response Headers](https://docs.aws.amazon.com/AmazonS3/latest/API/RESTCommonResponseHeaders.html). 
* `expiration` - the object expiration is configured, the response includes this header. It includes the expiry-date and rule-id key-value pairs providing object expiration information. The value of the rule-id is URL-encoded. 
* `last_modified` - Date and time when the object was last modified.
* `parts_count` -  The count of parts this object has. This value is only returned if you specify partNumber in your request and the object was uploaded as a multipart upload.
* `version_id` - Unique version ID value for the object, if bucket versioning is enabled.
* `website_redirect_location` - Target URL for website redirect.