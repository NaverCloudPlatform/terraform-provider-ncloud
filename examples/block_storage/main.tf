provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "ncloud_regions" "regions" {}
data "ncloud_server_images" "server_images" {}

resource "ncloud_instance" "instance" {
  "server_name"               = "${var.server_name}"
  "server_image_product_code" = "${var.server_image_product_code}"
  "server_product_code"       = "${var.server_product_code}"
}

resource "ncloud_block_storage" "storage" {
  "server_instance_no" = "${ncloud_instance.instance.id}"
  "block_storage_name" = "${var.block_storage_name}"
  "block_storage_size_gb" = "10"
}