provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

data "ncloud_mysql_image_products" "all" {
  output_file = "mysql_image_products.json"
}