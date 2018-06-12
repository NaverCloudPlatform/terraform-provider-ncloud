provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "ncloud_server_products" "all" {
  "server_image_product_code" = "SPSW0LINUX000032"
  "output_file" = "server_products.json"
}
