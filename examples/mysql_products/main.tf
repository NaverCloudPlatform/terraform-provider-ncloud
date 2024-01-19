provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

data "ncloud_mysql_products" "all" {
  image_product_code = "SW.VDBAS.DBAAS.LNX64.CNTOS.0708.MYSQL.8025.B050"
  output_file = "mysql_products.json"
}

