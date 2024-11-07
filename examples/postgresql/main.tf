terraform {
  required_version = ">= 0.13"
  required_providers {
    ncloud = {
      source = "navercloudplatform/ncloud"
    }
  }
}

provider "ncloud" {
  region      = var.region
  site        = var.site
  support_vpc = true
  access_key  = var.access_key
  secret_key  = var.secret_key
}

resource "ncloud_vpc" "vpc" {
  name               = "example-vpc"
  ipv4_cidr_block    = "10.0.0.0/16"
}

resource "ncloud_subnet" "subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.0.8.0/24"
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "node-subnet"
  usage_type     = "GEN"
}

data "ncloud_postgresql_image_products" "images_by_code" {
  output_file = "image.json"
}

output "image_list" {
  value = {
    for image in data.ncloud_postgresql_image_products.images_by_code.image_product_list:
    image.product_name => image.product_code
  }
}

data "ncloud_postgresql_products" "all" {
  postgresql_image_product_code = "SW.VSVR.DBMS.LNX64.CNTOS.0708.PSTGR.1403.B050"
  output_file = "products.json"
}

output "product_list" {
  value = {
    for product in data.ncloud_postgresql_products.all.product_list:
    product.product_name => product.product_code
  }
}

resource "ncloud_postgresql" "postgresql" {
  service_name = "example-postgresql"
  server_name_prefix = "example-svr"
  vpc_no = ncloud_vpc.vpc.id
  subnet_no = ncloud_subnet.subnet.id
  user_name = var.user_name
  user_password = var.user_password
  database_name = var.database_name
  client_cir = "0.0.0.0/0"
}

data "ncloud_postgresql" "by_id" {
    id = ncloud_postgresql.postgresql.id
}