provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

resource "ncloud_public_ip" "public_ip" {
  "server_instance_no" = "${var.server_instance_no}"
  "region"             = "${var.region}"
  "zone"               = "${var.zone}"
}
