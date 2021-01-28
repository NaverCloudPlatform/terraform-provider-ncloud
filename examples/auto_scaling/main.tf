provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

data "ncloud_access_control_group" "acg" {
  name = var.acg_name
}

resource "ncloud_launch_configuration" "lc" {
  name = "example-lc"
  server_image_product_code = var.server_image_product_code
  server_product_code       = var.server_product_code
  access_control_group_configuration_no_list = [data.ncloud_access_control_group.acg.id]
}