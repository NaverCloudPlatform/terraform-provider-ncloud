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
  name               = "ex-vpc"
  ipv4_cidr_block    = "10.0.0.0/16"
}

resource "ncloud_subnet" "pri-subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.0.8.0/24"
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "cdb-pri-subnet"
  usage_type     = "GEN"
}

resource "ncloud_redis_config_group" "example" {
  name = "ex-rcg"
  redis_version = "7.0.13-simple"
  description = "example"
}

data "ncloud_redis_image_products" "images_by_code" {
  output_file = "image.json"
}

output "image_list" {
  value = {
    for image in data.ncloud_redis_image_products.images_by_code.image_product_list:
    image.product_name => image.product_code
  }
}

data "ncloud_redis_products" "all" {
  redis_image_product_code = "SW.VRDS.OS.LNX64.ROCKY.0810.REDIS.B050"
  output_file = "products.json"
}

output "product_list" {
  value = {
  for product in data.ncloud_redis_products.all.product_list:
    product.product_name => product.product_code
  }
}

resource "ncloud_redis" "ex-redis" {
  service_name = "ex-redis"
  server_name_prefix = "ex-svr"
  vpc_no = ncloud_vpc.vpc.id
  subnet_no = ncloud_subnet.pri-subnet.id
  config_group_no = ncloud_redis_config_group.example.id
  mode = "SIMPLE"
}
