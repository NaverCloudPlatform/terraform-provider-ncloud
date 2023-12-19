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


data "ncloud_nks_versions" "version" {
  filter {
    name = "value"
    values = [var.nks_version]
    regex = true
  }
}
resource "ncloud_login_key" "loginkey" {
  key_name = var.login_key
}

resource "ncloud_nks_cluster" "cluster" {
  cluster_type                = "SVR.VNKS.STAND.C002.M008.NET.SSD.B050.G002"
  k8s_version                 = data.ncloud_nks_versions.version.versions.0.value
  login_key_name              = ncloud_login_key.loginkey.key_name
  name                        = "sample-cluster"
  lb_private_subnet_no        = ncloud_subnet.lb_subnet.id
  kube_network_plugin         = "cilium"
  subnet_no_list              = [ ncloud_subnet.node_subnet.id ]
  vpc_no                      = ncloud_vpc.vpc.id
  zone                        = "KR-1"
  log {
    audit = true
  }
}

data "ncloud_nks_server_images" "image"{
  hypervisor_code = "XEN"
  filter {
    name = "label"
    values = ["ubuntu-20.04"]
    regex = true
  }
}

data "ncloud_nks_server_products" "product"{
  software_code = data.ncloud_nks_server_images.image.images[0].value
  zone = "KR-1"

  filter {
    name = "product_type"
    values = [ "STAND"]
  }

  filter {
    name = "cpu_count"
    values = [ "2"]
  }

  filter {
    name = "memory_size"
    values = [ "8GB" ]
  }
}

resource "ncloud_nks_node_pool" "node_pool" {
  cluster_uuid = ncloud_nks_cluster.cluster.uuid
  node_pool_name = "pool1"
  node_count     = 1
  software_code  = data.ncloud_nks_server_images.image.images[0].value
  product_code   = data.ncloud_nks_server_products.product[0].value
  subnet_no      = ncloud_subnet.subnet.id
  autoscale {
    enabled = true
    min = 1
    max = 2
  }
  label {
    key = "foo"
    value = "bar"
  }
  taints{
    key = "foo"
    value = "bar"
    effect = "NoExecute"
  }
}
