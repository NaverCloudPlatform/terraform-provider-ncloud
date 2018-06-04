provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "ncloud_zone" "zones" {
  "region_no" = "1"
}
