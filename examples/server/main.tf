provider "ncloud" {
  access_key = "C9zxQyBQVqcRNQmcAXKn"
  secret_key = "NDgSdssXg4RiMMN1f5dqIe775GJXzSe8hWpd2X3g"
  region     = "${var.region}"
}

data "ncloud_regions" "all" {}

resource "ncloud_instance" "terraform-test" {
  "server_image_product_code" = "${var.server_image_product_code}"
  "server_product_code"       = "${var.server_product_code}"
}
