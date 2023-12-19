provider "ncloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

data "ncloud_hadoop_images" "image" {
  product_code = "SW.VCHDP.LNX64.CNTOS.0708.HDP.15.B050"
  output_file = "ncloud_hadoop_images.json"
}

data "ncloud_hadoop_images" "all" {
  output_file = "ncloud_hadoop_images.json"
}