provider "ncloud" {
  support_vpc = true
  region      = "KR"        # FKR
  access_key  = var.access_key
  secret_key  = var.secret_key
}

resource "ncloud_vpc" "vpc" {
  name            = "tf-vpc"
  ipv4_cidr_block = "172.16.0.0/16"
}

resource "ncloud_subnet" "node_subnet" {
  vpc_no         = ncloud_vpc.vpc.id
  subnet         = "172.16.1.0/24"
  zone           = "KR-1"           # FKR-1
  network_acl_no = ncloud_vpc.vpc.default_network_acl_no
  subnet_type    = "PRIVATE"
  name           = "tf-private-subnet"
  usage_type     = "GEN"
}

data "ncloud_ses_versions" "version" {
  filter {
    name = "value"
    values = [var.ses_version]
    regex = true
  }
}

data "ncloud_ses_software_product" "os_version" {
  filter {
    name = "value"
    values = [var.os_version]
    regex = true
  }
}

data "ncloud_ses_node_product" "product_codes" {
  software_product_code = "SW.VELST.OS.LNX64.CNTOS.0708.B050"
  subnet_no             = 860
  filter {
    name   = "value"
    values = ["SVR.VELST.STAND.C002.M008.NET.SSD.B050.G002"]
  }
}

resource "ncloud_login_key" "loginkey" {
  key_name = var.login_key
}

resource "ncloud_ses_cluster" "cluster" {
  cluster_name                  = "tf-ses"
  software_product_code         = data.ncloud_ses_software_product.os_version.codes.0.value
  vpc_no                        = ncloud_vpc.vpc.id
  search_engine {
    version_code    			= "133"
    user_name       			= "admin"
    user_password   			= "qwe123!@#1"
  }
  manager_node {
    is_dual_manager           = false
    product_code     			= data.ncloud_ses_node_product.product_codes.codes.0.value
    subnet_no        			= ncloud_subnet.node_subnet.id
  }
  data_node {
    product_code       		= data.ncloud_ses_node_product.product_codes.codes.0.value
    subnet_no           		= ncloud_subnet.node_subnet.id
    count            		    = 5
    storage_size        		= 100
  }
  master_node {
    is_master_only_node_activated = false
  }
  login_key_name                = ncloud_login_key.loginkey.key_name
}

data "ncloud_ses_cluster" "cluster" {
  service_group_instance_no = ncloud_ses_cluster.cluster.uuid
}