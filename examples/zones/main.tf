provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

data "ncloud_zones" "zones" {
  region      = "KR"
  output_file = "zones.json"
}

