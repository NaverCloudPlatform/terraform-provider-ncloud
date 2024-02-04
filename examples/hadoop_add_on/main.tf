provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

data "ncloud_hadoop_add_on" "addon" {
  image_product_code = var.addon_image_product_code
  cluster_type_code = var.addon_cluster_type_code
  output_file = "ncloud_hadoop_addons.json"
}