provider "ncloud" {
  support_vpc = true
  access_key  = var.access_key
  secret_key  = var.secret_key
  region      = var.region
}

resource "ncloud_vpc" "vpc" {
  name            = "vpc"
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.0.1.0/24"
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "node-subnet"
  usage_type     = "GEN"
}

resource "ncloud_mongodb" "mongodb" {
    vpc_no            = ncloud_vpc.id
    subnet_no         = ncloud_subnet.subnet.id
    service_name      = "sample-mongodb"
    user_name         = "sample-user"
    user_password     = "sample1234!"
    cluster_type_code = "STAND_ALONE"
}
