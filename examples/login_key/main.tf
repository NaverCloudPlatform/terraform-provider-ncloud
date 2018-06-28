provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region = "${var.region}"
}

resource "ncloud_login_key" "loginkey" {
  "key_name" = "${var.login_key_name}"
}

resource "ncloud_server" "server" {
  "server_name" = "${var.server_name}"
  "server_image_product_code" = "${var.server_image_product_code}"
  "server_product_code" = "${var.server_product_code}"
  "login_key_name" = "${ncloud_login_key.loginkey.key_name}"
}
