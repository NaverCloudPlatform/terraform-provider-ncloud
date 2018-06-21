provider "ncloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "ncloud_zones" "zones" {
  "region_no" = "1"
  "output_file" = "zones.json"
}
