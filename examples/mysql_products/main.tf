provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

data "ncloud_mysql_products" "all" {
  image_product_code = "SW.VMYSL.OS.LNX64.ROCKY.0810.MYSQL.B050"
  output_file = "mysql_products.json"
}

