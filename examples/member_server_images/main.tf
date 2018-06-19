provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "ncloud_member_server_images" "all" {
  "output_file" = "server_images.json"
}
