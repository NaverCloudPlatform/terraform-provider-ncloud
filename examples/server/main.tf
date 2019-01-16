provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "ncloud_regions" "regions" {}
data "ncloud_server_images" "server_images" {
}

resource "ncloud_server" "server" {
  "name"               = "${var.server_name}"
  "server_image_product_code" = "${var.server_image_product_code}"
  "server_product_code"       = "${var.server_product_code}"
}
