provider "ncloud" {
    support_vpc = true
    access_key = var.access_key
    secret_key = var.secret_key
    site        = var.site # if needed
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
    bucket_name             = ncloud_objectstorage.bucket.bucket_name
    rule                    = "private"
}

resource "ncloud_objectstorage_object_acl" "object_acl" {
    object_id               = ncloud_objectstorage.object.id
    rule                    = "public-read-write"
}