provider "ncloud" {
  support_vpc = true
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region
}

resource "ncloud_vpc" "vpc" {
  name               = var.vpc_name
  ipv4_cidr_block    = "10.5.0.0/16"
}

resource "ncloud_subnet" "subnet" {
  vpc_no             = ncloud_vpc.vpc.vpc_no
  name               = var.subnet_name
  subnet             = "10.5.0.0/24"
  zone               = "KR-2"
  network_acl_no     = ncloud_vpc.vpc.default_network_acl_no
  subnet_type        = "PUBLIC"
}

resource "ncloud_mysql" "mysql" {
  subnet_no = ncloud_subnet.subnet.id
  service_name = var.service_name
  server_name_prefix = var.name_prefix
  user_name = var.user_name
  user_password = var.password
  database_name = var.database_name
  host_ip = "192.168.0.1"
}

data "ncloud_mysql" "by_id" {
  id = ncloud_mysql.mysql.id
}
