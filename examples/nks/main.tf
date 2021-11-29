# VPC > User scenario > Scenario 1. Single Public Subnet
# https://docs.ncloud.com/ko/networking/vpc/vpc_userscenario1.html

provider "ncloud" {
  support_vpc = true
  region      = "KR"
  access_key  = var.access_key
  secret_key  = var.secret_key
}

resource "ncloud_vpc" "vpc" {
  name            = "vpc"
  ipv4_cidr_block = "10.0.0.0/16"
}

resource "ncloud_subnet" "node_subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.0.1.0/24"
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "node-subnet"
  usage_type     = "GEN"
}

resource "ncloud_subnet" "lb_subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "10.0.100.0/24"
  zone           = "KR-1"
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "lb-subnet"
  usage_type     = "LOADB"
}


data "ncloud_nks_version" "version" {
  filter {
    name = "value"
    values = [var.nks_version]
    regex = true
  }
}
resource "ncloud_login_key" "loginkey" {
  key_name = "login-key"
}


resource "ncloud_nks_cluster" "cluster" {
  cluster_type                = "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
  k8s_version                 = data.ncloud_nks_version.version.versions.0.value
  login_key_name              = ncloud_login_key.loginkey.key_name
  name                        = "sample-cluster"
  subnet_lb_no                = ncloud_subnet.subnet_lb.id
  subnet_no_list              = [ ncloud_subnet.subnet.id ]
  vpc_no                      = ncloud_vpc.vpc.id
  zone                        = "KR-1"

}
resource "ncloud_nks_node_pool" "node_pool" {
  cluster_name = ncloud_nks_cluster.cluster.name
  node_pool_name = "sample-nodepool"
  node_count     = 2
  product_code   = "SVR.VSVR.STAND.C002.M008.NET.SSD.B050.G002"
  subnet_no      = ncloud_subnet.subnet.id
  autoscale {
    enabled = true
    min = 1
    max = 2
  }
}
